package models_test

import (
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGenerateTagLength(t *testing.T) {
	tag := models.GenerateTag()
	assert.Len(t, tag, 22, "tag is not at length 22")
}
