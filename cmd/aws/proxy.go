package aws

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/batch"
	"github.com/golang/protobuf/jsonpb"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/server"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"regexp"
	"strings"
)

const cancelReason = "tes_canceled"

func awsAuthInterceptor() grpc.UnaryServerInterceptor {
	// Return a function that is the interceptor.
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

			if md, ok := metadata.FromContext(ctx); ok {
				if len(md["authorization"]) > 0 {
					raw := md["authorization"][0]
					key, secret, ok := server.ParseBasicAuth(raw)
					if ok {
						ctxv := context.WithValue(ctx, "auth", &credentials.Value{
							AccessKeyID: key,
							SecretAccessKey: secret,
						})
						return handler(ctxv, req)
					}
				}
			}
			return nil, grpc.Errorf(codes.Unauthenticated, "")
	}
}

func batchServer(conf config.Config) *server.Server {
	srv := server.DefaultServer(conf)
	srv.DisableHTTPCache = true
	srv.ServerOptions = []grpc.ServerOption{
		grpc.UnaryInterceptor(awsAuthInterceptor()),
	}
	return srv
}

func runProxy(conf config.Config) error {
	p := proxy{
		client: newBatchSvc(DefaultConfig()),
	}

	srv := batchServer(conf)
	srv.TaskServiceServer = &p
	return srv.Serve(context.Background())
}

type proxy struct {
	client *batchsvc
}

func (p *proxy) CreateTask(ctx context.Context, task *tes.Task) (*tes.CreateTaskResponse, error) {
	log.Debug("CreateTask called", "task", task)

	if err := tes.Validate(task); err != nil {
		log.Error("Invalid task message", "error", err)
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	if len(task.Resources.Zones) == 0 {
		err := errors.New("Task.Resources.Zones cannot be empty")
		log.Error("Invalid task message", "error", err)
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	taskId, err := p.client.CreateJob(ctx, task)
	if err != nil {
		return nil, err
	}
	return &tes.CreateTaskResponse{Id: taskId}, nil
}

func (p *proxy) GetTask(ctx context.Context, req *tes.GetTaskRequest) (*tes.Task, error) {
	// Get the AWS Batch job description.
	result, err := p.client.DescribeJob(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	j := result.Jobs[0]

	// The original TES task message is stored as a parameter in the job description.
	// Load that into task.
	var task tes.Task
	jsonpb.UnmarshalString(*j.Parameters["task"], &task)
	// Translate the task state from AWS status.
	task.State = translateStatus(j.Status, j.StatusReason)

	// Rebuild the task logs.
	for _, attempt := range j.Attempts {
		l := p.buildTaskLog(ctx, req.Id, j, attempt)
		task.Logs = append(task.Logs, l)
	}

	return &task, nil
}

// NOTE: ListTasks does not yet support pagination. For each call, it returns all tasks.
func (p *proxy) ListTasks(ctx context.Context, req *tes.ListTasksRequest) (*tes.ListTasksResponse, error) {

	// AWS Batch's ListJobs endpoint requires a job status query.
	// The loop below makes a query for each of these job statuses, in this order.
	statuses := []string{
		batch.JobStatusSubmitted,
		batch.JobStatusPending,
		batch.JobStatusRunnable,
		batch.JobStatusStarting,
		batch.JobStatusRunning,
		batch.JobStatusSucceeded,
		batch.JobStatusFailed,
	}

	var resp tes.ListTasksResponse

	// Query AWS Batch ListJobs endpoint for each job status.
	for _, status := range statuses {
		page := ""
		// Loop over ListJobs pages.
		for {
			result, err := p.client.ListJobs(ctx, status, page, 100)
			if err != nil {
				return nil, err
			}

			// No results returned, so break to the next job status
			if len(result.JobSummaryList) == 0 {
				break
			}

			for _, summary := range result.JobSummaryList {
				resp.Tasks = append(resp.Tasks, &tes.Task{
					Id: *summary.JobId,
				})
			}

			// No next page, so break to the next job status.
			if result.NextToken == nil {
				break
			}
			page = *result.NextToken
		}
	}

	return &resp, nil
}

func (p *proxy) CancelTask(ctx context.Context, req *tes.CancelTaskRequest) (*tes.CancelTaskResponse, error) {
	_, err := p.client.TerminateJob(ctx, req.Id)
	return &tes.CancelTaskResponse{}, err
}

func (p *proxy) GetServiceInfo(ctx context.Context, info *tes.ServiceInfoRequest) (*tes.ServiceInfo, error) {
	return &tes.ServiceInfo{}, nil
}

// Task/Executor logs are written to CloudWatchLogs as a sequence of events.
// This processes those events and rebuilds them into a TES TaskLog.
func (p *proxy) buildTaskLog(ctx context.Context, id string, j *batch.JobDetail, a *batch.AttemptDetail) *tes.TaskLog {
	t := &tes.TaskLog{}
	arn := *a.Container.TaskArn
	attemptID := strings.Split(arn, "/")[1]
	logs, _ := p.client.GetTaskLogs(ctx, id, *j.JobName, *j.JobId, attemptID)

	for _, l := range logs.Events {
		m := logmsg{}
		if err := json.Unmarshal([]byte(*l.Message), &m); err != nil {
			log.Error("Error processing task log message", err)
			continue
		}

		switch m.Msg {
		case "StartTime":
			t.StartTime = m.StartTime
		case "EndTime":
			t.EndTime = m.EndTime
		case "Outputs":
			t.Outputs = m.Outputs
		case "Metadata":
			t.Metadata = m.Metadata

		case "ExecutorExitCode":
			e := getExec(t, m.ExecutorIndex)
			e.ExitCode = m.ExecutorExitCode
		case "ExecutorHostIP":
			e := getExec(t, m.ExecutorIndex)
			e.HostIp = m.ExecutorHostIP
		case "ExecutorStartTime":
			e := getExec(t, m.ExecutorIndex)
			e.StartTime = m.ExecutorStartTime
		case "ExecutorEndTime":
			e := getExec(t, m.ExecutorIndex)
			e.EndTime = m.ExecutorEndTime
		case "ExecutorPorts":
			e := getExec(t, m.ExecutorIndex)
			e.Ports = m.ExecutorPorts

		case "AppendExecutorStdout":
			e := getExec(t, m.ExecutorIndex)
			e.Stdout += m.AppendExecutorStdout
		case "AppendExecutorStderr":
			e := getExec(t, m.ExecutorIndex)
			e.Stderr += m.AppendExecutorStderr
		}
	}
	return t
}

// Get or create an ExecutorLog entry in the given TaskLog.
func getExec(tl *tes.TaskLog, i int) *tes.ExecutorLog {

	// Grow slice length if necessary
	if len(tl.Logs) <= i {
		desired := i + 1
		tl.Logs = append(tl.Logs, make([]*tes.ExecutorLog, desired-len(tl.Logs))...)
	}

	if tl.Logs[i] == nil {
		tl.Logs[i] = &tes.ExecutorLog{}
	}

	return tl.Logs[i]
}

// Translate AWS job status into TES task state.
func translateStatus(status, reason *string) tes.State {
	if status == nil {
		return tes.State_UNKNOWN
	}

	switch *status {
	case batch.JobStatusSubmitted:
		return tes.State_QUEUED

	case batch.JobStatusPending:
		return tes.State_QUEUED

	case batch.JobStatusRunnable:
		return tes.State_INITIALIZING

	case batch.JobStatusStarting:
		return tes.State_INITIALIZING

	case batch.JobStatusRunning:
		return tes.State_RUNNING

	case batch.JobStatusSucceeded:
		return tes.State_COMPLETE

	case batch.JobStatusFailed:
		if reason != nil && *reason == cancelReason {
			return tes.State_CANCELED
		}
		return tes.State_ERROR

	default:
		return tes.State_UNKNOWN
	}
}

// logmsg represents a task log message written to CloudWatchLogs
// by the funnel worker.
type logmsg struct {
	Level string
	Msg   string
	Ns    string
	Task  string

	StartTime string
	EndTime   string

	ExecutorIndex        int
	ExecutorStartTime    string
	ExecutorEndTime      string
	ExecutorHostIP       string
	AppendExecutorStdout string
	AppendExecutorStderr string
	ExecutorExitCode     int32

	ExecutorPorts []*tes.Ports
	Outputs       []*tes.OutputFileLog
	Metadata      map[string]string
}

// AWS limits the characters allowed in job names,
// so replace invalid characters with underscores.
func safeJobName(s string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9\-_]`)
	return re.ReplaceAllString(s, "_")
}
