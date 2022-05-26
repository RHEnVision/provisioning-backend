package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagFormat(t *testing.T) {

	pubkey := PubkeyResource{12, "randomtag", 25, 1, "aws-handle"}

	// assert equality
	assert.Equal(t, "pk-12-randomtag", pubkey.FormattedTag(), "the tag format is incorrect")
}
