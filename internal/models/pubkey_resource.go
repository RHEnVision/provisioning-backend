package models

import (
	"fmt"
)

// PubkeyResource represents uploaded SSH key-pair resource.
type PubkeyResource struct {
	// Required auto-generated PK.
	ID uint64 `db:"id" json:"id"`

	// Auto-generated random string, not DB-unique but pseudo-unique (128 bits),
	// not returned through API, read-only. Used to "tag" resources in clouds for
	// later use (e.g. search a resource by a tag). Use FormattedTag function to
	// format it with a proper prefix.
	Tag string `db:"tag" json:"-"`

	// Associated Account model. Required.
	PubkeyID uint64 `db:"pubkey_id" json:"pubkey_id"`

	// Provider constant (for example ProviderTypeAWS). Required.
	Provider ProviderType `db:"provider" json:"provider"`

	// Resource handle (id). Format is provider-dependant. Required.
	Handle string `db:"handle" json:"handle"`
}

// FormattedTag returns Tag concatenated in a safe way for clouds. That means
// it does not start with a number, only includes alpha-num characters and dash.
// Tag is
func (p *PubkeyResource) FormattedTag() string {
	return fmt.Sprintf("pk-%s", p.Tag)
}

// RandomizeTag sets Tag field via GenerateTag function if and only if the Tag is
// currently blank.
func (p *PubkeyResource) RandomizeTag() {
	if p.Tag == "" {
		p.Tag = GenerateTag()
	}
}
