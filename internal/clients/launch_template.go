package clients

// LaunchTemplate represents a generic launch template for a hyperscaler.
type LaunchTemplate struct {
	// ID is an identifier, for example "lt-94397398248932342" for AWS EC2.
	ID string `json:"id" yaml:"id"`

	// Name describes the launch template, user defined.
	Name string `json:"name" yaml:"name"`
}
