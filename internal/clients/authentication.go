package clients

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/rs/zerolog"
)

type Authentication struct {
	SourceApplictionID string              `json:"source_application_id"`
	ProviderType       models.ProviderType `json:"type"`
	Payload            string              `json:"payload"`
}

func NewAuthentication(str string, provType models.ProviderType) *Authentication {
	a := Authentication{
		Payload:      str,
		ProviderType: provType,
	}
	return &a
}

func NewAuthenticationFromSourceAuthType(ctx context.Context, str, authType, appID string) (*Authentication, error) {
	a := Authentication{Payload: str, SourceApplictionID: appID}
	switch authType {
	case "provisioning-arn":
		a.ProviderType = models.ProviderTypeAWS
	case "provisioning_lighthouse_subscription_id":
		a.ProviderType = models.ProviderTypeAzure
	case "provisioning_project_id":
		a.ProviderType = models.ProviderTypeGCP
	default:
		zerolog.Ctx(ctx).Warn().Msgf("Unknown auth type returned from sources: %s", authType)
		a.ProviderType = models.ProviderTypeUnknown
		return &a, ErrUnknownAuthenticationType
	}
	return &a, nil
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
		return fmt.Errorf("%w: %s", ErrUnknownAuthenticationType, auth.ProviderType.String())
	}

	return nil
}

// String returns authentication payload string (ARN, Subscription UUID, Project-ID...)
func (auth *Authentication) String() string {
	return auth.Payload
}
