package models

import "strings"

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

func ProviderTypeFromString(str string) ProviderType {
	switch strings.ToLower(str) {
	case "noop":
		return ProviderTypeNoop
	case "aws":
		return ProviderTypeAWS
	case "azure":
		return ProviderTypeAzure
	case "gce":
		return ProviderTypeGCE
	default:
		return ProviderTypeUnknown
	}
}
