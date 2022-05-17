package models

// Pubkey represents a SSH public key that can be deployed to clouds.
type Pubkey struct {
	// Required auto-generated PK.
	ID uint64 `db:"id" json:"id"`

	// Associated Account model. Required.
	AccountID uint64 `db:"account_id" json:"account_id"`

	// User-facing name. Required.
	Name string `db:"name" json:"name"`

	// Public key body encoded in base64 (.pub format). Required.
	Body string `db:"body" json:"body"`
}
