package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func ListSources(w http.ResponseWriter, r *http.Request) {
	var sourcesList []*clients.Source

	client, err := clients.GetSourcesClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), err))
		return
	}
	provider := r.URL.Query().Get("provider")
	asProviderType := models.ProviderTypeFromString(provider)
	if asProviderType != models.ProviderTypeUnknown {
		sourcesList, err = client.ListProvisioningSourcesByProvider(r.Context(), asProviderType)
		if err != nil {
			renderError(w, r, payloads.NewClientError(r.Context(), err))
			return
		}
	} else if provider == "" {
		sourcesList, err = client.ListAllProvisioningSources(r.Context())
		if err != nil {
			renderError(w, r, payloads.NewClientError(r.Context(), err))
			return
		}
	} else {
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), fmt.Sprintf("unknown provider: %s", provider), clients.UnknownProviderErr))
		return
	}

	if err := render.RenderList(w, r, payloads.NewListSourcesResponse(sourcesList)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render sources list", err))
		return
	}
}

func GetAWSAccountIdentity(w http.ResponseWriter, r *http.Request) {
	sourceId := chi.URLParam(r, "ID")

	sourcesClient, err := clients.GetSourcesClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), err))
		return
	}

	authentication, err := sourcesClient.GetAuthentication(r.Context(), sourceId)
	if err != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), err))
		return
	}

	if typeErr := authentication.MustBe(models.ProviderTypeAWS); typeErr != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), typeErr))
		return
	}

	uploadInfo, err := getAWSUploadInfo(r.Context(), authentication)
	if err != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), err))
		return
	}

	if err := render.Render(w, r, payloads.NewAccountIdentityResponse(uploadInfo)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render account id", err))
		return
	}
}

func GetSourceUploadInfo(w http.ResponseWriter, r *http.Request) {
	sourceId := chi.URLParam(r, "ID")

	sourcesClient, err := clients.GetSourcesClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), err))
		return
	}

	authentication, err := sourcesClient.GetAuthentication(r.Context(), sourceId)
	if err != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), err))
		return
	}

	payload := payloads.SourceUploadInfoResponse{Provider: models.ProviderTypeAzure.String()}
	switch authentication.ProviderType {
	case models.ProviderTypeAWS:
		if payload.AwsInfo, err = getAWSUploadInfo(r.Context(), authentication); err != nil {
			renderError(w, r, payloads.NewAWSError(r.Context(), "unable to get AWS upload info", err))
			return
		}
	case models.ProviderTypeAzure:
		if payload.AzureInfo, err = getAzureUploadInfo(r.Context(), authentication); err != nil {
			renderError(w, r, payloads.NewAzureError(r.Context(), "unable to fetch Azure upload info", err))
			return
		}
	case models.ProviderTypeGCP, models.ProviderTypeNoop, models.ProviderTypeUnknown:
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "provider is not supported", ProviderTypeNotImplementedError))
		return
	}

	if rndrErr := render.Render(w, r, payload); rndrErr != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render Azure source details", rndrErr))
		return
	}
}

func getAWSUploadInfo(ctx context.Context, authentication *clients.Authentication) (*clients.AccountDetailsAWS, error) {
	ec2Client, err := clients.GetEC2Client(ctx, authentication, "")
	if err != nil {
		return nil, fmt.Errorf("unable to initialize AWS client: %w", err)
	}

	accountId, err := ec2Client.GetAccountId(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get account id: %w", err)
	}
	return &clients.AccountDetailsAWS{AccountID: accountId}, nil
}

func getAzureUploadInfo(ctx context.Context, authentication *clients.Authentication) (*clients.AzureSourceDetail, error) {
	azureClient, err := clients.GetAzureClient(ctx, authentication)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Azure client: %w", err)
	}

	tenantId, err := azureClient.TenantId(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch Tenant ID: %w", err)
	}
	groupList, err := azureClient.ListResourceGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch Resource Group list: %w", err)
	}

	return &clients.AzureSourceDetail{TenantID: tenantId, SubscriptionID: authentication.Payload, ResourceGroups: groupList}, nil
}
