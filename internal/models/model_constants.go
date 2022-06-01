package models

type ProviderType int

const (
	// ProviderTypeUnknown is reserved
	ProviderTypeUnknown ProviderType = iota

	ProviderTypeAWS
	ProviderTypeAzure
	ProviderTypeGCE
)

// AllProviders is a slice of all supported providers
var AllProviders []ProviderType

func init() {
	AllProviders = []ProviderType{
		ProviderTypeAWS,
		ProviderTypeAzure,
		ProviderTypeGCE,
	}
}
