package clients

type AzureInstanceID string

// InstanceDescription defines a model for an instance description
type InstanceDescription struct {
	// The id of the instance
	ID string `json:"id,omitempty" yaml:"id"`

	// The public IPv4 dns of the instance or empty when not available
	DNS string `json:"dns,omitempty" yaml:"dns"`

	// The public IPv4 of the instance or empty when not available
	IPv4 string `json:"ipv4,omitempty" yaml:"ipv4"`

	// The IPv4 of the instance or empty when not available
	PrivateIPv4 string `json:"private_ipv4,omitempty" yaml:"private_ipv4"`

	// The IPv6 of the instance or empty when not available
	PrivateIPv6 string `json:"private_ipv6,omitempty" yaml:"private_ipv6"`
}
