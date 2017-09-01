package worker

import (
	"context"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"time"
)

// Worker is a type which runs a task.
type Worker interface {
	Run(context.Context)
}

// TaskReader is a type which reads and writes task information
// during task execution.
type TaskReader interface {
	Task() (*tes.Task, error)
	State() tes.State
}

// TaskLogger provides write access to a task's logs.
type TaskLogger interface {
	SetState(tes.State) error
	StartTime(t time.Time)
	EndTime(t time.Time)
	Outputs(o []*tes.OutputFileLog)
	Metadata(m map[string]string)

	ExecutorExitCode(i int, code int)
	ExecutorPorts(i int, ports []*tes.Ports)
	ExecutorHostIP(i int, ip string)
	ExecutorStartTime(i int, t time.Time)
	ExecutorEndTime(i int, t time.Time)

	AppendExecutorStdout(i int, s string)
	AppendExecutorStderr(i int, s string)
}
