package image_builder_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	httpClients "github.com/RHEnVision/provisioning-backend/internal/clients/http"
	"github.com/RHEnVision/provisioning-backend/internal/clients/http/image_builder"
	"github.com/RHEnVision/provisioning-backend/internal/preload"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func composeStatusServer(t *testing.T) *httptest.Server {
	t.Helper()

	uploadOptions := &image_builder.UploadStatus_Options{}
	err := uploadOptions.FromAWSUploadStatus(image_builder.AWSUploadStatus{
		Ami:    "ami-1234-test",
		Region: "us-east-1",
	})
	require.NoError(t, err)
	response := image_builder.ComposeStatus{
		ImageStatus: image_builder.ImageStatus{
			Status: image_builder.ImageStatusStatusSuccess,
			UploadStatus: &image_builder.UploadStatus{
				Status:  image_builder.UploadStatusStatusSuccess,
				Type:    image_builder.UploadTypesAws,
				Options: *uploadOptions,
			},
		},
		Request: image_builder.ComposeRequest{
			Distribution: image_builder.Rhel9,
			ImageRequests: []image_builder.ImageRequest{{
				Architecture: image_builder.ImageRequestArchitectureAarch64,
			}},
		},
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		responseString, marshalErr := json.Marshal(response)
		require.NoError(t, marshalErr)
		_, err := io.WriteString(w, string(responseString))
		require.NoError(t, err, "failed to write http body for stubbed server")
	}))
}

func Test_GetAWSAmi(t *testing.T) {
	t.Run("fails to resolve AMI for mismatching architecture image", func(t *testing.T) {
		composeUUID, err := uuid.NewRandom()
		require.NoError(t, err)
		instanceType := preload.EC2InstanceType.FindInstanceType("t3.nano")
		require.NotNil(t, instanceType, "failed to find instance type")

		ts := composeStatusServer(t)
		defer ts.Close()

		ctx := context.Background()
		client, err := image_builder.NewImageBuilderClientWithUrl(ctx, ts.URL)
		require.NoError(t, err, "failed to initialize sources client with test server")

		_, amiErr := client.GetAWSAmi(ctx, composeUUID, *instanceType)

		require.ErrorIs(t, amiErr, httpClients.ErrImageArchInvalid, "Expected an architecture mismatch")
	})

	t.Run("resolves AMI correctly for matching image and instance architecture", func(t *testing.T) {
		composeUUID, err := uuid.NewRandom()
		require.NoError(t, err)
		instanceType := preload.EC2InstanceType.FindInstanceType("t4g.nano")
		require.NotNil(t, instanceType, "failed to find instance type")

		ts := composeStatusServer(t)
		defer ts.Close()

		ctx := context.Background()
		client, err := image_builder.NewImageBuilderClientWithUrl(ctx, ts.URL)
		require.NoError(t, err, "failed to initialize sources client with test server")

		ami, amiErr := client.GetAWSAmi(ctx, composeUUID, *instanceType)
		require.NoError(t, amiErr, "expected to resolve AMI correctly")
		assert.Equal(t, "ami-1234-test", ami)
	})
}
