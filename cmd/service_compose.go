package cmd

import (
	awsecs "github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/cobra"
	"github.com/turnerlabs/fargate/console"
	"github.com/turnerlabs/fargate/dockercompose"
	ECS "github.com/turnerlabs/fargate/ecs"
)

// ServiceComposeOperation represents a deploy operation
type ServiceComposeOperation struct {
	ComposeFile      string
	ComposeImageOnly bool
	Region           string
	ServiceName      string
	WaitForService   bool
}

type serviceComposeContainer struct {
	Name          string
	DockerService *dockercompose.Service
}

const ignoreDockerComposeLabel = "aws.ecs.fargate.ignore"

var flagServiceComposeFile string
var flagServiceComposeImageOnly bool
var flagServiceComposeWaitForService bool

var composeCmd = &cobra.Command{
	Use:   "compose",
	Short: "Deploy one or more container definitions defined in a docker compose file",
	Long: `Deploy one or more container definitions defined in a docker compose file

Each service name in the docker compose file will be used as the container definition
name. Note that the image, environments variables, and secrets are replaced with what's in the
compose file. Only pre-existing container definitions will be updated.
`,
	Example: `
fargate service compose -f docker-compose.yml
`,
	Run: func(cmd *cobra.Command, args []string) {
		operation := &ServiceComposeOperation{
			ComposeFile:      flagServiceComposeFile,
			ComposeImageOnly: flagServiceComposeImageOnly,
			Region:           region,
			ServiceName:      getServiceName(),
			WaitForService:   flagServiceComposeWaitForService,
		}

		if operation.ComposeFile == "" {
			cmd.Help()
			return
		}

		composeService(operation)
	},
}

func init() {
	composeCmd.Flags().StringVarP(&flagServiceComposeFile, "file", "f", "docker-compose.yml", "Specify a docker-compose.yml file to deploy. The image and environment variables in the file will be deployed.")
	composeCmd.Flags().BoolVar(&flagServiceComposeImageOnly, "image-only", false, "Only deploy the image in the docker-compose.yml file.")
	composeCmd.Flags().BoolVarP(&flagServiceComposeWaitForService, "wait-for-service", "w", false, "Wait for all services to reach a steady state after deploying the new task definition.")

	serviceCmd.AddCommand(composeCmd)
}

func composeService(operation *ServiceComposeOperation) {
	//read the compose file configuration
	composeFile := dockercompose.Read(operation.ComposeFile)
	services := getComposeServicesToDeploy(&composeFile.Data)

	if len(services) == 0 {
		console.IssueExit("Please specify at least one service to deploy")
	}

	taskDefinitionArn := deployDockerComposeContainers(operation.ServiceName, services, operation.ComposeImageOnly)

	if operation.WaitForService {
		waitForService(operation.ServiceName, taskDefinitionArn)
	}
}

//determine which docker-compose service/container to deploy
func getComposeServicesToDeploy(dc *dockercompose.DockerCompose) []*serviceComposeContainer {
	services := []*serviceComposeContainer{}

	for name, service := range dc.Services {
		if service.Labels[ignoreDockerComposeLabel] != "1" {
			services = append(services, &serviceComposeContainer{
				Name:          name,
				DockerService: service,
			})
		}
	}

	return services
}

func deployDockerComposeContainers(serviceName string, services []*serviceComposeContainer, imagesOnly bool) string {
	var containers []*awsecs.ContainerDefinition

	ecs := ECS.New(sess, getClusterName())
	ecsService := ecs.DescribeService(serviceName)

	for _, s := range services {
		envvars := convertDockerComposeEnvVarsToECSEnvVars(s.DockerService)
		secrets := convertDockerComposeSecretsToECSSecrets(s.DockerService)

		container := ecs.CreateContainerDefinition(&ECS.ContainerDefinitionInput{
			EnvVars:    envvars,
			Name:       s.Name,
			Image:      s.DockerService.Image,
			SecretVars: secrets,
		})
		containers = append(containers, container)
	}

	taskDefinitionArn := ecs.UpdateTaskDefinitionContainers(ecsService.TaskDefinitionArn, containers, imagesOnly)

	//update service with new task definition
	ecs.UpdateServiceTaskDefinition(serviceName, taskDefinitionArn)

	console.Info("Deployed revision %s to service %s.", ecs.GetRevisionNumber(taskDefinitionArn), serviceName)

	return taskDefinitionArn
}
