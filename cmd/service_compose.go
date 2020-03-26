package cmd

import (
	"sync"

	"github.com/spf13/cobra"
	"github.com/turnerlabs/fargate/console"
	"github.com/turnerlabs/fargate/dockercompose"
)

// ServiceComposeOperation represents a deploy operation
type ServiceComposeOperation struct {
	ComposeFile      string
	ComposeImageOnly bool
	Region           string
	WaitForService   bool
}

// ServiceComposeService represents a Fargate service from a compose file
type ServiceComposeService struct {
	Name    string
	Service *dockercompose.Service
}

const ignoreDockerComposeLabel = "aws.ecs.fargate.ignore"

var flagServiceComposeFile string
var flagServiceComposeImageOnly bool
var flagServiceComposeWaitForService bool

var composeCmd = &cobra.Command{
	Use:   "compose",
	Short: "Deploy one or more services defined in a docker compose file",
	Long: `Deploy one or more services defined in a docker compose file

Each service name in the docker compose file will be used as the Fargate
service name. Note that environments variables and secrets are replaced
with what's in the compose file.
`,
	Example: `
fargate service compose -f docker-compose.yml
`,
	Run: func(cmd *cobra.Command, args []string) {
		operation := &ServiceComposeOperation{
			Region:           region,
			ComposeFile:      flagServiceComposeFile,
			ComposeImageOnly: flagServiceComposeImageOnly,
			WaitForService:   flagServiceComposeWaitForService,
		}

		if operation.ComposeFile == "" {
			cmd.Help()
			return
		}

		composeServices(operation)
	},
}

func init() {
	composeCmd.Flags().StringVarP(&flagServiceComposeFile, "file", "f", "docker-compose.yml", "Specify a docker-compose.yml file to deploy. The image and environment variables in the file will be deployed.")
	composeCmd.Flags().BoolVar(&flagServiceComposeImageOnly, "image-only", false, "Only deploy the image in the docker-compose.yml file.")
	composeCmd.Flags().BoolVarP(&flagServiceComposeWaitForService, "wait-for-service", "w", false, "Wait for all services to reach a steady state after deploying the new task definition.")

	serviceCmd.AddCommand(composeCmd)
}

func composeServices(operation *ServiceComposeOperation) {
	//read the compose file configuration
	composeFile := dockercompose.Read(operation.ComposeFile)
	services := getComposeServicesToDeploy(&composeFile.Data)

	if len(services) == 0 {
		console.IssueExit("Please specify at least one service to deploy")
	}

	//run deploy for each service in a goroutine
	var wg sync.WaitGroup
	for _, service := range services {
		wg.Add(1)
		go deployComposeService(service, operation, &wg)
	}

	//wait for all services to deploy
	wg.Wait()
}

//determine which docker-compose service/container to deploy
func getComposeServicesToDeploy(dc *dockercompose.DockerCompose) []*ServiceComposeService {
	services := []*ServiceComposeService{}

	for name, service := range dc.Services {
		if service.Labels[ignoreDockerComposeLabel] != "1" {
			services = append(services, &ServiceComposeService{
				Name:    name,
				Service: service,
			})
		}
	}

	return services
}

func deployComposeService(service *ServiceComposeService, operation *ServiceComposeOperation, wg *sync.WaitGroup) {
	defer wg.Done()

	taskDefinitionArn := deployDockerComposeService(service.Name, service.Service, operation.ComposeImageOnly)

	if operation.WaitForService {
		waitForService(service.Name, taskDefinitionArn)
	}
}
