package worker

import (
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"time"
)

// NewLogTaskLogger returns a TaskLogger/TaskReader which writes task logs
// to the given logger instance (level info).
func NewLogTaskLogger(t *tes.Task, l logger.Logger) TaskLogger {
	return &logTaskLogger{l, t}
}

// logTaskLogger implements a TaskLogger/TaskReader which writes task logs
// to the given logger instance (level info).
type logTaskLogger struct {
	log  logger.Logger
	task *tes.Task
}

func (ts *logTaskLogger) Task() (*tes.Task, error) {
	return ts.task, nil
}

func (ts *logTaskLogger) State() tes.State {
	return ts.task.State
}

func (ts *logTaskLogger) SetState(s tes.State) error {
	ts.log.Info("SetState", "State", s)
	ts.task.State = s
	return nil
}

func (ts *logTaskLogger) StartTime(t time.Time) {
	ts.log.Info("StartTime", "StartTime", t)
}

func (ts *logTaskLogger) EndTime(t time.Time) {
	ts.log.Info("EndTime", "EndTime", t)
}

func (ts *logTaskLogger) Outputs(o []*tes.OutputFileLog) {
	ts.log.Info("Outputs", "Outputs", o)
}

func (ts *logTaskLogger) Metadata(m map[string]string) {
	ts.log.Info("Metadata", "Metadata", m)
}

func (ts *logTaskLogger) ExecutorExitCode(i int, code int) {
	ts.log.Info("ExecutorExitCode", "ExecutorIndex", i, "ExecutorExitCode", code)
}

func (ts *logTaskLogger) ExecutorPorts(i int, ports []*tes.Ports) {
	ts.log.Info("ExecutorPorts", "ExecutorIndex", i, "ExecutorPorts", ports)
}

func (ts *logTaskLogger) ExecutorHostIP(i int, ip string) {
	ts.log.Info("ExecutorHostIP", "ExecutorIndex", i, "ExecutorHostIP", ip)
}

func (ts *logTaskLogger) ExecutorStartTime(i int, t time.Time) {
	ts.log.Info("ExecutorStartTime", "ExecutorIndex", i, "ExecutorStartTime", t)
}

func (ts *logTaskLogger) ExecutorEndTime(i int, t time.Time) {
	ts.log.Info("ExecutorEndTime", "ExecutorIndex", i, "ExecutorEndTime", t)
}

func (ts *logTaskLogger) AppendExecutorStdout(i int, s string) {
	ts.log.Info("AppendExecutorStdout", "ExecutorIndex", i, "AppendExecutorStdout", s)
}

func (ts *logTaskLogger) AppendExecutorStderr(i int, s string) {
	ts.log.Info("AppendExecutorStderr", "ExecutorIndex", i, "AppendExecutorStderr", s)
}
