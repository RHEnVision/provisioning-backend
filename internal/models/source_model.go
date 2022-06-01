package models

type Source struct {
	// Required auto-generated PK.
	ID uint64 `db:"id" json:"id"`

	// Source ID in sources. Unique.
	SourceID uint64 `db:"source_id" json:"source_id"`

	// User-facing name. Required.
	Name string `db:"name" json:"name"`

	// Associated Account model. Required.
	AccountID uint64 `db:"account_id" json:"account_id"`

	// Authentication ID in sources lined to a spesific source in order to fetch the ARN
	AuthID string `db:"auth_id" json:"auth_id"`

	CreatedAt string `db:"created_at" json:"created_at"`
}
