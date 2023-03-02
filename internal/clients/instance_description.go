package clients

// InstanceInfo defines a model for an instance description
type InstanceDescription struct {
	// The id of the instance
	ID string `json:"id,omitempty" yaml:"id"`

	// The public ipv4 dns of the instance
	PublicDNS string `json:"dns,omitempty" yaml:"dns"`

	// the public ipv4 of the instance
	PublicIPv4 string `json:"ipv4,omitempty" yaml:"ipv4"`
}
