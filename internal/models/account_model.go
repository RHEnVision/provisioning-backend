package models

// Account represents a Red Hat Console account
type Account struct {
	// Required auto-generated PK.
	ID uint64 `db:"id"`

	// Organization ID. Required.
	OrgID string `db:"org_id"`

	// EBS account number. Can be NULL but not blank.
	AccountNumber string `db:"account_number"`
}
