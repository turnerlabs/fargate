package cmd

import (
	"github.com/spf13/cobra"
	"github.com/turnerlabs/fargate/console"
	"github.com/turnerlabs/fargate/docker"
	"github.com/turnerlabs/fargate/dockercompose"
	ECR "github.com/turnerlabs/fargate/ecr"
	ECS "github.com/turnerlabs/fargate/ecs"
	"github.com/turnerlabs/fargate/git"
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

var serviceDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy applications to services",
	Long: `Deploy applications to services

The Docker container image to use in the service can be specified
via the --image flag. 

The docker-compose.yml format is also supported using the --file flag.
If -f is specified, the image and the environment variables in the 
docker-compose.yml file will be deployed.

If no args are specified, fargate will build a new Docker container image from the 
current working directory and push it to Amazon ECR in a repository named for
the task group. If the current working directory is a git repository, 
the container image will be tagged with the short ref of the HEAD commit. 
If not, a timestamp in the format of YYYYMMDDHHMMSS will be used.
`,
	Example: `
fargate service deploy
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
	serviceDeployCmd.Flags().StringVarP(&flagServiceDeployImage, "image", "i", "", "Docker image to run in the service; if omitted Fargate will build an image from the Dockerfile in the current directory")

	serviceDeployCmd.Flags().StringVarP(&flagServiceDeployDockerComposeFile, "file", "f", "", "Specify a docker-compose.yml file to deploy. The image and environment variables in the file will be deployed.")

	serviceCmd.AddCommand(serviceDeployCmd)
}

func deployService(operation *ServiceDeployOperation) {

	if operation.ComposeFile != "" {
		deployDockerComposeFile(operation)
		return
	}

	ecs := ECS.New(sess, getClusterName())
	service := ecs.DescribeService(operation.ServiceName)

	if operation.Image == "" {
		var tag string

		ecr := ECR.New(sess)
		repositoryUri := ecr.GetRepositoryUri(operation.ServiceName)
		repository := docker.Repository{Uri: repositoryUri}
		username, password := ecr.GetUsernameAndPassword()

		if git.IsCwdGitRepo() {
			tag = git.GetShortSha()
		} else {
			tag = docker.GenerateTag()
		}

		repository.Login(username, password)
		repository.Build(tag)
		repository.Push(tag)

		operation.Image = repository.UriFor(tag)
	}

	taskDefinitionArn := ecs.UpdateTaskDefinitionImage(service.TaskDefinitionArn, operation.Image)
	ecs.UpdateServiceTaskDefinition(operation.ServiceName, taskDefinitionArn)
	console.Info("Deployed %s to service %s", operation.Image, operation.ServiceName)
}

//deploy a docker-compose.yml file to fargate
func deployDockerComposeFile(operation *ServiceDeployOperation) {

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

	//register a new task definition based on the image and environment variables from the compose file
	taskDefinitionArn := ecs.UpdateTaskDefinitionImageAndEnvVars(ecsService.TaskDefinitionArn, dockerService.Image, dockerService.Environment)

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
