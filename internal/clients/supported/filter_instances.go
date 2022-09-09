package supported

import (
	"embed"
	"errors"

	"gopkg.in/yaml.v2"
)

var awsSupported = make(map[string]interface{})
var ErrArchitectureNotSupported = errors.New("architecture is not supported")

type supportedInstanceTypes struct {
	Aws []string `yaml:"aws"`
}

//go:embed supported_instance_types.yml
var f embed.FS

func init() {
	var supported supportedInstanceTypes
	it, err := f.ReadFile("supported_instance_types.yml")
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(it, &supported); err != nil {
		panic(err)
	}

	for _, i := range supported.Aws {
		awsSupported[i] = nil
	}
}

func IsSupported(name string) bool {
	_, ok := awsSupported[name]
	return ok
}
