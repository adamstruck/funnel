package ccc

import (
	"fmt"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	pbf "github.com/ohsu-comp-bio/funnel/proto/funnel"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/scheduler"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

var log = logger.New("ccc")

// prefix is a string prefixed to condor worker IDs, so that condor
// workers can be identified by ShouldStartWorker() below.
// TODO move to worker metadata to be consistent with GCE
const prefix = "ccc-worker-"

// Plugin provides the HTCondor scheduler backend plugin.
var Plugin = &scheduler.BackendPlugin{
	Name:   "ccc",
	Create: NewBackend,
}

// NewBackend returns a new HTCondor Backend instance.
func NewBackend(conf config.Config) (scheduler.Backend, error) {
	return scheduler.Backend(&Backend{conf}), nil
}

// Backend represents the CCC backend.
type Backend struct {
	conf config.Config
}

// Schedule schedules a task on the HTCondor queue and returns a corresponding Offer.
func (s *Backend) Schedule(t *tes.Task) *scheduler.Offer {
	log.Debug("Running ccc scheduler")

	disk := s.conf.Worker.Resources.DiskGb
	if disk == 0.0 {
		disk = t.GetResources().GetSizeGb()
	}

	cpus := s.conf.Worker.Resources.Cpus
	if cpus == 0 {
		cpus = t.GetResources().GetCpuCores()
	}

	ram := s.conf.Worker.Resources.RamGb
	if ram == 0.0 {
		ram = t.GetResources().GetRamGb()
	}

	// Setup strategy
	// Worker config take precedence
	metadata := make(map[string]string)
	metadata["strategy"] = "routed_file"
	if val, ok := t.Tags["strategy"]; ok {
		metadata["strategy"] = val
	}
	// for _, store := range s.conf.Worker.Storage {
	// 	if store.CCC.Valid() {
	// 		if store.CCC.Strategy != "" {
	// 			metadata["strategy"] = store.CCC.Strategy
	// 		}
	// 	}
	// }

	// TODO could we call condor_submit --dry-run to test if a task would succeed?
	w := &pbf.Worker{
		Id: prefix + t.Id,
		Resources: &pbf.Resources{
			Cpus:   cpus,
			RamGb:  ram,
			DiskGb: disk,
		},
		Metadata: metadata,
	}
	return scheduler.NewOffer(w, t, scheduler.Scores{})
}

// ShouldStartWorker is part of the Scaler interface and returns true
// when the given worker needs to be started by Backend.StartWorker
func (s *Backend) ShouldStartWorker(w *pbf.Worker) bool {
	return strings.HasPrefix(w.Id, prefix) &&
		w.State == pbf.WorkerState_UNINITIALIZED
}

// StartWorker submits a task via "condor_submit" to start a new worker.
func (s *Backend) StartWorker(w *pbf.Worker) error {
	log.Debug("Starting ccc worker", "worker", w)
	var err error

	// TODO document that these working dirs need manual cleanup
	workdir := path.Join(s.conf.WorkDir, w.Id)
	workdir, _ = filepath.Abs(workdir)
	err = os.MkdirAll(workdir, 0755)
	if err != nil {
		return err
	}

	wc := s.conf
	wc.Worker.ID = w.Id
	wc.Worker.Timeout = 5 * time.Second
	wc.Worker.Resources.Cpus = w.Resources.Cpus
	wc.Worker.Resources.RamGb = w.Resources.RamGb
	wc.Worker.Resources.DiskGb = w.Resources.DiskGb

	storage := setupCCCStorage(w.Metadata["strategy"], wc.Storage)
	wc.Worker.Storage = storage

	confPath := path.Join(workdir, "worker.conf.yml")
	wc.ToYamlFile(confPath)

	workerPath, err := scheduler.DetectWorkerPath()
	if err != nil {
		return err
	}

	submitPath := path.Join(workdir, "condor.submit")
	f, err := os.Create(submitPath)
	if err != nil {
		return err
	}

	submitTpl, err := template.New("condor.submit").Parse(`
universe       = vanilla
executable     = {{.Executable}}
arguments      = worker --config worker.conf.yml
getenv         = True
log            = {{.WorkDir}}/condor-event-log
error          = {{.WorkDir}}/funnel-worker-stderr
output         = {{.WorkDir}}/funnel-worker-stdout
input          = {{.Config}}
{{.Resources}}
should_transfer_files   = YES
when_to_transfer_output = ON_EXIT
queue
`)
	if err != nil {
		return err
	}

	err = submitTpl.Execute(f, map[string]string{
		"Executable": workerPath,
		"WorkDir":    workdir,
		"Config":     confPath,
		"Resources":  resolveCondorResourceRequest(int(w.Resources.Cpus), w.Resources.RamGb, w.Resources.DiskGb),
	})
	if err != nil {
		return err
	}
	f.Close()

	cmd := exec.Command("condor_submit")
	stdin, err := os.Open(submitPath)
	if err != nil {
		return err
	}
	cmd.Stdin = stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func resolveCondorResourceRequest(cpus int, ram float64, disk float64) string {
	var resources = []string{}
	if cpus != 0 {
		resources = append(resources, fmt.Sprintf("request_cpus   = %d", cpus))
	}
	if ram != 0.0 {
		resources = append(resources, fmt.Sprintf("request_memory = %f GB", ram))
	}
	if disk != 0.0 {
		// Convert GB to KiB
		resources = append(resources, fmt.Sprintf("request_disk   = %f", disk*976562))
	}
	return strings.Join(resources, "\n")
}

func setupCCCStorage(strategy string, sconf []*config.StorageConfig) []*config.StorageConfig {
	// Use 'strategy' to setup  CCC Storage Backend
	var storage []*config.StorageConfig
	var sc *config.StorageConfig
	for _, s := range sconf {
		if s.CCC.Valid() {
			sc = s
			sc.CCC.Strategy = strategy
			storage = append(storage, sc)
		}
	}
	return storage
}
