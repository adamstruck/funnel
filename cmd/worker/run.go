package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/imdario/mergo"
	"github.com/golang/protobuf/jsonpb"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/worker"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strings"
)

var taskID string
var rawTask string
var rawTaskFile string

func init() {
	f := runCmd.Flags()
	f.StringVar(&taskID, "task-id", "", "Task ID")
	f.StringVar(&rawTask, "task-json", "", "Task JSON")
	f.StringVar(&rawTaskFile, "task-file", "", "Task JSON file path")
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a task directly, bypassing the server.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var taskString string
		var task = &tes.Task{}

		switch {
		case rawTaskFile != "":
			b, err := ioutil.ReadFile(rawTaskFile)
			if err != nil {
				return err
			}
			taskString = string(b)

		case rawTask != "":
			rawTask = strings.Join(append([]string{rawTask}, args...), " ")
			var anything interface{}
			err := json.Unmarshal([]byte(rawTask), &anything)
			if err != nil {
				return fmt.Errorf("Error parsing Task JSON: %v", err)
			}
			if t, ok := anything.(string); ok {
				taskString = t
			} else {
				b, err := json.Marshal(&anything)
				if err != nil {
					return fmt.Errorf("Error cleaning Task JSON: %v", err)
				}
				taskString = string(b)
			}
		
		case taskID != "":
			task.Id = taskID
			
		default:
			fmt.Printf("error: No task was provided.\n\n")
			return cmd.Help()
		}

		if taskString != "" {
			err := jsonpb.UnmarshalString(rawTask, task)
			if err != nil {
				return err
			}

			if err := tes.Validate(task); err != nil {
				return fmt.Errorf("Invalid task message: %v", err)
			}
		}

		// parse config file
		conf := config.DefaultConfig()
		config.ParseFile(configFile, &conf)

		// make sure server address and password is inherited by the worker
		conf = config.InheritServerProperties(conf)
		flagConf = config.InheritServerProperties(flagConf)

		// file vals <- cli val
		err := mergo.MergeWithOverwrite(&conf, flagConf)
		if err != nil {
			return err
		}

		return Run(conf.Worker, task)
	},
}

// Run configures and runs a Worker
func Run(conf config.Worker, task *tes.Task) error {
	logger.Configure(conf.Logger)
	w, err := worker.NewDefaultWorker(conf, task.Id)
	if err != nil {
		return err
	}
	w.Run(context.Background())
	return nil
}
