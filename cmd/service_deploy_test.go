package cmd

import (
	"testing"

	"github.com/turnerlabs/fargate/console"
	"github.com/turnerlabs/fargate/dockercompose"
)

func TestGetDockerServiceToDeploy_Happy(t *testing.T) {

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
	compose, err := dockercompose.UnmarshalComposeYAML([]byte(yml))
	if err != nil {
		console.ErrorExit(err, "error unmarshalling docker-compose.yml")
	}

	//test
	got, _ := getDockerServiceToDeploy(&compose)

	//assert
	expected := "web"
	if got != expected {
		t.Errorf("expected: %s, got: %s", expected, got)
	}
}

func TestGetDockerServiceToDeploy_Two(t *testing.T) {

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
	compose, err := dockercompose.UnmarshalComposeYAML([]byte(yml))
	if err != nil {
		console.ErrorExit(err, "error unmarshalling docker-compose.yml")
	}

	//test
	got, _ := getDockerServiceToDeploy(&compose)

	//assert
	expected := "web"
	if got != expected {
		t.Errorf("expected: %s, got: %s", expected, got)
	}
}

func TestGetDockerServiceToDeploy_NoLabel(t *testing.T) {

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
	compose, err := dockercompose.UnmarshalComposeYAML([]byte(yml))
	if err != nil {
		console.ErrorExit(err, "error unmarshalling docker-compose.yml")
	}

	//test
	got, _ := getDockerServiceToDeploy(&compose)
	t.Log(got)

	//assert
	if got != "" {
		t.Error("expected nil service")
	}
}

func TestGetDockerServiceToDeploy_Two_37(t *testing.T) {

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
	compose, err := dockercompose.UnmarshalComposeYAML([]byte(yml))
	if err != nil {
		console.ErrorExit(err, "error unmarshalling docker-compose.yml")
	}

	//test
	got, _ := getDockerServiceToDeploy(&compose)

	//assert
	expected := "web"
	if got != expected {
		t.Errorf("expected: %s, got: %s", expected, got)
	}
}

func TestGetDockerServiceToDeploy_Two_38(t *testing.T) {

	//create a docker-compose.yml representation
	yml := `
version: "3.8"
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
	compose, err := dockercompose.UnmarshalComposeYAML([]byte(yml))
	if err != nil {
		console.ErrorExit(err, "error unmarshalling docker-compose.yml")
	}

	//test
	got, _ := getDockerServiceToDeploy(&compose)

	//assert
	expected := "web"
	if got != expected {
		t.Errorf("expected: %s, got: %s", expected, got)
	}
}
