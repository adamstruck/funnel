package htcondor

import (
	"github.com/ohsu-comp-bio/funnel/compute"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"os"
	"os/exec"
)

// Name of the compute backend
const Name = "htcondor"

var log = logger.Sub(Name)

// NewBackend returns a new HTCondor Backend instance.
func NewBackend(conf config.Config) (compute.Backend, error) {
	return &Backend{
		name:     "htcondor",
		conf:     conf,
		template: conf.Backends.HTCondor.Template,
	}, nil
}

// Backend represents the HTCondor backend.
type Backend struct {
	name     string
	conf     config.Config
	template string
}

// Submit submits a task via "condor_submit"
func (s *Backend) Submit(task *tes.Task) error {
	log.Debug("Submitting to htcondor", "taskID", task.Id)

	submitPath, err := compute.SetupTemplatedHPCSubmit(s.name, s.template, s.conf, task)
	if err != nil {
		return err
	}

	cmd := exec.Command("condor_submit", submitPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
