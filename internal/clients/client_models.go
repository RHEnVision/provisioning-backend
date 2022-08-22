package clients

type ID = string

// Source defines model for Source.
type Source struct {
	// ID of the resource
	Id *ID `json:"id,omitempty"`

	// The name of the source
	Name *string `json:"name,omitempty"`

	// ID of the resource
	SourceTypeId *ID `json:"source_type_id,omitempty"`

	// Unique ID of the inventory source installation
	Uid *string `json:"uid,omitempty"`
}

// An instance type defines a model for an instance type that corresponds to one in a cloud provider.
type InstanceType struct {
	// The name of the instance type
	Name string `json:"name,omitempty"`

	// Virtual CPU (maps to hypervisor hyper-thread)
	VCPUs *int32 `json:"vcpus,omitempty"`

	// Core (physical or virtual core)
	Cores *int32 `json:"cores,omitempty"`

	// The size of the memory, in MiB.
	MemoryMiB int64 `json:"memory,omitempty"`

	// Does the instance type supports RHEL
	Supported bool `json:"supported,omitempty"`

	// Instance type's Architecture: i386, arm64, x86_64
	Architecture string `json:"architecture,omitempty"`
}
