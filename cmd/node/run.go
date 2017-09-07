package node

import (
	"context"
	"github.com/imdario/mergo"
	"github.com/ohsu-comp-bio/funnel/compute/basic"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/spf13/cobra"
)

var configFile string
var flagConf = config.Config{}

// runCmd represents the node run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a Funnel node.",
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

func init() {
	flags := runCmd.Flags()
	flags.StringVar(&flagConf.Backends.Basic.Node.ID, "id", flagConf.Backends.Basic.Node.ID, "Node ID")
	flags.StringVar(&flagConf.Backends.Basic.Node.ServerAddress, "server-address", flagConf.Backends.Basic.Node.ServerAddress, "Address of scheduler gRPC endpoint")
	flags.DurationVar(&flagConf.Backends.Basic.Node.Timeout, "timeout", flagConf.Backends.Basic.Node.Timeout, "Timeout in seconds")
	flags.StringVar(&flagConf.Backends.Basic.Node.WorkDir, "work-dir", flagConf.Backends.Basic.Node.WorkDir, "Working Directory")
	flags.StringVar(&flagConf.Backends.Basic.Node.Logger.Level, "log-level", flagConf.Backends.Basic.Node.Logger.Level, "Level of logging")
	flags.StringVar(&flagConf.Backends.Basic.Node.Logger.OutputFile, "log-path", flagConf.Backends.Basic.Node.Logger.OutputFile, "File path to write logs to")
	flags.StringVarP(&configFile, "config", "c", "", "Config File")
}

// Run runs a node with the given config, blocking until the node exits.
func Run(conf config.Config) error {
	logger.Configure(conf.Backends.Basic.Node.Logger)

	if conf.Backends.Basic.Node.ID == "" {
		conf.Backends.Basic.Node.ID = basic.GenNodeID("manual")
	}

	n, err := basic.NewNode(conf)
	if err != nil {
		return err
	}

	n.Run(context.Background())
	return nil
}
