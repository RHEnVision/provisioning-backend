package logging

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTruncateEmpty(t *testing.T) {
	result := truncateText("", 13)
	assert.Equal(t, "", result)
}

func TestTruncateNothing(t *testing.T) {
	result := truncateText("test", 13)
	assert.Equal(t, "test", result)
}

func TestTruncateExactly(t *testing.T) {
	result := truncateText("test", 4)
	assert.Equal(t, "test", result)
}

func TestTruncateOne(t *testing.T) {
	result := truncateText("test", 3)
	assert.Equal(t, "tes...\"", result)
}
