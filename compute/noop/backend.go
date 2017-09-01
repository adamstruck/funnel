package noop

import (
	"github.com/ohsu-comp-bio/funnel/compute"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
)

// Name of the compute backend.
const Name = "noop"

var log = logger.Sub(Name)

// NewBackend returns a new Backend instance.
func NewBackend(conf config.Config) (compute.Backend, error) {
	return &Backend{Name, conf}, nil
}

// Backend is a scheduler backend that doesn't do anything
// which is useful for testing.
type Backend struct {
	name string
	conf config.Config
}

// Submit submits a task. For the noop backend this does nothing.
func (b *Backend) Submit(task *tes.Task) error {
	log.Debug("Submitting to noop", "taskID", task.Id)
	return nil
}
