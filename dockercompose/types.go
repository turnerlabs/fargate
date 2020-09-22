package dockercompose

import (
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// DockerCompose represents a docker-compose.yml file
type DockerCompose struct {
	Version  string              `yaml:"version"`
	Services map[string]*Service `yaml:"services"`
}

// Service represents a docker container
type Service struct {
	Image       string            `yaml:"image,omitempty"`
	Ports       []Port            `yaml:"ports,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Secrets     map[string]string `yaml:"x-fargate-secrets,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
}

// Port represents a port
type Port struct {
	Published int64 `yaml:"published"`
	Target    int64 `yaml:"target"`
}

// used to parse the short syntax
type dockerComposeShortPortSyntax struct {
	Version  string                             `yaml:"version"`
	Services map[string]*serviceShortPortSyntax `yaml:"services"`
}

type serviceShortPortSyntax struct {
	Image       string            `yaml:"image,omitempty"`
	Ports       []string          `yaml:"ports,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Secrets     map[string]string `yaml:"x-fargate-secrets,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
}

//UnmarshalComposeYAML unmarshals yaml into a DockerCompose struct
//handles versioning and schema issues
func UnmarshalComposeYAML(yamlBytes []byte) (DockerCompose, error) {
	var result DockerCompose

	//does the yaml use the long or short port syntax?
	//note: docker-compose config used to only output the long syntax after version 3.2
	//but they've since changed to only output the short syntax
	//so we need to support both versions depending on
	//what version of docker-compose the user has installed

	var short dockerComposeShortPortSyntax
	err := yaml.Unmarshal(yamlBytes, &short)
	if err == nil {

		//copy data from short types to result types
		result.Version = short.Version
		result.Services = make(map[string]*Service, len(short.Services))
		for s, svc := range short.Services {

			//convert ports
			ports := []Port{}
			for _, p := range svc.Ports {
				portString := strings.Split(p, ":")
				published, err := strconv.ParseInt(portString[0], 10, 64)
				if err != nil {
					return result, err
				}
				target, err := strconv.ParseInt(strings.Split(portString[1], "/")[0], 10, 64)
				if err != nil {
					return result, err
				}
				ports = append(ports, Port{
					Published: published,
					Target:    target,
				})
			}

			result.Services[s] = &Service{
				Image:       svc.Image,
				Environment: svc.Environment,
				Secrets:     svc.Secrets,
				Labels:      svc.Labels,
				Ports:       ports,
			}
		}
	} else { //error unmarshaling short syntax

		//try long syntax
		err := yaml.Unmarshal(yamlBytes, &result)
		if err != nil {
			return result, nil
		}
	}

	return result, nil
}
