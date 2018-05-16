package cmd

import (
	"github.com/turnerlabs/fargate/console"
	ECS "github.com/turnerlabs/fargate/ecs"
	"github.com/spf13/cobra"
)

type ServiceRestartOperation struct {
	ServiceName string
}

var serviceRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart service",
	Long: `Restart service

Creates a new set of tasks for the service and stops the previous tasks. This
is useful if your service needs to reload data cached from an external source,
for example.`,
	Run: func(cmd *cobra.Command, args []string) {
		operation := &ServiceRestartOperation{
			ServiceName: getServiceName(),
		}

		restartService(operation)
	},
}

func init() {
	serviceCmd.AddCommand(serviceRestartCmd)
}

func restartService(operation *ServiceRestartOperation) {
	ecs := ECS.New(sess, getClusterName())

	ecs.RestartService(operation.ServiceName)
	console.Info("Restarted %s", operation.ServiceName)
}
