package sources_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/clients/http/sources"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSourcesClient_GetAuthentication(t *testing.T) {
	t.Run("source with missing Provisioning auth", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := io.WriteString(w, `{"data":[{"id":"256144","authtype":"provisioning-arn","username":"arn:aws:asdfasdfsdfsdafsdf","availability_status":"in_progress","resource_type":"Source","resource_id":"304935"}],"meta":{"count":1,"limit":100,"offset":0},"links":{"first":"/api/sources/v3.1/sources/304935/authentications?limit=100\u0026offset=0","last":"/api/sources/v3.1/sources/304935/authentications?limit=100\u0026offset=100"}}`)
			require.NoError(t, err, "failed to write http body for stubbed server")
		}))
		defer ts.Close()

		ctx := context.Background()
		client, err := sources.NewSourcesClientWithUrl(ctx, ts.URL)
		require.NoError(t, err, "failed to initialize sources client with test server")

		_, err = client.GetAuthentication(ctx, "256144")
		assert.Error(t, err, "Authentication should not succeed with missing link to Provisioning service")
	})
}
