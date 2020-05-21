package cmd

import (
	"testing"

	"github.com/turnerlabs/fargate/console"
	"github.com/turnerlabs/fargate/dockercompose"
	yaml "gopkg.in/yaml.v2"
)

func TestGetDockerServicesFromComposeFile_Happy(t *testing.T) {
	//create a docker-compose.yml representation
	yml := `
version: "2"
services:
  web:
    build: .
    image: 1234567890.dkr.ecr.us-east-1.amazonaws.com/my-service:0.1.0
    ports:
    - "80:8080"
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
	got, _ := getDockerServicesFromComposeFile(&compose, false)
	expected := map[string]*dockercompose.Service{
		"web": &dockercompose.Service{
			Image: "1234567890.dkr.ecr.us-east-1.amazonaws.com/my-service:0.1.0",
		},
	}

	if len(expected) != len(got) {
		t.Errorf("expected: %d, got: %d", len(expected), len(got))
	}

	for name, service := range expected {
		actual := got[name]

		if actual == nil {
			t.Errorf("expected: %s, got nil", name)
		}

		if actual.Image != service.Image {
			t.Errorf("expected: %s, got %s", service.Image, actual.Image)
		}
	}
}

func TestGetDockerServicesFromComposeFile_Two(t *testing.T) {
	//create a docker-compose.yml representation
	yml := `
version: "2"
services:
  web:
    build: .
    image: 1234567890.dkr.ecr.us-east-1.amazonaws.com/my-service:0.1.0
    ports:
    - "80:8080"
    environment:
      FOO: bar
    labels:
      aws.ecs.fargate.deploy: 1
  redis:
    image: redis
`

	//unmarshal the yaml
	var compose dockercompose.DockerCompose
	err := yaml.Unmarshal([]byte(yml), &compose)
	if err != nil {
		console.ErrorExit(err, "error unmarshalling docker-compose.yml")
	}

	//test
	got, _ := getDockerServicesFromComposeFile(&compose, false)
	expected := map[string]*dockercompose.Service{
		"web": &dockercompose.Service{
			Image: "1234567890.dkr.ecr.us-east-1.amazonaws.com/my-service:0.1.0",
		},
	}

	if len(expected) != len(got) {
		t.Errorf("expected: %d, got: %d", len(expected), len(got))
	}

	for name, service := range expected {
		actual := got[name]

		if actual == nil {
			t.Errorf("expected: %s, got nil", name)
		}

		if actual.Image != service.Image {
			t.Errorf("expected: %s, got %s", service.Image, actual.Image)
		}
	}
}

func TestGetDockerServicesFromComposeFile_NoLabel(t *testing.T) {
	//create a docker-compose.yml representation
	yml := `
version: "2"
services:
  web:
    build: .
    image: 1234567890.dkr.ecr.us-east-1.amazonaws.com/my-service:0.1.0
    ports:
    - "80:8080"
    environment:
      FOO: bar
  redis:
    image: redis
`

	//unmarshal the yaml
	var compose dockercompose.DockerCompose
	err := yaml.Unmarshal([]byte(yml), &compose)
	if err != nil {
		console.ErrorExit(err, "error unmarshalling docker-compose.yml")
	}

	//test
	expected := "Please indicate which docker container you'd like to deploy using the label \"aws.ecs.fargate.deploy: 1\""
	_, got := getDockerServicesFromComposeFile(&compose, false)

	if got.Error() != expected {
		t.Errorf("expected error: %v, got error: %v", expected, got)
	}
}

func TestGetDockerServicesFromComposeFile_Two_37(t *testing.T) {

	//create a docker-compose.yml representation
	yml := `
version: "3.7"
services:
  web:
    build: .
    image: 1234567890.dkr.ecr.us-east-1.amazonaws.com/my-service:0.1.0
    ports:
    - published: 80
    - target: 8080
    environment:
      FOO: bar
    labels:
      aws.ecs.fargate.deploy: 1
  redis:
    image: redis
`

	//unmarshal the yaml
	var compose dockercompose.DockerCompose
	err := yaml.Unmarshal([]byte(yml), &compose)
	if err != nil {
		console.ErrorExit(err, "error unmarshalling docker-compose.yml")
	}

	//test
	got, _ := getDockerServicesFromComposeFile(&compose, false)
	expected := map[string]*dockercompose.Service{
		"web": &dockercompose.Service{
			Image: "1234567890.dkr.ecr.us-east-1.amazonaws.com/my-service:0.1.0",
		},
	}

	if len(expected) != len(got) {
		t.Errorf("expected: %d, got: %d", len(expected), len(got))
	}

	for name, service := range expected {
		actual := got[name]

		if actual == nil {
			t.Errorf("expected: %s, got nil", name)
		}

		if actual.Image != service.Image {
			t.Errorf("expected: %s, got %s", service.Image, actual.Image)
		}
	}
}

func TestGetDockerServicesFromComposeFile_ComposeAllHappy(t *testing.T) {
	//create a docker-compose.yml representation
	yml := `
version: "2.0"
services:
  app:
    build: .
    image: 1234567890.dkr.ecr.us-east-1.amazonaws.com/my-service:0.1.0
    ports:
    - "80:8080"
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
	got, _ := getDockerServicesFromComposeFile(&compose, true)
	expected := map[string]*dockercompose.Service{
		"app": &dockercompose.Service{
			Image: "1234567890.dkr.ecr.us-east-1.amazonaws.com/my-service:0.1.0",
		},
	}

	if len(expected) != len(got) {
		t.Errorf("expected: %d, got: %d", len(expected), len(got))
	}

	for name, service := range expected {
		actual := got[name]

		if actual == nil {
			t.Errorf("expected: %s, got nil", name)
		}

		if actual.Image != service.Image {
			t.Errorf("expected: %s, got %s", service.Image, actual.Image)
		}
	}
}

func TestGetDockerServicesFromComposeFile_ComposeAllTwo(t *testing.T) {
	//create a docker-compose.yml representation
	yml := `
version: "3.7"
services:
  app:
    image: 1234567890.dkr.ecr.us-east-1.amazonaws.com/my-app:0.1.0
    environment:
        FOO: bar
  api:
    image: 1234567890.dkr.ecr.us-east-1.amazonaws.com/my-api:0.1.0
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
	got, _ := getDockerServicesFromComposeFile(&compose, true)
	expected := map[string]*dockercompose.Service{
		"app": &dockercompose.Service{
			Image: "1234567890.dkr.ecr.us-east-1.amazonaws.com/my-app:0.1.0",
		},
		"api": &dockercompose.Service{
			Image: "1234567890.dkr.ecr.us-east-1.amazonaws.com/my-api:0.1.0",
		},
	}

	if len(expected) != len(got) {
		t.Errorf("expected: %d, got: %d", len(expected), len(got))
	}

	for name, service := range expected {
		actual := got[name]

		if actual == nil {
			t.Errorf("expected: %s, got nil", name)
		}

		if actual.Image != service.Image {
			t.Errorf("expected: %s, got %s", service.Image, actual.Image)
		}
	}
}

func TestGetDockerServicesFromComposeFile_ComposeAllWithIgnoreLabel(t *testing.T) {
	//create a docker-compose.yml representation
	yml := `
version: "3.7"
services:
  web:
    image: 1234567890.dkr.ecr.us-east-1.amazonaws.com/my-service:0.1.0
    environment:
      FOO: bar
  backend:
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
	got, _ := getDockerServicesFromComposeFile(&compose, true)
	expected := map[string]*dockercompose.Service{
		"web": &dockercompose.Service{
			Image: "1234567890.dkr.ecr.us-east-1.amazonaws.com/my-service:0.1.0",
		},
		"backend": &dockercompose.Service{
			Image: "1234567890.dkr.ecr.us-east-1.amazonaws.com/my-backend:0.1.0",
		},
	}

	if len(expected) != len(got) {
		t.Errorf("expected: %d, got: %d", len(expected), len(got))
	}

	for name, service := range expected {
		actual := got[name]

		if actual == nil {
			t.Errorf("expected: %s, got nil", name)
		}

		if actual.Image != service.Image {
			t.Errorf("expected: %s, got %s", service.Image, actual.Image)
		}
	}
}

func TestGetDockerServicesFromComposeFile_ComposeAllNoLabel(t *testing.T) {
	//create a docker-compose.yml representation
	yml := `
version: "3.7"
services:
  web:
    image: 1234567890.dkr.ecr.us-east-1.amazonaws.com/my-service:0.1.0
    environment:
      FOO: bar
    labels:
      aws.ecs.fargate.ignore: 1
  backend:
    image: 1234567890.dkr.ecr.us-east-1.amazonaws.com/my-backend:0.1.0
    environment:
      BAZ: qux
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
	expected := "Please indicate at least one docker container you'd like to deploy"
	_, got := getDockerServicesFromComposeFile(&compose, true)

	if got.Error() != expected {
		t.Errorf("expected error: %v, got error: %v", expected, got)
	}
}
