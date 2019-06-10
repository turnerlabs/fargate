package dockercompose

import (
	"strconv"
	"strings"
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
// long syntax introduced in v3.2
type Port struct {
	Published int64 `yaml:"published"`
	Target    int64 `yaml:"target"`
}

type dockerComposeVersion struct {
	Version string `yaml:"version"`
}

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

type dockerComposeLongPortSyntax struct {
	Version  string                            `yaml:"version"`
	Services map[string]*serviceLongPortSyntax `yaml:"services"`
}

type serviceLongPortSyntax struct {
	Image       string            `yaml:"image,omitempty"`
	Ports       []Port            `yaml:"ports,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Secrets     map[string]string `yaml:"x-fargate-secrets,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
}

//UnmarshalYAML is a custom unmarshaller
//to support different versions of docker compose
func (dc *DockerCompose) UnmarshalYAML(unmarshal func(interface{}) error) error {

	//unmarshal into dockerComposeVersion to detect version
	var v dockerComposeVersion
	err := unmarshal(&v)
	if err != nil {
		return err
	}

	//convert version to a number
	version, err := strconv.ParseFloat(v.Version, 32)
	if err != nil {
		return err
	}

	//unmarshal into a 3.7 compatible version
	dc.Version = "3.7"

	// config command for >= v3.2 will always return the long port syntax
	if version < 3.2 {

		//unmarshal to short port syntax format
		var short dockerComposeShortPortSyntax
		err := unmarshal(&short)
		if err != nil {
			return err
		}

		//copy data
		dc.Services = make(map[string]*Service, len(short.Services))
		for s, svc := range short.Services {

			//convert ports
			ports := []Port{}
			for _, p := range svc.Ports {
				portString := strings.Split(p, ":")
				published, err := strconv.ParseInt(portString[0], 10, 64)
				if err != nil {
					return err
				}
				target, err := strconv.ParseInt(strings.Split(portString[1], "/")[0], 10, 64)
				if err != nil {
					return err
				}
				ports = append(ports, Port{
					Published: published,
					Target:    target,
				})
			}

			dc.Services[s] = &Service{
				Image:       svc.Image,
				Environment: svc.Environment,
				Secrets:     svc.Secrets,
				Labels:      svc.Labels,
				Ports:       ports,
			}
		}
	} else { //long port syntax
		var long dockerComposeLongPortSyntax
		err := unmarshal(&long)
		if err != nil {
			return nil
		}

		//copy data
		dc.Services = make(map[string]*Service, len(long.Services))
		for s, svc := range long.Services {
			dc.Services[s] = &Service{
				Image:       svc.Image,
				Environment: svc.Environment,
				Secrets:     svc.Secrets,
				Labels:      svc.Labels,
				Ports:       svc.Ports,
			}
		}
	}

	return nil
}
