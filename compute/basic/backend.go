package basic

import (
	"context"
	"github.com/ohsu-comp-bio/funnel/compute"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
)

// Name of the basic backend
const Name = "funnel-basic-scheduler"

var log = logger.Sub(Name)

// NewBackend returns a new funnel basic basic Backend instance.
func NewBackend(conf config.Config) (compute.Backend, error) {
	client, err := NewClient(conf.Backends.Basic)
	if err != nil {
		log.Error("Can't connect scheduler client", err)
		return nil, err
	}

	return &Backend{Name, conf, client}, nil
}

// Backend represents the funnel basic basic backend.
type Backend struct {
	name   string
	conf   config.Config
	client Client
}

// Submit submits a task via gRPC call to the funnel basic basic
func (b *Backend) Submit(task *tes.Task) error {
	log.Debug("Submitting to basic basic", "taskID", task.Id)
	_, err := b.client.QueueTask(context.Background(), task)
	if err != nil {
		log.Error("Failed to submit task to the scheduler queue", err, "taskID", task.Id)
		return err
	}
	return nil
}
