package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	ECS "github.com/turnerlabs/fargate/ecs"
)

var flagTaskRegisterImage string
var flagTaskRegisterDockerComposeFile string
var flagTaskRegisterEnvVars []string
var flagTaskRegisterEnvFile string
var flagTaskRegisterSecretVars []string
var flagTaskRegisterSecretFile string

//represents a task register operation
type taskRegisterOperation struct {
	Cluster     string
	Task        string
	Image       string
	EnvVars     []string
	EnvFile     string
	ComposeFile string
	SecretVars  []string
	SecretFile  string
}

var taskRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Registers a new task definition revision for the specified docker image or environment variables based on the latest revision of the task family and returns the new revision number.",
	Run: func(cmd *cobra.Command, args []string) {

		operation := taskRegisterOperation{
			Cluster:     getClusterName(),
			Task:        getTaskName(),
			Image:       flagTaskRegisterImage,
			EnvVars:     flagTaskRegisterEnvVars,
			EnvFile:     flagTaskRegisterEnvFile,
			ComposeFile: flagTaskRegisterDockerComposeFile,
			SecretVars:  flagTaskRegisterSecretVars,
			SecretFile:  flagTaskRegisterSecretFile,
		}

		//valid cli arg combinations
		nonComposeOptions := (flagTaskRegisterImage != "" ||
			len(flagTaskRegisterEnvVars) > 0 ||
			flagTaskRegisterEnvFile != "" ||
			len(flagTaskRegisterSecretVars) > 0 ||
			flagTaskRegisterSecretFile != "")

		if (flagTaskRegisterDockerComposeFile != "" && nonComposeOptions) ||
			(flagTaskRegisterDockerComposeFile == "" && !nonComposeOptions) {
			cmd.Help()
			return
		}

		registerTask(operation)
	},
	Example: `
fargate task register --image 123456789.dkr.ecr.us-east-1.amazonaws.com/my-app:0.1.0
fargate task register --image 123456789.dkr.ecr.us-east-1.amazonaws.com/my-app:0.1.0 --env FOO=bar --env BAR=baz
fargate task register --image 123456789.dkr.ecr.us-east-1.amazonaws.com/my-app:0.1.0 --env FOO=bar --secret BAZ=qux
fargate task register --env-file dev.env
fargate task register --secret-file secrets.env
fargate task register --file docker-compose.yml
`,
}

func init() {
	taskRegisterCmd.Flags().StringVarP(&flagTaskRegisterImage, "image", "i", "", "Docker image to register")

	taskRegisterCmd.Flags().StringArrayVarP(&flagTaskRegisterEnvVars, "env", "e", []string{}, "Environment variables to set [e.g. -e KEY=value -e KEY2=value]")

	taskRegisterCmd.Flags().StringVar(&flagTaskRegisterEnvFile, "env-file", "", "File containing list of environment variables to set, one per line, of the form KEY=value")

	taskRegisterCmd.Flags().StringVarP(&flagTaskRegisterDockerComposeFile, "file", "f", "", "Docker Compose file containing image and environment variables to register.")

	taskRegisterCmd.Flags().StringArrayVar(&flagTaskRegisterSecretVars, "secret", []string{}, "Secret variables to set [e.g. --secret KEY=valueFrom --secret KEY2=valueFrom]")

	taskRegisterCmd.Flags().StringVar(&flagTaskRegisterSecretFile, "secret-file", "", "File containing list of secret variables to set, one per line, of the form KEY=valueFrom")

	taskCmd.AddCommand(taskRegisterCmd)
}

func registerTask(op taskRegisterOperation) {

	//are we registering from cli args or a compose file?
	image := op.Image
	var envvars []ECS.EnvVar
	var secrets []ECS.Secret
	replaceVars := false

	if op.ComposeFile != "" {
		dockerService := getDockerServiceFromComposeFile(op.ComposeFile)
		image = dockerService.Image
		envvars = convertDockerComposeEnvVarsToECSEnvVars(dockerService)
		secrets = convertDockerComposeSecretsToECSSecrets(dockerService)
		replaceVars = true

	} else {
		//read env file (if specified) and combine with other envvars
		envvars = processEnvVarArgs(op.EnvVars, op.EnvFile)

		//read secrets file (if specified) and combine with other secret vars
		secrets = processSecretVarArgs(op.SecretVars, op.SecretFile)

		//don't replace, just add, update where exists
		replaceVars = false
	}

	//update and register new task definition
	ecs := ECS.New(sess, op.Cluster)
	newTD := ecs.UpdateTaskDefinitionImageAndEnvVars(op.Task, image, envvars, replaceVars, secrets)

	//output new revision
	fmt.Println(ecs.GetRevisionNumber(newTD))
}
