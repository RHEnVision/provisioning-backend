package payloads

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/go-chi/render"
)

// ReservationRequest is empty, account comes in HTTP header and
// provider type in HTTP URL. All other fields are auto-generated.
type ReservationRequest struct {
}

type ReservationResponse struct {
	*models.Reservation
}

type NoopReservationRequest ReservationRequest

type NoopReservationResponse struct {
	*models.NoopReservation
}

type AWSReservationRequest struct {
	Pubkey struct {
		ExistingID *int64  `json:"existing_id"`
		NewName    *string `json:"new_name"`
		NewBody    *string `json:"new_body"`
	}
}

type AWSReservationResponse struct {
	*models.AWSReservation
}

func (p *ReservationRequest) Bind(_ *http.Request) error {
	return nil
}

func (p *ReservationResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (p *AWSReservationRequest) Bind(_ *http.Request) error {
	return nil
}

func (p *AWSReservationResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (p *NoopReservationRequest) Bind(_ *http.Request) error {
	return nil
}

func (p *NoopReservationResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewReservationResponse(reservation *models.Reservation) render.Renderer {
	return &ReservationResponse{Reservation: reservation}
}

func NewAWSReservationResponse(reservation *models.AWSReservation) render.Renderer {
	return &AWSReservationResponse{AWSReservation: reservation}
}

func NewNoopReservationResponse(reservation *models.NoopReservation) render.Renderer {
	return &NoopReservationResponse{NoopReservation: reservation}
}

func NewReservationListResponse(accounts []*models.Reservation) []render.Renderer {
	list := make([]render.Renderer, 0, len(accounts))
	for _, a := range accounts {
		list = append(list, &ReservationResponse{Reservation: a})
	}
	return list
}
