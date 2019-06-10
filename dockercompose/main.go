package dockercompose

import (
	"bytes"
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
func Read(file string) ComposeFile {
	result := ComposeFile{
		File: file,
	}
	result.Read()
	return result
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
func (composeFile *ComposeFile) Read() {
	console.Debug("running docker-compose config [%s]", composeFile.File)
	cmd := exec.Command("docker-compose", "-f", composeFile.File, "config")

	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	if err := cmd.Start(); err != nil {
		console.ErrorExit(err, errbuf.String())
	}

	if err := cmd.Wait(); err != nil {
		console.IssueExit(errbuf.String())
	}

	//unmarshal the yaml
	var compose DockerCompose
	err := yaml.Unmarshal(outbuf.Bytes(), &compose)
	if err != nil {
		console.ErrorExit(err, "error unmarshalling docker-compose.yml")
	}

	composeFile.Data = compose
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
