package services

import (
	"errors"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/clients/http/ec2"
	sources "github.com/RHEnVision/provisioning-backend/internal/clients/http/sources"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func ListInstanceTypes(w http.ResponseWriter, r *http.Request) {
	logger := ctxval.Logger(r.Context())

	sourceId := chi.URLParam(r, "ID")
	region := r.URL.Query().Get("region")
	if region == "" {
		renderError(w, r, payloads.NewNotFoundError(r.Context(), ec2.RegionNotFoundErr))
	}

	sourcesClient, err := clients.GetSourcesClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientInitializationError(r.Context(), "can't init sources client", err))
		return
	}

	arn, err := sourcesClient.GetArn(r.Context(), sourceId)
	if err != nil {
		if errors.Is(err, sources.ApplicationNotFoundErr) {
			renderError(w, r, payloads.ClientError(r.Context(), "Sources", "can't fetch arn from sources: application not found", err, 404))
			return
		}
		if errors.Is(err, sources.AuthenticationForSourcesNotFoundErr) {
			renderError(w, r, payloads.ClientError(r.Context(), "Sources", "can't fetch arn from sources: authentication not found", err, 404))
			return
		}
		renderError(w, r, payloads.ClientError(r.Context(), "Sources", "can't fetch arn from sources", err, 500))
		return
	}

	ec2Client, err := clients.GetCustomerEC2Client(r.Context(), arn, region)
	if err != nil {
		renderError(w, r, payloads.NewAWSError(r.Context(), "failed to establish ec2 connection", err))
		return
	}

	res, err := ec2Client.ListInstanceTypesWithPaginator()
	if err != nil {
		renderError(w, r, payloads.NewAWSError(r.Context(), "can't list EC2 instance types", err))
		return
	}

	numBefore := len(res)
	instances, err := ec2.NewInstanceTypes(r.Context(), res)
	logger.Trace().Msgf("Total AWS EC2 instance types: %d (%d after architecture breakdown)", numBefore, len(*instances))
	if err != nil {
		renderError(w, r, payloads.NewAWSError(r.Context(), "can't convertAWSTypes", err))
		return
	}

	if err != nil {
		renderError(w, r, payloads.NewAWSError(r.Context(), "can't FilterUnsupportedTypes", err))
		return
	}
	if err := render.RenderList(w, r, payloads.NewListInstanceTypeResponse(instances)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "list instance types", err))
		return
	}
}
