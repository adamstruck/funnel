package gridengine

import (
	"github.com/ohsu-comp-bio/funnel/compute"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"os"
	"os/exec"
)

// Name of the scheduler backend
const Name = "gridengine"

var log = logger.Sub(Name)

// NewBackend returns a new grid engine Backend instance.
func NewBackend(conf config.Config) (compute.Backend, error) {
	return &Backend{
		name:     Name,
		conf:     conf,
		template: conf.Backends.GridEngine.Template,
	}, nil
}

// Backend represents the grid engine backend.
type Backend struct {
	name     string
	conf     config.Config
	template string
}

// Submit submits a task via "qsub"
func (s *Backend) Submit(task *tes.Task) error {
	log.Debug("Submitting to Grid Engine", "taskID", task.Id)

	submitPath, err := compute.SetupTemplatedHPCSubmit(s.name, s.template, s.conf, task)
	if err != nil {
		return err
	}

	cmd := exec.Command("qsub", submitPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
