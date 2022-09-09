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
