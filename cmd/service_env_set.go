package cmd

import (
	"bufio"
	"github.com/spf13/cobra"
	"github.com/turnerlabs/fargate/console"
	ECS "github.com/turnerlabs/fargate/ecs"
	"os"
)

type ServiceEnvSetOperation struct {
	ServiceName string
	EnvVars     []ECS.EnvVar
}

func (o *ServiceEnvSetOperation) Validate() {
	if len(o.EnvVars) == 0 {
		console.IssueExit("No environment variables specified")
	}
}

func (o *ServiceEnvSetOperation) SetEnvVars(inputEnvVars []string, envVarFile string) {
	if envVarFile != "" {
		file, err := os.Open(envVarFile)
		if err != nil {
			return
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			inputEnvVars = append(inputEnvVars, scanner.Text())
		}
	}
	o.EnvVars = extractEnvVars(inputEnvVars)
}

var flagServiceEnvSetEnvVars []string
var flagServiceEnvSetEnvFile string

var serviceEnvSetCmd = &cobra.Command{
	Use:   "set --env <key=value> [--env <key=value] [--file filename] ...",
	Short: "Set environment variables",
	Long: `Set environment variables

At least one environment variable must be specified via either the --env or
--file flags. You may specify any number of variables on the command line by
repeating --env before each one, or else place multiple variables in a file, one
per line, and specify the filename with --file.  Each --env parameter string or line in the file must be of the form
"key=value", with no quotation marks and no whitespace around the "=" unless you want
literal leading whitespace in the value.`,
	Run: func(cmd *cobra.Command, args []string) {
		operation := &ServiceEnvSetOperation{
			ServiceName: getServiceName(),
		}

		operation.SetEnvVars(flagServiceEnvSetEnvVars, flagServiceEnvSetEnvFile)
		operation.Validate()
		serviceEnvSet(operation)
	},
}

func init() {
	serviceEnvSetCmd.Flags().StringArrayVarP(&flagServiceEnvSetEnvVars, "env", "e", []string{}, "Environment variables to set [e.g. KEY=value]")
	serviceEnvSetCmd.Flags().StringVarP(&flagServiceEnvSetEnvFile, "file", "f", "", "File containing list of environment variables to set, one per line, of the form KEY=value")

	serviceEnvCmd.AddCommand(serviceEnvSetCmd)
}

func serviceEnvSet(operation *ServiceEnvSetOperation) {
	ecs := ECS.New(sess, getClusterName())
	service := ecs.DescribeService(operation.ServiceName)
	taskDefinitionArn := ecs.AddEnvVarsToTaskDefinition(service.TaskDefinitionArn, operation.EnvVars)

	ecs.UpdateServiceTaskDefinition(operation.ServiceName, taskDefinitionArn)

	console.Info("Set %s environment variables:", operation.ServiceName)

	for _, envVar := range operation.EnvVars {
		console.Info("- %s=%s", envVar.Key, envVar.Value)
	}

}
