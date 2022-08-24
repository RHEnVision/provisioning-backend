package payloads

import (
	"net/http"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/go-chi/render"
)

// ReservationRequest is empty, account comes in HTTP header and
// provider type in HTTP URL. All other fields are auto-generated.

type GenericReservationResponsePayload struct {
	ID int64 `json:"id"`

	// Provider type. Required.
	Provider int `json:"provider"`

	// Time when reservation was made.
	CreatedAt time.Time `json:"created_at"`

	// Textual status of the reservation or error when there was a failure
	Status string `json:"status"`

	// Time when reservation was finished or nil when it's still processing.
	FinishedAt *time.Time `json:"finished_at"`

	// Flag indicating success, error or unknown state (NULL). See Status for the actual error.
	Success *bool `json:"success"`
}

type AWSReservationResponsePayload struct {
	ID int64 `json:"reservation_id"`

	// Pubkey ID.
	PubkeyID int64 `json:"pubkey_id"`

	// Source ID.
	SourceID string `json:"source_id"`

	//AWS Instance type.
	InstanceType string `json:"instance_type"`

	// Amount of instances to provision of type: Instance type.
	Amount int32 `json:"amount"`

	// The ID of the image from which the instance is created.
	ImageID string `json:"image_id"`

	// The ID of the aws reservation which was created.
	AWSReservationID string `json:"aws_reservation_id"`
}

type NoopReservationResponsePayload struct {
	ID int64 `json:"reservation_id"`
}

type AWSReservationRequestPayload struct {
	// Pubkey ID.
	PubkeyID int64 `json:"pubkey_id"`

	// Source ID.
	SourceID string `json:"source_id"`

	// Optional name of the instance(s).
	Name string `json:"name"`

	// AWS Instance type.
	InstanceType string `json:"instance_type"`

	// Amount of instances to provision of type: Instance type.
	Amount int32 ` json:"amount"`

	// Image Builder UUID of the image that should be launched. AMI is also supported.
	ImageID string `json:"image_id"`

	// Immediately power off the system after initialization
	PowerOff bool `json:"poweroff"`
}

func (p *GenericReservationResponsePayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (p *AWSReservationRequestPayload) Bind(_ *http.Request) error {
	return nil
}

func (p *AWSReservationResponsePayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (p *NoopReservationResponsePayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewAWSReservationResponse(reservation *models.AWSReservation) render.Renderer {
	response := AWSReservationResponsePayload{
		ImageID:          reservation.ImageID,
		SourceID:         reservation.SourceID,
		Amount:           reservation.Amount,
		InstanceType:     reservation.InstanceType,
		AWSReservationID: reservation.AWSReservationID,
		ID:               reservation.ID,
	}
	return &response
}

func NewNoopReservationResponse(reservation *models.NoopReservation) render.Renderer {
	return &NoopReservationResponsePayload{
		ID: reservation.ID,
	}
}

func NewReservationListResponse(reservations []*models.Reservation) []render.Renderer {
	list := make([]render.Renderer, 0, len(reservations))
	for _, reservation := range reservations {
		list = append(list, &GenericReservationResponsePayload{
			ID:         reservation.ID,
			Provider:   int(reservation.Provider),
			CreatedAt:  reservation.CreatedAt,
			FinishedAt: &reservation.FinishedAt.Time,
			Status:     reservation.Status,
			Success:    &reservation.Success.Bool,
		})
	}
	return list
}
