package aws

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/batch"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/golang/protobuf/jsonpb"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"strings"
)

type awsTaskId struct {
	Id     string
	Region string
}

func (ati *awsTaskId) encode() string {
	s := fmt.Sprintf("%s:%s", ati.Region, ati.Id)
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func decodeAwsTaskId(id string) *awsTaskId {
	data, _ := base64.StdEncoding.DecodeString(id)
	parts := strings.SplitN(string(data), ":", 2)
	if len(parts) != 2 {
		return &awsTaskId{}
	}
	return &awsTaskId{parts[1], parts[0]}
}

func newBatchSvc(conf Config) *batchsvc {
	return &batchsvc{
		conf: conf,
	}
}

func newSession(ctx context.Context) *session.Session {
	auth := ctx.Value("auth")
	config := aws.NewConfig()
	config.WithCredentials(credentials.NewStaticCredentialsFromCreds(auth.(credentials.Value)))
	return session.Must(session.NewSession(config))
}

type batchsvc struct {
	conf Config
}

func (b *batchsvc) CreateJob(ctx context.Context, task *tes.Task) (string, error) {

	marshaler := jsonpb.Marshaler{}
	taskJSON, err := marshaler.MarshalToString(task)
	if err != nil {
		return "", err
	}

	sess := newSession(ctx)
	var foundJobDef bool
	var batchCli Batch

	for _, r := range task.Resources.Zones {
		sess.Config.Region = aws.String(r)
		batchCli = batch.New(sess)
		resp, err := batchCli.DescribeJobDefinitions(&batch.DescribeJobDefinitionsInput{
			JobDefinitionName: aws.String(b.conf.JobDef.Name),
		})
		if err != nil {
			return "", err
		}
		if len(resp.JobDefinitions) > 0 {
			foundJobDef = true
			break
		}
	}

	if !foundJobDef {
		err := fmt.Errorf("JobDefinition [%s] not found in %s", b.conf.JobDef.Name, task.Resources.Zones)
		log.Error("CreateTask failed", "error", err)
		return "", grpc.Errorf(codes.NotFound, err.Error())
	}

	req := &batch.SubmitJobInput{
		JobDefinition: aws.String(b.conf.JobDef.Name),
		JobName:       aws.String(safeJobName(task.Name)),
		JobQueue:      aws.String(b.conf.JobQueue.Name),
		Parameters: map[string]*string{
			// Include the entire task message, encoded as a JSON string,
			// in the job parameters. This gets used by the AWS Batch
			// task runner.
			"task": aws.String(taskJSON),
		},
	}

	// convert ram from GB to MiB
	ram := int64(task.Resources.RamGb * 953.674)
	vcpus := int64(task.Resources.CpuCores)
	if ram > b.conf.JobDef.Memory {
		req.ContainerOverrides.Memory = aws.Int64(ram)
	}
	if vcpus > b.conf.JobDef.VCPUs {
		req.ContainerOverrides.Vcpus = aws.Int64(vcpus)
	}

	result, err := batchCli.SubmitJob(req)
	if err != nil {
		return "", err
	}

	taskId := &awsTaskId{*result.JobId, *sess.Config.Region}
	// taskId := &awsTaskId{"1c73954d-47e4-47e2-b7aa-383234cca0e5", *sess.Config.Region}
	log.Debug("Created Task", "Id", taskId.encode(), "decoded", decodeAwsTaskId(taskId.encode()))
	return taskId.encode(), nil
}

func (b *batchsvc) DescribeJob(ctx context.Context, id string) (*batch.DescribeJobsOutput, error) {
	did := decodeAwsTaskId(id)
	sess := newSession(ctx)
	sess.Config.Region = aws.String(did.Region)
	batchCli := batch.New(sess)

	return batchCli.DescribeJobs(&batch.DescribeJobsInput{
		Jobs: []*string{
			aws.String(did.Id),
		},
	})
}

func (b *batchsvc) TerminateJob(ctx context.Context, id string) (*batch.TerminateJobOutput, error) {
	did := decodeAwsTaskId(id)
	sess := newSession(ctx)
	sess.Config.Region = aws.String(did.Region)
	batchCli := batch.New(sess)

	return batchCli.TerminateJob(&batch.TerminateJobInput{
		JobId:  aws.String(did.Id),
		Reason: aws.String(cancelReason),
	})
}

func (b *batchsvc) ListJobs(ctx context.Context, status, token string, size int64) (*batch.ListJobsOutput, error) {
	sess := newSession(ctx)
	// TODO remove hard-coded value
	sess.Config.Region = aws.String("us-west-2")
	batchCli := batch.New(sess)

	return batchCli.ListJobs(&batch.ListJobsInput{
		JobQueue:   aws.String(b.conf.JobQueue.Name),
		JobStatus:  aws.String(status),
		MaxResults: aws.Int64(size),
		NextToken:  aws.String(token),
	})
}

func (b *batchsvc) GetTaskLogs(ctx context.Context, taskID string, name string, jobID string, attemptID string) (*cloudwatchlogs.GetLogEventsOutput, error) {
	did := decodeAwsTaskId(taskID)
	sess := newSession(ctx)
	sess.Config.Region = aws.String(did.Region)
	cloudwatchCli := cloudwatchlogs.New(sess)

	return cloudwatchCli.GetLogEvents(&cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String("/aws/batch/job"),
		LogStreamName: aws.String(name + "/" + jobID + "/" + attemptID),
		StartFromHead: aws.Bool(true),
	})
}

func (b *batchsvc) CreateComputeEnvironment(region string) (*batch.CreateComputeEnvironmentOutput, error) {
	sess := session.Must(session.NewSession())
	sess.Config.Region = aws.String(region)
	batchCli := batch.New(sess)
	ec2Cli := ec2.New(sess)
	iamCli := iam.New(sess)

	sgres, err := ec2Cli.DescribeSecurityGroups(nil)
	if err != nil {
		return nil, err
	}
	securityGroupIds := []string{}
	for _, s := range sgres.SecurityGroups {
		securityGroupIds = append(securityGroupIds, *s.GroupId)
	}

	snres, err := ec2Cli.DescribeSubnets(nil)
	if err != nil {
		return nil, err
	}
	subnets := []string{}
	for _, s := range snres.Subnets {
		subnets = append(subnets, *s.SubnetId)
	}

	iamres, err := iamCli.GetRole(&iam.GetRoleInput{RoleName: aws.String("AWSBatchServiceRole")})
	if err != nil {
		return nil, err
	}
	serviceRole := *iamres.Role.Arn

	return batchCli.CreateComputeEnvironment(&batch.CreateComputeEnvironmentInput{
		ComputeEnvironmentName: aws.String(b.conf.ComputeEnv.Name),
		ComputeResources: &batch.ComputeResource{
			InstanceRole:     aws.String(b.conf.ComputeEnv.InstanceRole),
			InstanceTypes:    convertStringSlice(b.conf.ComputeEnv.InstanceTypes),
			MaxvCpus:         aws.Int64(b.conf.ComputeEnv.MaxVCPUs),
			MinvCpus:         aws.Int64(b.conf.ComputeEnv.MinVCPUs),
			SecurityGroupIds: convertStringSlice(securityGroupIds),
			Subnets:          convertStringSlice(subnets),
			Tags:             convertStringMap(b.conf.ComputeEnv.Tags),
			Type:             aws.String("EC2"),
		},
		ServiceRole: aws.String(serviceRole),
		State:       aws.String("ENABLED"),
		Type:        aws.String("MANAGED"),
	})
}

func (b *batchsvc) CreateJobQueue(region string) (*batch.CreateJobQueueOutput, error) {
	sess := session.Must(session.NewSession())
	sess.Config.Region = aws.String(region)
	batchCli := batch.New(sess)

	var envs []*batch.ComputeEnvironmentOrder
	for i, c := range b.conf.JobQueue.ComputeEnvs {
		envs = append(envs, &batch.ComputeEnvironmentOrder{
			ComputeEnvironment: aws.String(c),
			Order:              aws.Int64(int64(i)),
		})
	}

	return batchCli.CreateJobQueue(&batch.CreateJobQueueInput{
		ComputeEnvironmentOrder: envs,
		JobQueueName:            aws.String(b.conf.JobQueue.Name),
		Priority:                aws.Int64(1),
		State:                   aws.String("ENABLED"),
	})
}

func (b *batchsvc) CreateJobDef(region string) (*batch.RegisterJobDefinitionOutput, error) {
	sess := session.Must(session.NewSession())
	sess.Config.Region = aws.String(region)
	batchCli := batch.New(sess)

	return batchCli.RegisterJobDefinition(&batch.RegisterJobDefinitionInput{
		ContainerProperties: &batch.ContainerProperties{
			Image:      aws.String(b.conf.Container),
			Memory:     aws.Int64(b.conf.JobDef.Memory),
			Vcpus:      aws.Int64(b.conf.JobDef.VCPUs),
			Privileged: aws.Bool(true),
			MountPoints: []*batch.MountPoint{
				{
					SourceVolume:  aws.String("docker_sock"),
					ContainerPath: aws.String("/var/run/docker.sock"),
				},
			},
			Volumes: []*batch.Volume{
				{
					Name: aws.String("docker_sock"),
					Host: &batch.Host{
						SourcePath: aws.String("/var/run/docker.sock"),
					},
				},
			},
			Command: []*string{
				aws.String("aws"),
				aws.String("runtask"),
				aws.String("--task"),
				// This is a template variable that will be replaced with
				// the full TES task message in JSON form.
				aws.String("Ref::task"),
			},
		},
		JobDefinitionName: aws.String(b.conf.JobDef.Name),
		Type:              aws.String("container"),
	})
}

func logErr(err error) {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case batch.ErrCodeClientException:
			log.Error(batch.ErrCodeClientException, aerr.Error())
		case batch.ErrCodeServerException:
			log.Error(batch.ErrCodeServerException, aerr.Error())
		default:
			log.Error("Error", aerr.Error())
		}
	} else {
		log.Error("Error", err)
	}
}

func convertStringSlice(s []string) []*string {
	var ret []*string
	for _, t := range s {
		ret = append(ret, aws.String(t))
	}
	return ret
}

func convertStringMap(s map[string]string) map[string]*string {
	m := map[string]*string{}
	for k, v := range s {
		m[k] = aws.String(v)
	}
	return m
}
