package compute

import (
	"fmt"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"strings"
)

// Backend is responsible for submitting a task. For some backends such as HtCondor,
// Slurm, and AWS Batch this amounts to scheduling the task. For others such as
// Openstack this may include provisioning a VM and then running the task.
type Backend interface {
	Submit(*tes.Task) error
}

// BackendFactory is a function which creates a new scheduler backend.
// Various backends (Condor, GCE, local, etc.) implement this, which
// allows BackendLoader to load a backend by name (e.g. config.Scheduler).
type BackendFactory func(config.Config) (Backend, error)

// BackendLoader helps load a scheduler backend by name (e.g. config.Scheduler).
type BackendLoader map[string]BackendFactory

// Load finds a scheduler backend by name and returns a new instance.
func (bl BackendLoader) Load(name string, conf config.Config) (Backend, error) {
	name = strings.ToLower(name)
	factory, ok := bl[name]

	if !ok {
		log.Error("Unknown scheduler backend", "name", name)
		return nil, fmt.Errorf("Unknown scheduler backend %s", name)
	}

	return factory(conf)
}
