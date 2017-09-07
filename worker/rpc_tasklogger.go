package worker

import (
	"context"
	"github.com/ohsu-comp-bio/funnel/config"
	tl "github.com/ohsu-comp-bio/funnel/proto/tasklogger"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/util"
	"google.golang.org/grpc"
	"time"
)

// TODO document behavior of slow consumer of task log updates

// RPCTaskLogger provides access to writing task logs over gRPC to the funnel server.
type RPCTaskLogger struct {
	client        tl.TaskLoggerServiceClient
	taskID        string
	updateTimeout time.Duration
}

// NewRPCTaskLogger returns a TaskLogger that writes task logs over gRPC to the funnel server.
func NewRPCTaskLogger(conf config.Worker, taskID string) (*RPCTaskLogger, error) {
	client, err := newTaskLoggerClient(conf)
	if err != nil {
		return nil, err
	}
	return &RPCTaskLogger{client, taskID, conf.UpdateTimeout}, nil
}

// SetState sets the state of the task.
func (r *RPCTaskLogger) SetState(s tes.State) error {
	_, err := r.client.UpdateTaskState(context.Background(), &tl.UpdateTaskStateRequest{
		Id:    r.taskID,
		State: s,
	})
	return err
}

// StartTime updates the task's start time log.
func (r *RPCTaskLogger) StartTime(t time.Time) {
	r.updateTaskLogs(&tl.UpdateTaskLogsRequest{
		Id: r.taskID,
		TaskLog: &tes.TaskLog{
			StartTime: t.Format(time.RFC3339),
		},
	})
}

// EndTime updates the task's end time log.
func (r *RPCTaskLogger) EndTime(t time.Time) {
	r.updateTaskLogs(&tl.UpdateTaskLogsRequest{
		Id: r.taskID,
		TaskLog: &tes.TaskLog{
			EndTime: t.Format(time.RFC3339),
		},
	})
}

// Outputs updates the task's output file log.
func (r *RPCTaskLogger) Outputs(f []*tes.OutputFileLog) {
	r.updateTaskLogs(&tl.UpdateTaskLogsRequest{
		Id: r.taskID,
		TaskLog: &tes.TaskLog{
			Outputs: f,
		},
	})
}

// Metadata updates the task's metadata log.
func (r *RPCTaskLogger) Metadata(m map[string]string) {
	r.updateTaskLogs(&tl.UpdateTaskLogsRequest{
		Id: r.taskID,
		TaskLog: &tes.TaskLog{
			Metadata: m,
		},
	})
}

// ExecutorStartTime updates an executor's start time log.
func (r *RPCTaskLogger) ExecutorStartTime(i int, t time.Time) {
	r.updateExecutorLogs(&tl.UpdateExecutorLogsRequest{
		Id:   r.taskID,
		Step: int64(i),
		Log: &tes.ExecutorLog{
			StartTime: t.Format(time.RFC3339),
		},
	})
}

// ExecutorEndTime updates an executor's end time log.
func (r *RPCTaskLogger) ExecutorEndTime(i int, t time.Time) {
	r.updateExecutorLogs(&tl.UpdateExecutorLogsRequest{
		Id:   r.taskID,
		Step: int64(i),
		Log: &tes.ExecutorLog{
			EndTime: t.Format(time.RFC3339),
		},
	})
}

// ExecutorExitCode updates an executor's exit code log.
func (r *RPCTaskLogger) ExecutorExitCode(i int, x int) {
	r.updateExecutorLogs(&tl.UpdateExecutorLogsRequest{
		Id:   r.taskID,
		Step: int64(i),
		Log: &tes.ExecutorLog{
			ExitCode: int32(x),
		},
	})
}

// ExecutorPorts updates an executor's ports log.
func (r *RPCTaskLogger) ExecutorPorts(i int, ports []*tes.Ports) {
	r.updateExecutorLogs(&tl.UpdateExecutorLogsRequest{
		Id:   r.taskID,
		Step: int64(i),
		Log: &tes.ExecutorLog{
			Ports: ports,
		},
	})
}

// ExecutorHostIP updates an executor's host IP log.
func (r *RPCTaskLogger) ExecutorHostIP(i int, ip string) {
	r.updateExecutorLogs(&tl.UpdateExecutorLogsRequest{
		Id:   r.taskID,
		Step: int64(i),
		Log: &tes.ExecutorLog{
			HostIp: ip,
		},
	})
}

// AppendExecutorStdout appends to an executor's stdout log.
func (r *RPCTaskLogger) AppendExecutorStdout(i int, s string) {
	r.updateExecutorLogs(&tl.UpdateExecutorLogsRequest{
		Id:   r.taskID,
		Step: int64(i),
		Log: &tes.ExecutorLog{
			Stdout: s,
		},
	})
}

// AppendExecutorStderr appends to an executor's stderr log.
func (r *RPCTaskLogger) AppendExecutorStderr(i int, s string) {
	r.updateExecutorLogs(&tl.UpdateExecutorLogsRequest{
		Id:   r.taskID,
		Step: int64(i),
		Log: &tes.ExecutorLog{
			Stderr: s,
		},
	})
}

func (r *RPCTaskLogger) updateExecutorLogs(up *tl.UpdateExecutorLogsRequest) error {
	ctx, cleanup := context.WithTimeout(context.Background(), r.updateTimeout)
	_, err := r.client.UpdateExecutorLogs(ctx, up)
	if err != nil {
		log.Error("Couldn't update executor logs", err)
	}
	cleanup()
	return err
}

func (r *RPCTaskLogger) updateTaskLogs(up *tl.UpdateTaskLogsRequest) error {
	ctx, cleanup := context.WithTimeout(context.Background(), r.updateTimeout)
	_, err := r.client.UpdateTaskLogs(ctx, up)
	if err != nil {
		log.Error("Couldn't update task logs", err)
	}
	cleanup()
	return err
}

func newTaskLoggerClient(conf config.Worker) (tl.TaskLoggerServiceClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	conn, err := grpc.DialContext(ctx,
		conf.ServerAddress,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		util.PerRPCPassword(conf.ServerPassword),
	)
	if err != nil {
		return nil, err
	}

	return tl.NewTaskLoggerServiceClient(conn), nil
}

// RPCTaskReader implements the TaskReader interface via gRPC calls to the //
// funnel server.
type RPCTaskReader struct {
	client tes.TaskServiceClient
	taskID string
}

// NewRPCTaskReader returns a TaskReader that reads tasks and task states from
// the funnel server via gRPC.
func NewRPCTaskReader(conf config.Worker, taskID string) (*RPCTaskReader, error) {
	client, err := newTaskClient(conf)
	if err != nil {
		return nil, err
	}
	return &RPCTaskReader{client, taskID}, nil
}

// Task returns the task descriptor.
func (r *RPCTaskReader) Task() (*tes.Task, error) {
	return r.client.GetTask(context.Background(), &tes.GetTaskRequest{
		Id:   r.taskID,
		View: tes.TaskView_FULL,
	})
}

// State returns the current state of the task.
func (r *RPCTaskReader) State() tes.State {
	t, _ := r.client.GetTask(context.Background(), &tes.GetTaskRequest{
		Id:   r.taskID,
		View: tes.TaskView_MINIMAL,
	})
	return t.GetState()
}

func newTaskClient(conf config.Worker) (tes.TaskServiceClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	conn, err := grpc.DialContext(ctx,
		conf.ServerAddress,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		util.PerRPCPassword(conf.ServerPassword),
	)
	if err != nil {
		return nil, err
	}

	return tes.NewTaskServiceClient(conn), nil
}
