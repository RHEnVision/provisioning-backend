package payloads

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/go-chi/render"
)

type ReservationPayload struct {
	*models.Reservation
}

type ReservationRequest struct {
	Pubkey struct {
		ExistingID *int64  `json:"existing_id"`
		NewName    *string `json:"new_name"`
		NewBody    *string `json:"new_body"`
	}
}

type ReservationResponse ReservationPayload

func (p *ReservationRequest) Bind(_ *http.Request) error {
	return nil
}

func (p *ReservationResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewReservationResponse(account *models.Reservation) render.Renderer {
	return &ReservationResponse{Reservation: account}
}

func NewReservationListResponse(accounts []*models.Reservation) []render.Renderer {
	list := make([]render.Renderer, 0, len(accounts))
	for _, a := range accounts {
		list = append(list, &ReservationResponse{Reservation: a})
	}
	return list
}
