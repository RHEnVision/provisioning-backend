package clients

import (
	"context"
	"embed"
	"errors"
	"strings"

	"gopkg.in/yaml.v2"
)

var awsSupported = make(map[string]interface{})
var ErrArchitectureNotSupported = errors.New("architecture is not supported")

const x86_64 = "x86_64"
const arm64 = "arm64"
const i386 = "i386"

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

func MapArchitectures(_ context.Context, arch string) (string, error) {
	switch {
	case strings.Contains(arch, "i386"):
		return i386, nil
	case strings.Contains(arch, "86"):
		return x86_64, nil
	case strings.Contains(arch, "arm"):
		return arm64, nil
	}
	return "", ErrArchitectureNotSupported
}

func IsSupported(name string) bool {
	_, ok := awsSupported[name]
	return ok
}
