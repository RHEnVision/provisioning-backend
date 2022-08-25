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

type ArchitectureType string

const (
	ArchitectureTypeI386       ArchitectureType = "i386"
	ArchitectureTypeX8664      ArchitectureType = "x86_64"
	ArchitectureTypeArm64      ArchitectureType = "arm64"
	ArchitectureTypeAppleX8664 ArchitectureType = "apple-x86_64"
	ArchitectureTypeAppleArm64 ArchitectureType = "apple-arm64"
)

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
	case strings.Contains(arch, "86") && strings.Contains(arch, "mac"):
		return ArchitectureTypeAppleX8664, nil
	case strings.Contains(arch, "arm") && strings.Contains(arch, "mac"):
		return ArchitectureTypeAppleArm64, nil
	case strings.Contains(arch, "i386"):
		return ArchitectureTypeI386, nil
	case strings.Contains(arch, "86"):
		return ArchitectureTypeX8664, nil
	case strings.Contains(arch, "arm"):
		return ArchitectureTypeArm64, nil
	}
	return "", ErrArchitectureNotSupported
}

func IsSupported(name string) bool {
	_, ok := awsSupported[name]
	return ok
}
