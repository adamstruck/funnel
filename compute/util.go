package compute

import (
	"fmt"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/util"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

// DetectFunnelBinaryPath detects the path to the "funnel" binary
func DetectFunnelBinaryPath() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("Failed to detect path of funnel binary")
	}
	return path, err
}

// SetupTemplatedHPCSubmit sets up a task submission in a HPC environment with
// a shared file system. It generates a submission file based on a template for
// schedulers such as SLURM, HTCondor, SGE, PBS/Torque, etc.
func SetupTemplatedHPCSubmit(name string, tpl string, conf config.Config, task *tes.Task) (string, error) {
	var err error

	// TODO document that these working dirs need manual cleanup
	workdir := path.Join(conf.Worker.WorkDir, task.Id)
	workdir, _ = filepath.Abs(workdir)
	err = util.EnsureDir(workdir)
	if err != nil {
		return "", err
	}

	confPath := path.Join(workdir, "worker.conf.yml")
	conf.ToYamlFile(confPath)

	funnelPath, err := DetectFunnelBinaryPath()
	if err != nil {
		return "", err
	}

	submitName := fmt.Sprintf("%s.submit", name)

	submitPath := path.Join(workdir, submitName)
	f, err := os.Create(submitPath)
	if err != nil {
		return "", err
	}

	submitTpl, err := template.New(submitName).Parse(tpl)
	if err != nil {
		return "", err
	}

	var zone string
	zones := task.Resources.GetZones()
	if zones != nil {
		zone = zones[0]
	}

	err = submitTpl.Execute(f, map[string]interface{}{
		"TaskId":     task.Id,
		"Executable": funnelPath,
		"Config":     confPath,
		"WorkDir":    workdir,
		"Cpus":       int(task.Resources.CpuCores),
		"RamGb":      task.Resources.RamGb,
		"DiskGb":     task.Resources.SizeGb,
		"Zone":       zone,
		"Project":    task.Project,
	})
	if err != nil {
		return "", err
	}
	f.Close()

	return submitPath, nil
}
