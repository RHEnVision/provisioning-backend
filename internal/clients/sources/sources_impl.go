package sources

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/config"
)

func getSourcesClient(ctx context.Context) (SourcesIntegration, error) {
	return NewClientWithResponses(*config.SourcesURL)
}

func init() {
	GetSourcesClient = getSourcesClient
}
