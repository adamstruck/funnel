package worker

import (
	"github.com/golang/protobuf/jsonpb"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/tests/e2e"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func TestInMemoryLogger(t *testing.T) {
	task := &tes.Task{
		Name:        "Hello world",
		Description: "Demonstrates the most basic echo task.",
		Executors: []*tes.Executor{
			{
				ImageName: "alpine",
				Cmd:       []string{"echo", "hello world"},
				Stdout:    "/tmp/stdout",
			},
		},
	}

	conf := e2e.DefaultConfig().Worker

	r, w, _ := os.Pipe()
	os.Stdout = w

	err := Run(task, conf, "in-memory")
	if err != nil {
		t.Fatal(err)
	}

	w.Close()
	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}

	var taskout tes.Task
	err = jsonpb.UnmarshalString(string(out), &taskout)
	if err != nil {
		t.Fatal(err)
	}
	actual := taskout.Logs[0].Logs[0].Stdout
	expected := "hello world\n"

	if actual != expected {
		t.Fatal("\n", "expected:", expected, "\n", "actual:", actual)
	}
}

func TestLogLogger(t *testing.T) {
	task := &tes.Task{
		Id:          "testid-1234",
		Name:        "Hello world",
		Description: "Demonstrates the most basic echo task.",
		Executors: []*tes.Executor{
			{
				ImageName: "alpine",
				Cmd:       []string{"echo", "hello world"},
				Stdout:    "/tmp/stdout",
			},
		},
	}

	conf := e2e.DefaultConfig().Worker
	conf.Logger.Level = "info"
	conf.Logger.Formatter = "json"
	conf.Logger.OutputFile = path.Join(conf.WorkDir, "worker.log")
	conf.Logger.JSONFormat = logger.JSONFormatConfig{
		DisableTimestamp: true,
	}

	err := Run(task, conf, "log")
	if err != nil {
		t.Fatal(err)
	}

	out, err := ioutil.ReadFile(path.Join(conf.WorkDir, "worker.log"))
	if err != nil {
		t.Fatal(err)
	}

	actual := string(out)
	expected := []string{
		"{\"State\":3,\"level\":\"info\",\"msg\":\"SetState\",\"ns\":\"worker\",\"taskID\":\"testid-1234\"}",
		"{\"ExecutorIndex\":0,\"ExecutorPorts\":null,\"level\":\"info\",\"msg\":\"ExecutorPorts\",\"ns\":\"worker\",\"taskID\":\"testid-1234\"}",
		"{\"ExecutorExitCode\":0,\"ExecutorIndex\":0,\"level\":\"info\",\"msg\":\"ExecutorExitCode\",\"ns\":\"worker\",\"taskID\":\"testid-1234\"}",
		"{\"AppendExecutorStdout\":\"hello world\\n\",\"ExecutorIndex\":0,\"level\":\"info\",\"msg\":\"AppendExecutorStdout\",\"ns\":\"worker\",\"taskID\":\"testid-1234\"}",
		"{\"Outputs\":null,\"level\":\"info\",\"msg\":\"Outputs\",\"ns\":\"worker\",\"taskID\":\"testid-1234\"}",
		"{\"State\":5,\"level\":\"info\",\"msg\":\"SetState\",\"ns\":\"worker\",\"taskID\":\"testid-1234\"}",
	}

	for _, s := range expected {
		if !strings.Contains(actual, s) {
			t.Fatal("\n", "missing:", s, "\n", "in:", actual)
		}
	}
}