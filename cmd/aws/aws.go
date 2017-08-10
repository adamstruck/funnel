package aws

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strings"
)

var log = logger.New("aws cmd")

// Capture AWS Batch config (compute env, job queue, etc.)
var flagConf = DefaultConfig()
var region string

// Funnel's AWS Batch proxy passes the task message to this
// command as a JSON string via a CLI flag.
var rawTask string
var rawTaskFile string

func init() {
	f := runTaskCmd.Flags()
	f.StringVar(&rawTask, "task", "", "Task JSON")
	f.StringVar(&rawTaskFile, "task-file", "", "Task JSON file path")

	d := deployCmd.Flags()
	d.StringVar(&flagConf.Container, "container", flagConf.Container,
		"Funnel worker Docker container to run.")

	Cmd.AddCommand(deployCmd)
	Cmd.AddCommand(runTaskCmd)
	Cmd.AddCommand(proxyCmd)
}

// Cmd is the aws command
var Cmd = &cobra.Command{
	Use: "aws",
}

var deployCmd = &cobra.Command{
	Use: "deploy REGION",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Println("You must provide a region\n")
			return cmd.Help()
		}
		return deploy(flagConf, args[0])
	},
}

var runTaskCmd = &cobra.Command{
	Use: "runtask",
	RunE: func(cmd *cobra.Command, args []string) error {

		if rawTask != "" {
			rawTask = strings.Join(append([]string{rawTask}, args...), " ")

			var anything interface{}
			err := json.Unmarshal([]byte(rawTask), &anything)
			if err != nil {
				return fmt.Errorf("Error parsing Task JSON: %v", err)
			}
			if t, ok := anything.(string); ok {
				rawTask = t
			} else {
				b, err := json.Marshal(&anything)
				if err != nil {
					return fmt.Errorf("Error cleaning Task JSON: %v", err)
				}
				rawTask = string(b)
			}
		}

		if rawTaskFile != "" {
			b, err := ioutil.ReadFile(rawTaskFile)
			if err != nil {
				return err
			}
			rawTask = string(b)
			fmt.Println("3: ", rawTask)
		}

		// Load tes.Task from raw string (comes from CLI flag).
		var task tes.Task
		err := jsonpb.UnmarshalString(rawTask, &task)
		if err != nil {
			return err
		}

		if err := tes.Validate(&task); err != nil {
			return fmt.Errorf("Invalid task message: %v", err)
		}

		conf := config.DefaultConfig()
		log.Configure(conf.Logger)
		return runTask(&task, conf)
	},
}

var proxyCmd = &cobra.Command{
	Use: "proxy",
	RunE: func(cmd *cobra.Command, args []string) error {
		conf := config.DefaultConfig()
		conf.Logger.Level = "debug"
		log.Configure(conf.Logger)
		return runProxy(conf)
	},
}
