package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateTagLength(t *testing.T) {
	tag := GenerateTag()
	assert.Len(t, tag, 20, "tag is not at length 20")
}
