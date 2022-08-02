package sources

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/config"
)

func getSourcesAPIClient(ctx context.Context) (APIClient, error) {
	return NewClientWithResponses(config.Sources.URL)
}

func init() {
	GetSourcesAPIClient = getSourcesAPIClient
}
