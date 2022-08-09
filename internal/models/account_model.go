package models

import "database/sql"

// Account represents a Red Hat Console account
type Account struct {
	// Required auto-generated PK.
	ID int64 `db:"id"`

	// Organization ID. Required.
	OrgID string `db:"org_id"`

	// EBS account number. Can be NULL but not blank.
	AccountNumber sql.NullString `db:"account_number"`
}
