package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/payloads/validation"

	"github.com/RHEnVision/provisioning-backend/internal/cache"
	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/page"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func ListSources(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("provider")
	asProviderType := models.ProviderTypeFromString(provider)

	if provider == "" {
		ListAllProvisioningSources(w, r)
		return
	}
	switch asProviderType {
	case models.ProviderTypeAWS, models.ProviderTypeGCP, models.ProviderTypeAzure:
		ListProvisioningSourcesByProvider(w, r, asProviderType)
		return
	default:
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), fmt.Sprintf("unknown provider: %s", provider), clients.ErrUnknownProvider))
		return
	}
}

func ListAllProvisioningSources(w http.ResponseWriter, r *http.Request) {
	var sourcesList []*clients.Source

	client, err := clients.GetSourcesClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), err))
		return
	}

	sourcesList, total, err := client.ListAllProvisioningSources(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), err))
		return
	}

	info := page.NewOffsetMetadata(r.Context(), r, total)

	if err := render.Render(w, r, payloads.NewListSourcesResponse(sourcesList, info)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render sources list", err))
		return
	}
}

func ListProvisioningSourcesByProvider(w http.ResponseWriter, r *http.Request, asProviderType models.ProviderType) {
	var sourcesList []*clients.Source

	client, err := clients.GetSourcesClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), err))
		return
	}

	sourcesList, total, err := client.ListProvisioningSourcesByProvider(r.Context(), asProviderType)
	if err != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), err))
		return
	}

	info := page.NewOffsetMetadata(r.Context(), r, total)

	if err := render.Render(w, r, payloads.NewListSourcesResponse(sourcesList, info)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render sources list", err))
		return
	}
}

func GetSourceUploadInfo(w http.ResponseWriter, r *http.Request) {
	sourceId := chi.URLParam(r, "ID")
	if err := validation.DigitsOnly(sourceId); err != nil {
		renderError(w, r, payloads.NewURLParsingError(r.Context(), "id parameter invalid", err))
	}

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

	payload := payloads.SourceUploadInfoResponse{Provider: authentication.ProviderType.String()}
	switch authentication.ProviderType {
	case models.ProviderTypeAWS:
		if payload.AwsInfo, err = getAWSAccountDetails(r.Context(), sourceId, authentication); err != nil {
			renderError(w, r, payloads.NewAWSError(r.Context(), "unable to get AWS upload info", err))
			return
		}
	case models.ProviderTypeAzure:
		if payload.AzureInfo, err = getAzureAccountDetails(r.Context(), sourceId, authentication); err != nil {
			renderError(w, r, payloads.NewAzureError(r.Context(), "unable to fetch Azure upload info", err))
			return
		}
	case models.ProviderTypeGCP:
		if payload.GcpInfo, err = getGCPAccountDetails(r.Context(), sourceId, authentication); err != nil {
			renderError(w, r, payloads.NewGCPError(r.Context(), "unable to get GCP upload info", err))
			return
		}
	case models.ProviderTypeNoop, models.ProviderTypeUnknown:
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "provider is not supported", ErrProviderTypeNotImplemented))
		return
	}

	if rndrErr := render.Render(w, r, payload); rndrErr != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render Azure source details", rndrErr))
		return
	}
}

func getAWSAccountDetails(ctx context.Context, sourceId string, authentication *clients.Authentication) (*clients.AccountDetailsAWS, error) {
	result := &clients.AccountDetailsAWS{}

	err := cache.Find(ctx, sourceId, result)
	if errors.Is(err, cache.ErrNotFound) {
		ec2Client, clientErr := clients.GetEC2Client(ctx, authentication, "")
		if clientErr != nil {
			return nil, fmt.Errorf("unable to initialize AWS client: %w", clientErr)
		}

		result.AccountID, clientErr = ec2Client.GetAccountId(ctx)
		if clientErr != nil {
			return nil, fmt.Errorf("unable to get account id: %w", clientErr)
		}

		clientErr = cache.SetForever(ctx, sourceId, result)
		if clientErr != nil {
			return nil, fmt.Errorf("cache set error: %w", clientErr)
		}
	} else if err != nil {
		return nil, fmt.Errorf("cache find error: %w", err)
	}

	return result, nil
}

func getAzureAccountDetails(ctx context.Context, sourceId string, authentication *clients.Authentication) (*clients.AccountDetailsAzure, error) {
	var tenantId clients.AzureTenantId

	azureClient, err := clients.GetAzureClient(ctx, authentication)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Azure client: %w", err)
	}

	err = cache.Find(ctx, sourceId, &tenantId)
	if errors.Is(err, cache.ErrNotFound) {
		tenantId, err = azureClient.TenantId(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to fetch Tenant ID: %w", err)
		}

		err = cache.SetForever(ctx, sourceId, &tenantId)
		if err != nil {
			return nil, fmt.Errorf("cache set error: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("cache find error: %w", err)
	}

	groupList, err := azureClient.ListResourceGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch Resource Group list: %w", err)
	}

	return &clients.AccountDetailsAzure{
		TenantID:       tenantId,
		SubscriptionID: authentication.Payload,
		ResourceGroups: groupList,
	}, nil
}

func getGCPAccountDetails(ctx context.Context, sourceId string, authentication *clients.Authentication) (*clients.AccountDetailsGCP, error) {
	return &clients.AccountDetailsGCP{}, nil
}
