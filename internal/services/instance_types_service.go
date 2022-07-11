package services

import (
	"errors"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients/ec2"
	sources "github.com/RHEnVision/provisioning-backend/internal/clients/sources"
	"github.com/RHEnVision/provisioning-backend/internal/clients/sts"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func ListInstanceTypes(w http.ResponseWriter, r *http.Request) {
	sourceId := chi.URLParam(r, "source_id")
	ec2Client := ec2.NewEC2Client(r.Context())

	sourcesClient, err := sources.GetSourcesClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientInitializationError(r.Context(), "can't init sources client", err))
		return
	}

	arn, err := fetchARN(r.Context(), sourcesClient, sourceId)
	if err != nil {
		if errors.Is(err, sources.ApplicationNotFoundErr) {
			renderError(w, r, payloads.SourcesClientError(r.Context(), "can't fetch arn from sources: application not found", err, 404))
			return
		}
		if errors.Is(err, sources.AuthenticationForSourcesNotFoundErr) {
			renderError(w, r, payloads.SourcesClientError(r.Context(), "can't fetch arn from sources: authentication not found", err, 404))
			return
		}
		renderError(w, r, payloads.SourcesClientError(r.Context(), "can't fetch arn from sources", err, 500))
		return
	}

	stsClient, err := sts.NewSTSClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientInitializationError(r.Context(), "sts client", err))
		return
	}

	crd, err := stsClient.AssumeRole(arn)
	if err != nil {
		renderError(w, r, payloads.NewAWSError(r.Context(), "assume role sts", err))
		return
	}

	newEC2Client, err := ec2Client.CreateEC2ClientFromConfig(crd)
	if err != nil {
		renderError(w, r, payloads.NewAWSError(r.Context(), "can't create new ec2 client", err))
		return
	}

	res, err := newEC2Client.ListInstanceTypes()
	if err != nil {
		renderError(w, r, payloads.NewAWSError(r.Context(), "can't list EC2 instance types", err))
		return
	}

	if err := render.RenderList(w, r, payloads.NewListInstanceTypeResponse(&res)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "list instance types", err))
		return
	}
}
