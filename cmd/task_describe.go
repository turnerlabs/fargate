package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/turnerlabs/fargate/console"
	"github.com/turnerlabs/fargate/dockercompose"
	ECS "github.com/turnerlabs/fargate/ecs"
)

var taskDescribeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe a task definition in docker compose format",
	Run:   describe,
	Example: `
# with a fargate.yml present	
fargate task describe

# specify task definition family
fargate task describe -t my-app

# specify specific task definition family with revision
fargate task describe -t my-app:42
`,
}

func init() {
	taskCmd.AddCommand(taskDescribeCmd)
}

func describe(cmd *cobra.Command, args []string) {
	ecs := ECS.New(sess, "")

	//lookup latest/active task definition from family
	td := ecs.DescribeTaskDefinition(getTaskName()).TaskDefinition
	if len(td.ContainerDefinitions) == 0 {
		console.IssueExit("No container found in task definition")
	}
	container := td.ContainerDefinitions[0]

	//initialize a new compose file
	composeFile := dockercompose.New("")

	//add service for the 1st container
	service := composeFile.AddService(*container.Name)
	service.Image = *container.Image

	//ports
	for _, p := range container.PortMappings {
		service.Ports = append(service.Ports, dockercompose.Port{
			Published: *p.ContainerPort,
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

	yaml, err := composeFile.Yaml()
	if err != nil {
		console.IssueExit("marshalling error: ", err)
	}

	fmt.Println(string(yaml))
}
