package payloads

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/go-chi/render"
)

// See clients.Source
type SourceResponse struct {
	ID           string `json:"id" yaml:"id"`
	Name         string `json:"name,omitempty" yaml:"name"`
	SourceTypeID string `json:"source_type_id" yaml:"source_type_id"`
	Uid          string `json:"uid" yaml:"uid"`
}

func (s *SourceResponse) Bind(_ *http.Request) error {
	return nil
}

func (s *SourceResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewListSourcesResponse(sourceList []*clients.Source) []render.Renderer {
	list := make([]render.Renderer, len(sourceList))
	for i, source := range sourceList {
		list[i] = &SourceResponse{
			ID:           source.ID,
			Name:         source.Name,
			SourceTypeID: source.SourceTypeID,
			Uid:          source.Uid,
		}
	}
	return list
}

type SourceUploadInfoResponse struct {
	Provider  string                       `json:"provider" yaml:"provider"`
	AwsInfo   *clients.AccountDetailsAWS   `json:"aws" nullable:"true" yaml:"aws"`
	AzureInfo *clients.AccountDetailsAzure `json:"azure" nullable:"true" yaml:"azure"`
	GcpInfo   *clients.AccountDetailsGCP   `json:"gcp" nullable:"true" yaml:"gcp"`
}

func (s SourceUploadInfoResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
