package cmd

import (
	"github.com/spf13/cobra"
	"github.com/turnerlabs/fargate/console"
	"github.com/turnerlabs/fargate/dockercompose"
	ECS "github.com/turnerlabs/fargate/ecs"
)

// ServiceDeployOperation represents a deploy operation
type ServiceDeployOperation struct {
	ServiceName string
	Image       string
	ComposeFile string
}

const deployDockerComposeLabel = "aws.ecs.fargate.deploy"

var flagServiceDeployImage string
var flagServiceDeployDockerComposeFile string
var flagServiceDeployDockerComposeImageOnly bool

var serviceDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy applications to services",
	Long: `Deploy applications to services

The Docker container image to use in the service can be specified
via the --image flag.

The docker-compose.yml format is also supported using the --file flag.
If -f is specified, the image and the environment variables in the
docker-compose.yml file will be deployed.
`,
	Example: `
fargate service deploy -i 123456789.dkr.ecr.us-east-1.amazonaws.com/my-service:1.0
fargate service deploy -f docker-compose.yml
`,
	Run: func(cmd *cobra.Command, args []string) {
		operation := &ServiceDeployOperation{
			ServiceName: getServiceName(),
			Image:       flagServiceDeployImage,
			ComposeFile: flagServiceDeployDockerComposeFile,
		}
		deployService(operation)
	},
}

func init() {
	serviceDeployCmd.Flags().StringVarP(&flagServiceDeployImage, "image", "i", "", "Docker image to run in the service")

	serviceDeployCmd.Flags().StringVarP(&flagServiceDeployDockerComposeFile, "file", "f", "", "Specify a docker-compose.yml file to deploy. The image and environment variables in the file will be deployed.")

	serviceDeployCmd.Flags().BoolVar(&flagServiceDeployDockerComposeImageOnly, "image-only", false, "Only deploy the image when a docker-compose.yml file is specified.")

	serviceCmd.AddCommand(serviceDeployCmd)
}

func deployService(operation *ServiceDeployOperation) {

	if operation.ComposeFile != "" {
		deployDockerComposeFile(operation)
		return
	}

	ecs := ECS.New(sess, getClusterName())
	service := ecs.DescribeService(operation.ServiceName)

	taskDefinitionArn := ecs.UpdateTaskDefinitionImage(service.TaskDefinitionArn, operation.Image)
	ecs.UpdateServiceTaskDefinition(operation.ServiceName, taskDefinitionArn)
	console.Info("Deployed %s to service %s", operation.Image, operation.ServiceName)
}

//deploy a docker-compose.yml file to fargate
func deployDockerComposeFile(operation *ServiceDeployOperation) {
	var taskDefinitionArn string

	//read the compose file configuration
	composeFile := dockercompose.NewComposeFile(operation.ComposeFile)
	dockerCompose := composeFile.Config()

	//determine which docker-compose service/container to deploy
	_, dockerService := getDockerServiceToDeploy(dockerCompose)
	if dockerService == nil {
		console.IssueExit(`Please indicate which docker container you'd like to deploy using the label "%s: 1"`, deployDockerComposeLabel)
	}

	ecs := ECS.New(sess, getClusterName())
	ecsService := ecs.DescribeService(operation.ServiceName)

	//only update image if --image-only flag is set
	if flagServiceDeployDockerComposeImageOnly {
		//register a new task definition based on the image from the compose file
		taskDefinitionArn = ecs.UpdateTaskDefinitionImage(ecsService.TaskDefinitionArn, dockerService.Image)
	} else {
		//register a new task definition based on the image and environment variables from the compose file
		taskDefinitionArn = ecs.UpdateTaskDefinitionImageAndEnvVars(ecsService.TaskDefinitionArn, dockerService.Image, dockerService.Environment)
	}

	//update service with new task definition
	ecs.UpdateServiceTaskDefinition(operation.ServiceName, taskDefinitionArn)

	console.Info("Deployed %s to service %s as deployment %s", operation.ComposeFile, operation.ServiceName, ecs.GetDeploymentId(taskDefinitionArn))
}

//determine which docker-compose service/container to deploy
func getDockerServiceToDeploy(dc *dockercompose.DockerCompose) (string, *dockercompose.Service) {
	//look for label if there's more than 1
	var service *dockercompose.Service
	name := ""
	for k, v := range dc.Services {
		if len(dc.Services) == 1 {
			service = v
			name = k
			break
		}
		if v.Labels[deployDockerComposeLabel] == "1" {
			service = v
			name = k
			break
		}
	}
	return name, service
}
