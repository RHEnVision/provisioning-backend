package clients

import (
	"context"
	"embed"
	"errors"
	"fmt"

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

func MapArchitectures(_ context.Context, arch string) (ArchitectureType, error) {
	switch {
	case arch == "x86_64_mac":
		return ArchitectureTypeAppleX8664, nil
	case arch == "arm64_mac":
		return ArchitectureTypeAppleArm64, nil
	case arch == "i386":
		return ArchitectureTypeI386, nil
	case arch == "x86-64" || arch == "x86_64" || arch == "x64":
		return ArchitectureTypeX8664, nil
	case arch == "aarch64" || arch == "arm64" || arch == "Arm64" || arch == "arm":
		return ArchitectureTypeArm64, nil
	}
	return "", fmt.Errorf("%s: %w", arch, ErrArchitectureNotSupported)
}

func IsSupported(name string) bool {
	_, ok := awsSupported[name]
	return ok
}
