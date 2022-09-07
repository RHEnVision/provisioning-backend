package providers

type TypeProvider struct {
	PrintRegisteredTypes      func(string)
	PrintRegionalAvailability func(string, string)
	GenerateTypes             func() error
}

var TypeProviders = make(map[string]TypeProvider)
