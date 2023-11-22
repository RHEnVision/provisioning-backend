package payloads

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/page"
	"github.com/go-chi/render"
)

// See clients.Source
type SourceResponse struct {
	ID           string `json:"id" yaml:"id"`
	Name         string `json:"name,omitempty" yaml:"name"`
	Uid          string `json:"uid" yaml:"uid"`
	Provider     string `json:"provider" yaml:"provider" description:"One of ('azure', 'aws', 'gcp')"`
	Status       string `json:"status" yaml:"status"`
	SourceTypeID string `json:"source_type_id" yaml:"source_type_id" deprecated:"true"`
}

type SourceListResponse struct {
	Data     []*SourceResponse `json:"data" yaml:"data"`
	Metadata page.Metadata     `json:"metadata" yaml:"metadata"`
}

func (s *SourceResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (s *SourceListResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewListSourcesResponse(sourceList []*clients.Source, meta *page.Metadata) render.Renderer {
	list := make([]*SourceResponse, len(sourceList))
	for i, source := range sourceList {
		list[i] = &SourceResponse{
			ID:           source.ID,
			Name:         source.Name,
			SourceTypeID: source.SourceTypeID,
			Uid:          source.Uid,
			Provider:     source.Provider.String(),
			Status:       source.Status,
		}
	}
	return &SourceListResponse{Data: list, Metadata: *meta}
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
