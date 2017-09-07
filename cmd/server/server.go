package server

import (
	"github.com/spf13/cobra"
)

// Cmd represents the `funnel server` command set
var Cmd = &cobra.Command{
	Use:   "server",
	Short: "Funnel server commands.",
}

func init() {
	flags := Cmd.PersistentFlags()
	flags.StringVarP(&configFile, "config", "c", "", "Config File")
	flags.StringVar(&flagConf.Server.HostName, "hostname", flagConf.Server.HostName, "Host name or IP")
	flags.StringVar(&flagConf.Server.RPCPort, "rpc-port", flagConf.Server.RPCPort, "RPC Port")
	flags.StringVar(&flagConf.Server.HTTPPort, "http-port", flagConf.Server.HTTPPort, "HTTP Port")
	flags.StringVar(&flagConf.Server.Logger.Level, "log-level", flagConf.Server.Logger.Level, "Level of logging")
	flags.StringVar(&flagConf.Server.Logger.OutputFile, "log-path", flagConf.Server.Logger.OutputFile, "File path to write logs to")
	flags.StringVar(&flagConf.Server.Logger.Formatter, "log-format", flagConf.Server.Logger.Formatter, "Log format. ['json', 'text']")
	flags.StringVar(&flagConf.Server.DBPath, "db-path", flagConf.Server.DBPath, "Database path")
	flags.StringVar(&flagConf.Backend, "backend", flagConf.Backend, "Name of scheduler backend to enable")
	Cmd.AddCommand(runCmd)
}
