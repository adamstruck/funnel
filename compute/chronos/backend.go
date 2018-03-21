package chronos

import (
	"context"

	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/events"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
)

// NewBackend returns a new local Backend instance.
func NewBackend(ctx context.Context, conf config.Chronos, reader tes.ReadOnlyServer, writer events.Writer) (*Backend, error) {
	return nil, nil
}


// Backend represents the local backend.
type Backend struct {
	client   *Client
	conf     config.Chronos
	event    events.Writer
	database tes.ReadOnlyServer
}

// WriteEvent writes an event to the compute backend.
// Currently, TASK_CREATED and TASK_STATE events are handled.
func (b *Backend) WriteEvent(ctx context.Context, ev *events.Event) error {
	switch ev.Type {
	case events.Type_TASK_CREATED:
		return b.Submit(ev.GetTask())

	case events.Type_TASK_STATE:
		if ev.GetState() == tes.State_CANCELED {
			return b.Cancel(ctx, ev.Id)
		}
	}
	return nil
}

// Submit submits a task to Chronos.
func (b *Backend) Submit(task *tes.Task) error {
	// https://mesos.github.io/chronos/docs/api.html#adding-a-docker-job
	return nil
}

// Cancel removes a task from Chronos.
func (b *Backend) Cancel(ctx context.Context, taskID string) error {
	// https://mesos.github.io/chronos/docs/api.html#deleting-a-job
	return nil
}
