package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	ECS "github.com/turnerlabs/fargate/ecs"
)

type ServiceEnvListOperation struct {
	ServiceName string
}

var serviceEnvListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show environment variables",
	Run: func(cmd *cobra.Command, args []string) {
		operation := &ServiceEnvListOperation{
			ServiceName: getServiceName(),
		}

		serviceEnvList(operation)
	},
}

func init() {
	serviceEnvCmd.AddCommand(serviceEnvListCmd)
}

func serviceEnvList(operation *ServiceEnvListOperation) {
	ecs := ECS.New(sess, getClusterName())
	service := ecs.DescribeService(operation.ServiceName)
	envVars := ecs.GetEnvVarsFromTaskDefinition(service.TaskDefinitionArn)

	for _, envVar := range envVars {
		fmt.Printf("%s=%s\n", envVar.Key, envVar.Value)
	}
}
