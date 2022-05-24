package models

import "fmt"

// PubkeyResource represents uploaded SSH key-pair resource.
type PubkeyResource struct {
	// Required auto-generated PK.
	ID uint64 `db:"id" json:"id"`

	// Auto-generated random string, not unique, not returned through API, read-only.
	Tag string `db:"tag" json:"-"`

	// Associated Account model. Required.
	PubkeyID uint64 `db:"pubkey_id" json:"pubkey_id"`

	// Provider constant (for example ProviderTypeAWS). Required.
	Provider int `db:"provider" json:"provider"`

	// Resource handle (id). Format is provider-dependant. Required.
	Handle string `db:"handle" json:"handle"`
}

// FormattedTag returns ID and Tag concatenated in a safe way for clouds. That means
// it does not start with a number, only includes alpha-num characters and dash.
func (p *PubkeyResource) FormattedTag() string {
	return fmt.Sprintf("pk-%d-%s", p.ID, p.Tag)
}
