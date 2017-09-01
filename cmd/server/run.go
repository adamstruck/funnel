package server

import (
	"context"
	"github.com/imdario/mergo"
	"github.com/ohsu-comp-bio/funnel/compute"
	"github.com/ohsu-comp-bio/funnel/compute/basic"
	"github.com/ohsu-comp-bio/funnel/compute/basic/gce"
	"github.com/ohsu-comp-bio/funnel/compute/basic/manual"
	"github.com/ohsu-comp-bio/funnel/compute/basic/openstack"
	"github.com/ohsu-comp-bio/funnel/compute/gridengine"
	"github.com/ohsu-comp-bio/funnel/compute/htcondor"
	"github.com/ohsu-comp-bio/funnel/compute/local"
	"github.com/ohsu-comp-bio/funnel/compute/pbs"
	"github.com/ohsu-comp-bio/funnel/compute/slurm"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/server"
	"github.com/spf13/cobra"
)

var log = logger.New("server run cmd")
var configFile string
var flagConf = config.Config{}

// runCmd represents the `funnel server run` command.
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs a Funnel server.",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		// parse config file
		conf := config.DefaultConfig()
		config.ParseFile(configFile, &conf)

		// make sure server address and password is inherited by scheduler nodes and workers
		conf = config.InheritServerProperties(conf)
		flagConf = config.InheritServerProperties(flagConf)

		// file vals <- cli val
		err := mergo.MergeWithOverwrite(&conf, flagConf)
		if err != nil {
			return err
		}

		return Run(conf)
	},
}

// Run runs a default Funnel server.
// This opens a database, and starts an API server, scheduler and task logger.
// This blocks indefinitely.
func Run(conf config.Config) error {
	logger.Configure(conf.Server.Logger)

	computeLoader := compute.BackendLoader{
		basic.Name:      basic.NewBackend,
		gridengine.Name: gridengine.NewBackend,
		htcondor.Name:   htcondor.NewBackend,
		local.Name:      local.NewBackend,
		pbs.Name:        pbs.NewBackend,
		slurm.Name:      slurm.NewBackend,
	}

	backend, err := computeLoader.Load(conf.Backend, conf)
	if err != nil {
		return err
	}

	db, err := server.NewTaskBolt(conf, backend)
	if err != nil {
		log.Error("Couldn't open database", err)
		return err
	}

	srv := server.DefaultServer(db, conf.Server)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start server
	var srverr error
	go func() {
		srverr = srv.Serve(ctx)
		cancel()
	}()

	// Start scheduler
	if conf.Backend == "basic" {
		sloader := basic.BackendLoader{
			gce.Name:       gce.NewBackend,
			manual.Name:    manual.NewBackend,
			openstack.Name: openstack.NewBackend,
		}
		sbackend, err := sloader.Load(conf.Backends.Basic.Backend, conf)
		if err != nil {
			return err
		}
		sched := basic.NewScheduler(db, sbackend, conf.Backends.Basic)
		err = sched.Start(ctx)
		if err != nil {
			return err
		}
	}

	// Block
	<-ctx.Done()
	if srverr != nil {
		log.Error("Server error", srverr)
	}
	return srverr
}
