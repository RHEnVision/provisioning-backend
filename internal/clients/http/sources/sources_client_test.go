package sources_test

import (
	"context"
	"fmt"
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

	t.Run("source with Provisioning Azure auth", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := io.WriteString(w, `{"data":[{"id":"256144","authtype":"provisioning_lighthouse_subscription_id","username":"1a2b3c4d-5e6f-7a8b-9c8d-7e8f9a8b7c6d","availability_status":"in_progress","resource_type":"Application","resource_id":"40340"}],"meta":{"count":1,"limit":100,"offset":0},"links":{"first":"/api/sources/v3.1/sources/350934/authentications?limit=100\u0026offset=0","last":"/api/sources/v3.1/sources/350934/authentications?limit=100\u0026offset=100"}}`)
			require.NoError(t, err, "failed to write http body for stubbed server")
		}))
		defer ts.Close()

		ctx := context.Background()
		client, err := sources.NewSourcesClientWithUrl(ctx, ts.URL)
		require.NoError(t, err, "failed to initialize sources client with test server")

		authentication, clientErr := client.GetAuthentication(ctx, "256144")
		assert.NoError(t, clientErr, "Authentication should succeed with Azure auth for Provisioning")

		assert.Equal(t, "1a2b3c4d-5e6f-7a8b-9c8d-7e8f9a8b7c6d", authentication.Payload)
	})

	t.Run("source with multiple apps", func(t *testing.T) {
		testSourceId := "256144"

		mux := http.NewServeMux()
		mux.HandleFunc(fmt.Sprintf("/sources/%s/authentications", testSourceId), func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := io.WriteString(w, `{"data":[{"id":"1","authtype":"access_key_secret_key","username":"JDOEAWSUSER","availability_status":"in_progress","resource_type":"Source","resource_id":"1"},{"id":"2","authtype":"arn","username":"arn:aws:iam::123456789999:role/redhat-cost-management-role-0f60c5c","availability_status":"in_progress","resource_type":"Application","resource_id":"1"},{"id":"3","authtype":"provisioning-arn","username":"arn:aws:iam::123456789999:role/redhat-provisioning-role-2f6d01c","availability_status":"in_progress","resource_type":"Application","resource_id":"2"},{"id":"4","authtype":"cloud-meter-arn","username":"arn:aws:iam::123456789999:role/redhat-cloud-meter-role-6331a17","availability_status":"in_progress","resource_type":"Application","resource_id":"3"}],"meta":{"count":4,"limit":100,"offset":0},"links":{"first":"/api/sources/v3.1/sources/1/authentications?limit=100\u0026offset=0","last":"/api/sources/v3.1/sources/1/authentications?limit=100\u0026offset=100"}}`)
			require.NoError(t, err, "failed to write http body for stubbed server")
		})

		ts := httptest.NewServer(mux)
		defer ts.Close()

		ctx := context.Background()
		client, err := sources.NewSourcesClientWithUrl(ctx, ts.URL)
		require.NoError(t, err, "failed to initialize sources client with test server")

		authentication, clientErr := client.GetAuthentication(ctx, testSourceId)
		assert.NoError(t, clientErr, "Authentication should succeed with Provisioning as one of many apps")

		assert.Equal(t, "arn:aws:iam::123456789999:role/redhat-provisioning-role-2f6d01c", authentication.Payload)
	})
}
