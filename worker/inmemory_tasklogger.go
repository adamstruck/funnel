package worker

import (
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"time"
)

// NewInMemoryTaskLogger returns a TaskLogger/TaskReader that stores the state
// of the task in memory
func NewInMemoryTaskLogger(t *tes.Task) TaskLogger {
	t.Logs = append(t.Logs, &tes.TaskLog{})
	return &inMemoryTaskLogger{t}
}

// inMemoryTaskLogger implements a TaskLogger/TaskReader that stores the state
// of the task in memory
type inMemoryTaskLogger struct {
	task *tes.Task
}

func (ts *inMemoryTaskLogger) Task() (*tes.Task, error) {
	return ts.task, nil
}

func (ts *inMemoryTaskLogger) State() tes.State {
	return ts.task.State
}

func (ts *inMemoryTaskLogger) SetState(s tes.State) error {
	ts.task.State = s
	return nil
}

func (ts *inMemoryTaskLogger) StartTime(t time.Time) {
	ts.task.Logs[0].StartTime = t.String()
}

func (ts *inMemoryTaskLogger) EndTime(t time.Time) {
	ts.task.Logs[0].EndTime = t.String()
}

func (ts *inMemoryTaskLogger) Outputs(o []*tes.OutputFileLog) {
	ts.task.Logs[0].Outputs = o
}

func (ts *inMemoryTaskLogger) Metadata(m map[string]string) {
	ts.task.Logs[0].Metadata = m
}

func (ts *inMemoryTaskLogger) ExecutorExitCode(i int, code int) {
	exec := getExec(ts.task.Logs[0], i)
	exec.ExitCode = int32(code)
	ts.task.Logs[0].Logs[i] = exec
}

func (ts *inMemoryTaskLogger) ExecutorPorts(i int, ports []*tes.Ports) {
	exec := getExec(ts.task.Logs[0], i)
	exec.Ports = ports
	ts.task.Logs[0].Logs[i] = exec
}

func (ts *inMemoryTaskLogger) ExecutorHostIP(i int, ip string) {
	exec := getExec(ts.task.Logs[0], i)
	exec.HostIp = ip
	ts.task.Logs[0].Logs[i] = exec
}

func (ts *inMemoryTaskLogger) ExecutorStartTime(i int, t time.Time) {
	exec := getExec(ts.task.Logs[0], i)
	exec.StartTime = t.String()
	ts.task.Logs[0].Logs[i] = exec
}

func (ts *inMemoryTaskLogger) ExecutorEndTime(i int, t time.Time) {
	exec := getExec(ts.task.Logs[0], i)
	exec.EndTime = t.String()
	ts.task.Logs[0].Logs[i] = exec
}

func (ts *inMemoryTaskLogger) AppendExecutorStdout(i int, s string) {
	exec := getExec(ts.task.Logs[0], i)
	exec.Stdout += s
	ts.task.Logs[0].Logs[i] = exec
}

func (ts *inMemoryTaskLogger) AppendExecutorStderr(i int, s string) {
	exec := getExec(ts.task.Logs[0], i)
	exec.Stderr += s
	ts.task.Logs[0].Logs[i] = exec
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
