package dockercompose

import (
	"io/ioutil"

	"github.com/compose-spec/compose-go/loader"
	compose "github.com/compose-spec/compose-go/types"
	"github.com/turnerlabs/fargate/console"
	yaml "gopkg.in/yaml.v2"
)

//ComposeFile represents a docker-compose.yml file
//that can be manipulated
type ComposeFile struct {
	File []string
	Data DockerCompose
}

//Read loads a docker-compose.yml file
func Read(file []string) ComposeFile {
	result := ComposeFile{
		File: file,
	}
	result.Read()
	return result
}

//New returns an initialized compose file
func New(file string) ComposeFile {
	result := ComposeFile{
		File: []string{file},
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

	//Load Docker Compose yaml
	dcy, err := loadCompose(composeFile.File)
	if err != nil {
		console.ErrorExit(err, "error loading docker-compose yaml files")
	}

	console.Info(string(dcy))
	//unmarshal the yaml
	var compose DockerCompose
	err = yaml.Unmarshal(dcy, &compose)
	if err != nil {
		console.ErrorExit(err, "error unmarshalling docker-compose.yml")
	}

	composeFile.Data = compose
}

// Load and merge the docker-compose yaml files into one
func loadCompose(files []string) ([]byte, error) {
	var composeConfigFiles []compose.ConfigFile
	for _, f := range files {
		file, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, err
		}
		composeConfig, err := loader.ParseYAML(file)
		if err != nil {
			return nil, err
		}
		composeConfigFiles = append(composeConfigFiles, compose.ConfigFile{Filename: f, Config: composeConfig})

	}
	dc, err := loader.Load(compose.ConfigDetails{
		ConfigFiles: composeConfigFiles,
	})

	if err != nil {
		return nil, err
	}
	y, err := yaml.Marshal(dc)

	var remap map[string]interface{}
	err = yaml.Unmarshal(y, &remap)

	remap["version"] = composeConfigFiles[0].Config["version"]

	y, err = yaml.Marshal(remap)
	if err != nil {
		return nil, err
	}

	return y, nil

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
	err = ioutil.WriteFile(composeFile.File[0], bits, 0644)
	if err != nil {
		return err
	}
	return nil
}
