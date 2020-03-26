package cmd

import (
	"testing"

	"github.com/turnerlabs/fargate/console"
	"github.com/turnerlabs/fargate/dockercompose"
	yaml "gopkg.in/yaml.v2"
)

func TestGetComposeServicesToDeploy_Happy(t *testing.T) {
	//create a docker-compose.yml representation
	yml := `
version: "2.0"
services:
  web:
    build: ./web
    image: 1234567890.dkr.ecr.us-east-1.amazonaws.com/my-service:0.1.0
    ports:
      - 80:8080
    environment:
      FOO: bar
`

	//unmarshal the yaml
	var compose dockercompose.DockerCompose
	err := yaml.Unmarshal([]byte(yml), &compose)
	if err != nil {
		console.ErrorExit(err, "error unmarshalling docker-compose.yml")
	}

	//test
	services := getComposeServicesToDeploy(&compose)
	names := []string{"web"}

	for i, service := range services {
		got := service.Name
		expected := names[i]

		//assert
		if got != expected {
			t.Errorf("expected: %s, got: %s", expected, got)
		}
	}
}

func TestGetComposeServicesToDeploy_Two(t *testing.T) {
	//create a docker-compose.yml representation
	yml := `
version: "3.7"
services:
  web-dev:
    build: ./web
    image: 1234567890.dkr.ecr.us-east-1.amazonaws.com/my-service:0.1.0
    ports:
      - 80:8080
    environment:
      FOO: bar
  backend-dev:
    build: ./backend
    image: 1234567890.dkr.ecr.us-east-1.amazonaws.com/my-backend:0.1.0
    environment:
      BAZ: qux
    `

	//unmarshal the yaml
	var compose dockercompose.DockerCompose
	err := yaml.Unmarshal([]byte(yml), &compose)
	if err != nil {
		console.ErrorExit(err, "error unmarshalling docker-compose.yml")
	}

	//test
	services := getComposeServicesToDeploy(&compose)
	names := []string{"web-dev", "backend-dev"}
	images := []string{
		"1234567890.dkr.ecr.us-east-1.amazonaws.com/my-service:0.1.0",
		"1234567890.dkr.ecr.us-east-1.amazonaws.com/my-backend:0.1.0",
	}

	for i, service := range services {
		got := service.Name
		expected := names[i]

		//assert
		if got != expected {
			t.Errorf("expected: %s, got: %s", expected, got)
		}

		got = service.Service.Image
		expected = images[i]

		//assert
		if got != expected {
			t.Errorf("expected: %s, got: %s", expected, got)
		}
	}
}

func TestGetComposeServicesToDeploy_IgnoreLabel(t *testing.T) {
	//create a docker-compose.yml representation
	yml := `
version: "3.7"
services:
  web:
    build: ./web
    image: 1234567890.dkr.ecr.us-east-1.amazonaws.com/my-service:0.1.0
    ports:
      - 80:8080
    environment:
      FOO: bar
  backend:
    build: ./backend
    image: 1234567890.dkr.ecr.us-east-1.amazonaws.com/my-backend:0.1.0
    environment:
      BAZ: qux
  redis:
    image: redis
    labels:
      aws.ecs.fargate.ignore: 1
    `

	//unmarshal the yaml
	var compose dockercompose.DockerCompose
	err := yaml.Unmarshal([]byte(yml), &compose)
	if err != nil {
		console.ErrorExit(err, "error unmarshalling docker-compose.yml")
	}

	//test
	services := getComposeServicesToDeploy(&compose)
	names := []string{"web", "backend"}
	images := []string{
		"1234567890.dkr.ecr.us-east-1.amazonaws.com/my-service:0.1.0",
		"1234567890.dkr.ecr.us-east-1.amazonaws.com/my-backend:0.1.0",
	}

	for i, service := range services {
		got := service.Name
		expected := names[i]

		//assert
		if got != expected {
			t.Errorf("expected: %s, got: %s", expected, got)
		}

		got = service.Service.Image
		expected = images[i]

		//assert
		if got != expected {
			t.Errorf("expected: %s, got: %s", expected, got)
		}
	}
}
