package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	flagTaskLogsFilter            string
	flagTaskLogsEndTime           string
	flagTaskLogsStartTime         string
	flagTaskLogsFollow            bool
	flagTaskLogsTasks             []string
	flagTaskLogsContainerName     string
	flagTaskLogsTime              bool
	flagTaskLogsNoLogStreamPrefix bool
)

var taskLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Show logs from tasks",
	Long: `Show logs from tasks

Assumes a cloudwatch log group with the following convention: fargate/task/<task>
where task is specified via --task, or fargate.yml, or environment variable options

Return either a specific segment of task logs or tail logs in real-time using
the --follow option. Logs are prefixed by their log stream name which is in the
format of fargate/<container-name>/<task-id>.

--container-name allows you to specifiy the container within the task definition to get logs for
(defaults to app)

Follow will continue to run and return logs until interrupted by Control-C. If
--follow is passed --end cannot be specified.

Logs can be returned for specific tasks by passing a task
ID via the --task flag. Pass --task with a task ID multiple times in order to
retrieve logs from multiple specific tasks.

A specific window of logs can be requested by passing --start and --end options
with a time expression. The time expression can be either a duration or a
timestamp:

	- Duration (e.g. -1h [one hour ago], -1h10m30s [one hour, ten minutes, and
		thirty seconds ago], 2h [two hours from now])
	- Timestamp with optional timezone in the format of YYYY-MM-DD HH:MM:SS [TZ];
		timezone will default to UTC if omitted (e.g. 2017-12-22 15:10:03 EST)

You can filter logs for specific term by passing a filter expression via the
--filter flag. Pass a single term to search for that term, pass multiple terms
to search for log messages that include all terms.

--time includes the log timestamp in the output

--no-prefix excludes the log stream prefix from the output
`,
	Example: `
fargate task logs
fargate task logs --follow
fargate task logs --start "-10m"
fargate task logs --time --no-prefix
fargate task logs --task my-task-dev --container-name my-container
`,
	Run: func(cmd *cobra.Command, args []string) {

		operation := &GetLogsOperation{
			LogGroupName:      fmt.Sprintf(taskLogGroupFormat, getTaskName()),
			Filter:            flagTaskLogsFilter,
			Follow:            flagTaskLogsFollow,
			Namespace:         flagTaskLogsContainerName,
			IncludeTime:       flagTaskLogsTime,
			NoLogStreamPrefix: flagTaskLogsNoLogStreamPrefix,
		}

		operation.AddTasks(flagTaskLogsTasks)
		operation.AddStartTime(flagTaskLogsStartTime)
		operation.AddEndTime(flagTaskLogsEndTime)

		GetLogs(operation)
	},
}

func init() {
	taskCmd.AddCommand(taskLogsCmd)

	taskLogsCmd.Flags().BoolVarP(&flagTaskLogsFollow, "follow", "f", false, "Poll logs and continuously print new events")
	taskLogsCmd.Flags().StringVar(&flagTaskLogsFilter, "filter", "", "Filter pattern to apply")
	taskLogsCmd.Flags().StringVar(&flagTaskLogsStartTime, "start", "", "Earliest time to return logs (e.g. -1h, 2018-01-01 09:36:00 EST")
	taskLogsCmd.Flags().StringVar(&flagTaskLogsEndTime, "end", "", "Latest time to return logs (e.g. 3y, 2021-01-20 12:00:00 EST")
	taskLogsCmd.Flags().StringSliceVarP(&flagTaskLogsTasks, "task", "t", []string{}, "Show logs from specific task (can be specified multiple times)")
	taskLogsCmd.Flags().StringVarP(&flagTaskLogsContainerName, "container-name", "n", "app", "name of container in task defintion to get logs for")
	taskLogsCmd.PersistentFlags().BoolVarP(&flagTaskLogsTime, "time", "T", false, "append time to logs")
	taskLogsCmd.PersistentFlags().BoolVarP(&flagTaskLogsNoLogStreamPrefix, "no-prefix", "", false, "don't include log stream prefix in output")
}
