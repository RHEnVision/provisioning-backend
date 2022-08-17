package models_test

import (
	"testing"

	_ "github.com/RHEnVision/provisioning-backend/internal/testing/initialization"

	"github.com/RHEnVision/provisioning-backend/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestTagFormat(t *testing.T) {
	pubkey := models.PubkeyResource{
		ID:       12,
		Tag:      "ctNLLGsipCJYjeoGXdWy17",
		PubkeyID: 25,
		SourceID: 1,
		Provider: 1,
		Handle:   "aws-handle",
	}

	assert.Equal(t, "pk-ctNLLGsipCJYjeoGXdWy17", pubkey.FormattedTag(), "the tag format is incorrect")
}

func TestRandomizeTagLength(t *testing.T) {
	pk := models.PubkeyResource{Tag: ""}
	pk.RandomizeTag()
	assert.Len(t, pk.Tag, 22, "tag is not at length 20")
}

func TestRandomizeTagNonEmpty(t *testing.T) {
	pk := models.PubkeyResource{Tag: "something"}
	pk.RandomizeTag()
	assert.Equal(t, pk.Tag, "something", "tag must be not overwritten")
}
