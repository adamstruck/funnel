package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/batch"
)

func deploy(conf Config, region string) error {
	cli := newBatchSvc(conf)

	sess := session.Must(session.NewSession())
	sess.Config.Region = aws.String(region)
	batchCli := batch.New(sess)

	a, _ := batchCli.DescribeComputeEnvironments(&batch.DescribeComputeEnvironmentsInput{
		ComputeEnvironments: []*string{aws.String(conf.ComputeEnv.Name)},
	})
	if len(a.ComputeEnvironments) == 0 {
		_, err := cli.CreateComputeEnvironment(region)
		if err != nil {
			return err
		}
		fmt.Printf("Created ComputeEnvironment: %s\n", conf.ComputeEnv.Name)
	} else {
		fmt.Printf("ComputeEnvironment: %s already exists\n", conf.ComputeEnv.Name)
	}

	b, _ := batchCli.DescribeJobQueues(&batch.DescribeJobQueuesInput{
		JobQueues: []*string{aws.String(conf.JobQueue.Name)},
	})
	if len(b.JobQueues) == 0 {
		_, err := cli.CreateJobQueue(region)
		if err != nil {
			return err
		}
		fmt.Printf("Created JobQueue: %s\n", conf.JobQueue.Name)
	} else {
		fmt.Printf("JobQueue: %s already exists\n", conf.JobQueue.Name)
	}

	c, _ := batchCli.DescribeJobDefinitions(&batch.DescribeJobDefinitionsInput{
		JobDefinitionName: aws.String(conf.JobDef.Name),
		Status:            aws.String("ACTIVE"),
	})
	if len(c.JobDefinitions) == 0 {
		_, err := cli.CreateJobDef(region)
		if err != nil {
			return err
		}
		fmt.Printf("Created JobDef: %s\n", conf.JobDef.Name)
	} else {
		fmt.Printf("JobDef: %s already exists\n", conf.JobDef.Name)
	}

	return nil
}
