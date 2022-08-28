package services

import (
	"errors"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/clients/ec2"
	sources "github.com/RHEnVision/provisioning-backend/internal/clients/sources"
	"github.com/RHEnVision/provisioning-backend/internal/clients/sts"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func ListInstanceTypes(w http.ResponseWriter, r *http.Request) {
	sourceId := chi.URLParam(r, "ID")
	ec2Client := ec2.NewEC2Client(r.Context())
	logger := ctxval.Logger(r.Context())

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

	res, err := newEC2Client.ListInstanceTypesWithPaginator()
	if err != nil {
		renderError(w, r, payloads.NewAWSError(r.Context(), "can't list EC2 instance types", err))
		return
	}

	numBefore := len(res)
	instances, err := ec2.NewInstanceTypes(r.Context(), res)
	logger.Trace().Msgf("Number of AWS EC2 instance types %d, Number after filtering: %d", numBefore, len(*instances))
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
