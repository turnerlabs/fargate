package cmd

import (
	"github.com/spf13/cobra"
)

const taskLogGroupFormat = "/fargate/task/%s"

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage tasks",
	Long: `Manage tasks
	
Tasks are one-time executions of your container. Instances of your task are run
until you manually stop them either through AWS APIs, the AWS Management
Console, or until they are interrupted for any reason.`,
}

var taskName string

func init() {
	rootCmd.AddCommand(taskCmd)
	taskCmd.PersistentFlags().StringVarP(&taskName, "task", "t", "", `ECS task family name`)
	initPFlag(keyTask, taskCmd)
}