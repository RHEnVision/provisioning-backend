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

	// Pubkey ID.
	PubkeyID sql.NullInt64 `db:"pubkey_id" json:"pubkey_id"`

	// Time when reservation was made.
	CreatedAt time.Time `db:"created_at" json:"created_at"`

	// Textual status of the reservation or error when there was a failure
	Status string `db:"status" json:"status"`

	// Time when reservation was finished or nil when it's still processing.
	FinishedAt sql.NullTime `db:"finished_at" json:"finished_at"`

	// Flag indicating success, error or unknown state (NULL). See Status for the actual error.
	Success sql.NullBool `db:"success" json:"success"`
}
