package models

import (
	"database/sql"
	"time"
)

// Reservation represents an instance launch reservation. They are associated with a background
// job system with a particular job with its own ID. The function handlers update the reservation
// Status, Success and FinishedAt attributes until the job is considered finished.
type Reservation struct {
	// Required auto-generated PK.
	ID int64 `db:"id" json:"id"`

	// Provider type. Required.
	Provider ProviderType `db:"provider" json:"provider"`

	// Account ID. Required.
	AccountID int64 `db:"account_id" json:"-"`

	// Time when reservation was made.
	CreatedAt time.Time `db:"created_at" json:"created_at"`

	// Textual status of the reservation or error when there was a failure
	Status string `db:"status" json:"status"`

	// Time when reservation was finished or nil when it's still processing.
	FinishedAt sql.NullTime `db:"finished_at" json:"finished_at"`

	// Flag indicating success, error or unknown state (NULL). See Status for the actual error.
	Success sql.NullBool `db:"success" json:"success"`
}

type NoopReservation struct {
	Reservation
}

type AWSReservation struct {
	Reservation

	// Pubkey ID.
	PubkeyID sql.NullInt64 `db:"pubkey_id" json:"pubkey_id"`

	// Source ID.
	SourceID int64 `db:"source_id" json:"source_id"`

	//AWS Instance type.
	InstanceType string `db:"instance_type" json:"instance_type"`

	// Amount of instances to provision of type: Instance type.
	Amount int64 `db:"amount" json:"amount"`

	// The ID of the image from which the instance is created.
	ImageID int64 `db:"image_id" json:"image_id"`
}
