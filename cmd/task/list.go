package task

import (
	"fmt"
	"io"
	"strings"

	"github.com/ohsu-comp-bio/funnel/client"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"golang.org/x/net/context"
)

// List runs the "task list" CLI command, which connects to the server,
// calls ListTasks() and requests the given task view.
// Output is written to the given writer.
func List(server, taskView, pageToken, stateFilter string, tagsFilter []string, pageSize uint32, all bool, writer io.Writer) error {
	cli, err := client.NewClient(server)
	if err != nil {
		return err
	}

	view, err := getTaskView(taskView)
	if err != nil {
		return err
	}

	state, err := getTaskState(stateFilter)
	if err != nil {
		return err
	}

	tags := make(map[string]string)
	for _, v := range tagsFilter {
		parts := strings.Split(v, "=")
		if len(parts) != 2 {
			return fmt.Errorf("tags must be of the form: KEY=VALUE")
		}
		tags[parts[0]] = parts[1]
	}

	fmt.Fprintln(writer, "{")
	fmt.Fprintln(writer, `  "tasks": [`)

	var resp *tes.ListTasksResponse
	firstRequest := true

	for {
		resp, err = cli.ListTasks(context.Background(), &tes.ListTasksRequest{
			View:      tes.TaskView(view),
			PageToken: pageToken,
			PageSize:  pageSize,
			State:     state,
			Tags:      tags,
		})
		if err != nil {
			return err
		}

		// set up variables for next iteration if `--all` was requested
		pageToken = resp.NextPageToken
		firstRequest = false

		// append to last seen previous task
		if all && resp.NextPageToken != "" && !firstRequest {
			fmt.Fprintf(writer, ",\n")
		}

		for i, t := range resp.Tasks {
			tj, _ := tes.MarshalToString(t)
			tj = strings.Replace(tj, "\n", "\n    ", -1)
			if i != len(resp.Tasks)-1 {
				fmt.Fprintf(writer, "    %s,\n", tj)
			} else {
				fmt.Fprintf(writer, "    %s", tj)
			}
		}

		// all tasks have been collected
		if !all || (all && pageToken == "") {
			fmt.Fprintf(writer, "\n")
			break
		}
	}

	if resp.NextPageToken != "" {
		fmt.Fprintln(writer, "  ],")
		fmt.Fprintf(writer, "  \"nextPageToken\": \"%s\"\n", resp.NextPageToken)
	} else {
		fmt.Fprintln(writer, "  ]")
	}

	fmt.Fprintln(writer, "}")

	return nil
}
