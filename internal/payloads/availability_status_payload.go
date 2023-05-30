package payloads

import (
	"net/http"
)

type AvailabilityStatusRequest struct {
	SourceID string `json:"source_id" yaml:"source_id"`
}

func (p *AvailabilityStatusRequest) Bind(_ *http.Request) error {
	return nil
}
