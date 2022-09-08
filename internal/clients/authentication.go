package clients

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/models"
)

type Authentication struct {
	ProviderType models.ProviderType `json:"type"`
	Payload      string              `json:"payload"`
}

func NewAuthentication(str string, provType models.ProviderType) *Authentication {
	a := Authentication{
		Payload:      str,
		ProviderType: provType,
	}
	return &a
}

func NewAuthenticationFromSourceAuthType(ctx context.Context, str, authType string) *Authentication {
	a := Authentication{Payload: str}
	switch authType {
	case "provisioning-arn":
		a.ProviderType = models.ProviderTypeAWS
	case "provisioning_lighthouse_subscription_id":
		a.ProviderType = models.ProviderTypeAzure
	case "provisioning_project_id":
		a.ProviderType = models.ProviderTypeGCP
	default:
		ctxval.Logger(ctx).Warn().Msgf("Unknown auth type returned from sources: %s", authType)
		a.ProviderType = models.ProviderTypeUnknown
	}
	return &a
}

// Type returns authentication provider type
func (auth *Authentication) Type() models.ProviderType {
	return auth.ProviderType
}

// Is checks if Authentication is of a given provider type
func (auth *Authentication) Is(providerType models.ProviderType) bool {
	return auth.ProviderType == providerType
}

// MustBe returns nil, if authentication is of given type. Otherwise, returns an error.
func (auth *Authentication) MustBe(providerType models.ProviderType) error {
	if !auth.Is(providerType) {
		return fmt.Errorf("%w: %s", UnknownAuthenticationTypeErr, auth.ProviderType.String())
	}

	return nil
}

// String returns authentication payload string (ARN, Subscription UUID, Project-ID...)
func (auth *Authentication) String() string {
	return auth.Payload
}
