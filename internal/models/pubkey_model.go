package models

// Pubkey represents SSH public key that can be deployed to clients.
type Pubkey struct {
	// Required auto-generated PK.
	ID int64 `db:"id" json:"id"`

	// Associated Account model. Required.
	AccountID int64 `db:"account_id" json:"-"`

	// User-facing name. Required.
	Name string `db:"name" json:"name"`

	// Public key body encoded in base64 (.pub format). Required.
	Body string `db:"body" json:"body"`
}
