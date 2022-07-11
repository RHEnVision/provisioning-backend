package sources

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
)

var appTypeMutex sync.Mutex

type AppType struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type dataElement struct {
	Data []AppType `json:"data"`
}

func getSourcesClient(ctx context.Context) (SourcesIntegration, error) {
	client, err := NewClientWithResponses(config.Sources.URL)
	if err != nil {
		return client, err
	}
	// synchronize on the AppTypeId check
	appTypeMutex.Lock()
	defer appTypeMutex.Unlock()

	if config.Sources.AppTypeId == "" {
		ctxval.GetLogger(ctx).Info().Msg("Fetching the Application Type ID of Provisioning for Sources")
		resp, err := client.ListApplicationTypes(ctx, &ListApplicationTypesParams{})
		if err != nil {
			ctxval.GetLogger(ctx).Warn().Err(err).Msg("Failed to fetch ApplicationTypes from sources")
			return client, fmt.Errorf("Failed to fetch ApplicationTypes: %w", err)
		}
		defer resp.Body.Close()
		var appTypesData dataElement
		if err = json.NewDecoder(resp.Body).Decode(&appTypesData); err != nil {
			return client, fmt.Errorf("Could not unmarshal the response: %w", err)
		}
		for _, t := range appTypesData.Data {
			if t.Name == "/insights/platform/provisioning" {
				config.SetSourcesAppTypeId(t.Id)
				ctxval.GetLogger(ctx).Info().Msgf("The Application Type ID was set to %s", t.Id)
				break
			}
		}
	}
	return client, nil
}

func init() {
	GetSourcesClient = getSourcesClient
}
