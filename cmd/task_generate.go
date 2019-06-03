package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/turnerlabs/fargate/console"
	"github.com/turnerlabs/fargate/dockercompose"
	ECS "github.com/turnerlabs/fargate/ecs"
)

var flagTaskGenerateDockerComposeFile string

var taskGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a docker compose file based on a task definition",
	Run:   generate,
	Example: `
# with a fargate.yml present	
fargate task generate

# specify task definition family, and compose file
fargate task generate -t task -f docker-compose.yml

# specify specific task definition family with revision
fargate task generate -t my-app:42
`,
}

func init() {
	taskGenerateCmd.Flags().StringVarP(&flagTaskGenerateDockerComposeFile, "file", "f", "docker-compose.yml", "Ouptut Docker Compose file")

	taskCmd.AddCommand(taskGenerateCmd)
}

func generate(cmd *cobra.Command, args []string) {
	ecs := ECS.New(sess, "")

	//lookup latest/active task definition from family
	td := ecs.DescribeTaskDefinition(getTaskName()).TaskDefinition
	if len(td.ContainerDefinitions) == 0 {
		console.IssueExit("No container found in task definition")
	}
	container := td.ContainerDefinitions[0]

	//initialize a new compose file
	composeFile := dockercompose.New(flagTaskGenerateDockerComposeFile)

	//add service for the 1st container
	service := composeFile.AddService(*container.Name)
	service.Image = *container.Image

	//ports
	for _, p := range container.PortMappings {
		service.Ports = append(service.Ports, dockercompose.Port{
			Published: 80,
			Target:    *p.ContainerPort,
		})
	}

	//add envvars
	for _, e := range container.Environment {
		service.Environment[*e.Name] = *e.Value
	}

	//add secrets
	for _, s := range container.Secrets {
		service.Secrets[*s.Name] = *s.ValueFrom
	}

	//indicate that this container should be deployed
	service.Labels[deployDockerComposeLabel] = "1"

	//write object to file
	yes := true
	if _, err := os.Stat(composeFile.File); err == nil {
		fmt.Print(composeFile.File + " already exists. Overwrite? ")
		yes = askForConfirmation()
	}
	if yes {
		composeFile.Write()
		fmt.Println("wrote", composeFile.File)
	}
	fmt.Println("done")
}
