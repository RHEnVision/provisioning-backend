package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagFormat(t *testing.T) {
	pubkey := PubkeyResource{12, "randomtag", 25, 1, "aws-handle"}

	assert.Equal(t, "pk-12-randomtag", pubkey.FormattedTag(), "the tag format is incorrect")
}

func TestRandomizeTagLength(t *testing.T) {
	pk := PubkeyResource{Tag: ""}
	pk.RandomizeTag()
	assert.Len(t, pk.Tag, 20, "tag is not at length 20")
}

func TestRandomizeTagNonEmpty(t *testing.T) {
	pk := PubkeyResource{Tag: "something"}
	pk.RandomizeTag()
	assert.Equal(t, pk.Tag, "something", "tag must be not overwritten")
}
