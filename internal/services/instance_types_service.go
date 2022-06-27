package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients/ec2"
	sources "github.com/RHEnVision/provisioning-backend/internal/clients/sources"
	"github.com/RHEnVision/provisioning-backend/internal/clients/sts"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/render"
)

func ListInstanceTypes(w http.ResponseWriter, r *http.Request) {
	sourceId, err := getSourceId(r)
	if err != nil {
		renderError(w, r, payloads.New3rdPartyClientError(r.Context(), "get source id from client", err))
		return
	}

	ec2Client := ec2.NewEC2Client(r.Context())

	sourcesClient, err := sources.GetSourcesClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientInitializationError(r.Context(), "cant init sources client", err))
		return
	}

	arn, err := fetchARN(r.Context(), sourcesClient, sourceId)
	if err != nil {
		renderError(w, r, payloads.New3rdPartyClientError(r.Context(), "cant fetch arn from sources", err))
		return
	}

	stsClient, err := sts.NewSTSClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientInitializationError(r.Context(), "sts client", err))
		return
	}

	crd, err := stsClient.AssumeRole(arn)
	if err != nil {
		renderError(w, r, payloads.New3rdPartyClientError(r.Context(), "assume role sts", err))
		return
	}

	newEC2Client, err := ec2Client.CreateEC2ClientFromConfig(crd)
	if err != nil {
		renderError(w, r, payloads.NewAWSError(r.Context(), "Cant create new ec2 client", err))
		return
	}

	res, err := newEC2Client.ListInstanceTypes()
	if err != nil {
		renderError(w, r, payloads.NewAWSError(r.Context(), "Cant list EC2 instance types", err))
		return
	}

	if err := render.RenderList(w, r, payloads.NewListInstanceTypeResponse(&res)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "list instance types", err))
		return
	}
}

func getSourceId(r *http.Request) (string, error) {
	sourceId := &payloads.SourceID{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read response body %w", err)
	}
	if err := json.Unmarshal(b, &sourceId); err != nil {
		return "", fmt.Errorf("unable to unmarshel response body to source id %w", err)
	}
	return sourceId.SourceId, nil
}
