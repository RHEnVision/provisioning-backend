package providers

import "regexp"

type TypeProvider struct {
	PrintRegisteredTypes      func(string)
	PrintRegionalAvailability func(string, string)
	GenerateTypes             func() error
}

var TypeProviders = make(map[string]TypeProvider)

var ValidArchitectures = regexp.MustCompile(`^(x86[_-]64|aarch64|arm64)$`)
