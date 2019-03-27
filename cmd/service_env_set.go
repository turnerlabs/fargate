package cmd

import (
	"github.com/spf13/cobra"
	"github.com/turnerlabs/fargate/console"
	ECS "github.com/turnerlabs/fargate/ecs"
)

type ServiceEnvSetOperation struct {
	ServiceName string
	EnvVars     []ECS.EnvVar
	SecretVars  []ECS.Secret
}

func (o *ServiceEnvSetOperation) Validate() {
	if len(o.EnvVars) == 0 && len(o.SecretVars) == 0 {
		console.IssueExit("No environment variables or secrets specified")
	}
}

func (o *ServiceEnvSetOperation) SetEnvVars(inputEnvVars []string, envVarFile string) {
	o.EnvVars = processEnvVarArgs(inputEnvVars, envVarFile)
}

func (o *ServiceEnvSetOperation) SetSecretVars(inputSecretVars []string, secretVarFile string) {
	o.SecretVars = processSecretVarArgs(inputSecretVars, secretVarFile)
}

func processEnvVarArgs(inputEnvVars []string, envVarFile string) []ECS.EnvVar {
	if envVarFile != "" {
		inputEnvVars = append(inputEnvVars, readVarFile(envVarFile)...)
	}
	return extractEnvVars(inputEnvVars)
}

func processSecretVarArgs(inputSecretVars []string, secretVarFile string) []ECS.Secret {
	var result []ECS.Secret

	if secretVarFile != "" {
		inputSecretVars = append(inputSecretVars, readVarFile(secretVarFile)...)
	}

	for _, envVar := range extractEnvVars(inputSecretVars) {
		result = append(result, ECS.Secret{
			Key:       envVar.Key,
			ValueFrom: envVar.Value,
		})
	}

	return result
}

var flagServiceEnvSetEnvVars []string
var flagServiceEnvSetEnvFile string
var flagServiceEnvSetSecretVars []string
var flagServiceEnvSetSecretFile string

var serviceEnvSetCmd = &cobra.Command{
	Use:   "set --env <key=value> [--env <key=value>] [--file filename] [--secret <key=valueFrom>] [--secret-file filename]...",
	Short: "Set environment variables",
	Long: `Set environment variables

At least one environment variable must be specified via either the --env, --secret,
--file, or --secret-file flags. You may specify any number of variables on the command line by
repeating --env or --secret before each one, or else place multiple variables in a file, one
per line, and specify the filename with --file or --secret-file.

Each --env and --secret parameter string or line in the file must be of the form
"key=value", with no quotation marks and no whitespace around the "=" unless you want
literal leading whitespace in the value.  Additionally, the "key" side must be
a legal shell identifier, which means it must start with an ASCII letter A-Z or
underscore and consist of only letters, digits, and underscores.`,
	Run: func(cmd *cobra.Command, args []string) {
		operation := &ServiceEnvSetOperation{
			ServiceName: getServiceName(),
		}

		operation.SetEnvVars(flagServiceEnvSetEnvVars, flagServiceEnvSetEnvFile)
		operation.SetSecretVars(flagServiceEnvSetSecretVars, flagServiceEnvSetSecretFile)
		operation.Validate()
		serviceEnvSet(operation)
	},
}

func init() {
	serviceEnvSetCmd.Flags().StringArrayVarP(&flagServiceEnvSetEnvVars, "env", "e", []string{}, "Environment variables to set [e.g. KEY=value]")
	serviceEnvSetCmd.Flags().StringVarP(&flagServiceEnvSetEnvFile, "file", "f", "", "File containing list of environment variables to set, one per line, of the form KEY=value")
	serviceEnvSetCmd.Flags().StringArrayVar(&flagServiceEnvSetSecretVars, "secret", []string{}, "Secret variables to set [e.g. KEY=valueFrom]")
	serviceEnvSetCmd.Flags().StringVar(&flagServiceEnvSetSecretFile, "secret-file", "", "File containing list of secret variables to set, one per line, of the form KEY=valueFrom")

	serviceEnvCmd.AddCommand(serviceEnvSetCmd)
}

func serviceEnvSet(operation *ServiceEnvSetOperation) {
	ecs := ECS.New(sess, getClusterName())
	service := ecs.DescribeService(operation.ServiceName)
	taskDefinitionArn := ecs.AddEnvVarsToTaskDefinition(service.TaskDefinitionArn, operation.EnvVars, operation.SecretVars)

	ecs.UpdateServiceTaskDefinition(operation.ServiceName, taskDefinitionArn)

	if len(operation.EnvVars) > 0 {
		console.Info("Set %s environment variables:", operation.ServiceName)

		for _, envVar := range operation.EnvVars {
			console.Info("- %s=%s", envVar.Key, envVar.Value)
		}
	}

	if len(operation.SecretVars) > 0 {
		console.Info("Set %s secret variables:", operation.ServiceName)

		for _, envVar := range operation.SecretVars {
			console.Info("- %s=%s", envVar.Key, envVar.ValueFrom)
		}
	}
}
