package dockercompose

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/turnerlabs/fargate/console"
	yaml "gopkg.in/yaml.v2"
)

//ComposeFile represents a docker-compose.yml file
//that can be manipulated
type ComposeFile struct {
	File string
	Data DockerCompose
}

//Read loads a docker-compose.yml file
func Read(file string) (ComposeFile, error) {
	result := ComposeFile{
		File: file,
	}
	var err error
	err = result.Read()
	return result, err
}

//New returns an initialized compose file
func New(file string) ComposeFile {
	result := ComposeFile{
		File: file,
		Data: DockerCompose{
			Version:  "3.7",
			Services: make(map[string]*Service),
		},
	}
	return result
}

//Read reads the data structure from the file
//note that all variable interpolations are fully rendered
func (composeFile *ComposeFile) Read() error {
	console.Debug("running docker-compose config [%s]", composeFile.File)
	cmd := exec.Command("docker-compose", "-f", composeFile.File, "config")

	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("%v: %w", errbuf.String(), err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("%v: %w", errbuf.String(), err)
	}

	//unmarshal the yaml

	//is the yaml using the long port syntax?
	compose, err := UnmarshalComposeYAML(outbuf.Bytes())
	if err != nil {
		return fmt.Errorf("unmarshalling docker compose yaml: %w", err)
	}
	if compose.Version == "" || len(compose.Services) == 0 {
		return errors.New("unable to parse compose file")
	}
	composeFile.Data = compose
	return nil
}

//AddService adds a service to a compose file
func (composeFile *ComposeFile) AddService(name string) *Service {
	result := &Service{}
	result.Environment = make(map[string]string)
	result.Labels = make(map[string]string)
	result.Secrets = make(map[string]string)
	result.Ports = []Port{}
	composeFile.Data.Services[name] = result
	return result
}

//Yaml returns the yaml for this compose file
func (composeFile *ComposeFile) Yaml() ([]byte, error) {
	return yaml.Marshal(composeFile.Data)
}

//Write writes the data to a file
func (composeFile *ComposeFile) Write() error {
	bits, err := composeFile.Yaml()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(composeFile.File, bits, 0644)
	if err != nil {
		return err
	}
	return nil
}
