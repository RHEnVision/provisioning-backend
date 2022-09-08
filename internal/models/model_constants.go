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
	ProviderTypeGCP
)

func ProviderTypeFromString(str string) ProviderType {
	switch strings.ToLower(str) {
	case "noop":
		return ProviderTypeNoop
	case "aws":
		return ProviderTypeAWS
	case "azure":
		return ProviderTypeAzure
	case "gcp":
		return ProviderTypeGCP
	default:
		return ProviderTypeUnknown
	}
}

func (pt ProviderType) String() string {
	switch pt {
	case ProviderTypeNoop:
		return "noop"
	case ProviderTypeAWS:
		return "aws"
	case ProviderTypeAzure:
		return "azure"
	case ProviderTypeGCP:
		return "gcp"
	case ProviderTypeUnknown:
	default:
		return ""
	}
	return ""
}
