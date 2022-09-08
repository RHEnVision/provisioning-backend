package clients

// Region represents a provider's region (e.g. 'us-east-1' for EC2 or 'eastus' for Azure)
type Region string

func (r Region) String() string {
	return string(r)
}

// Zone represents a provider's zone. There are multiple types of zones (regional, wireless, cities)
// based on the provider. This type does not make any difference, as long as they have unique names.
// The name must include region in the name, so it is unique for each provider.
type Zone string

func (z Zone) String() string {
	return string(z)
}
