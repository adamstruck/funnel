package local

import (
	"context"
	"github.com/ohsu-comp-bio/funnel/compute"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/worker"
)

// Name of the compute backend.
const Name = "local"

var log = logger.Sub(Name)

// NewBackend returns a new Backend instance.
func NewBackend(conf config.Config) (compute.Backend, error) {
	return &Backend{Name, conf}, nil
}

// Backend represents the local backend.
type Backend struct {
	name string
	conf config.Config
}

// Submit submits a task. For the Local backend this results in the task
// running immediately.
func (b *Backend) Submit(task *tes.Task) error {
	log.Debug("Submitting to local", "taskID", task.Id)
	w := worker.NewDefaultWorker(b.conf.Worker, task.Id)
	go w.Run(context.Background())
	return nil
}
