package clients

import (
	"github.com/RHEnVision/provisioning-backend/internal/models"
)

// Source defines model for Source. Maps 1:1 to Source Database.
type Source struct {
	// ID of the resource
	ID string

	// The name of the source
	Name string

	// Source Type ID (number assigned to AWS source or Azure source)
	SourceTypeID string

	// UUID of the inventory source installation
	Uid string

	// Provider name as string
	Provider models.ProviderType

	// Status of the source
	Status string
}
