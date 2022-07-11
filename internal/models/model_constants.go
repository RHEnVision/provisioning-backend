package models

type ProviderType int

const (
	// ProviderTypeUnknown is reserved
	ProviderTypeUnknown ProviderType = iota

	// No operation (testing) provider
	ProviderTypeNoop

	// Amazon AWS provider
	ProviderTypeAWS

	// Microsoft Azure provider
	ProviderTypeAzure

	// Google Compute Engine provider
	ProviderTypeGCE
)

// AllProviders is a slice of all supported providers
var AllProviders []ProviderType

func init() {
	AllProviders = []ProviderType{
		ProviderTypeNoop,
		ProviderTypeAWS,
		ProviderTypeAzure,
		ProviderTypeGCE,
	}
}
