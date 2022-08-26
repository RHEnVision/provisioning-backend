package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCamelCaseReplacerEmpty(t *testing.T) {
	result := customReplacer{}.Replace("")
	assert.Equal(t, "", result)
}

func TestCamelCaseReplacerNothing(t *testing.T) {
	result := customReplacer{}.Replace("nothing")
	assert.Equal(t, "NOTHING", result)
}

func TestCamelCaseReplacerApp(t *testing.T) {
	result := customReplacer{}.Replace("APP.NAME")
	assert.Equal(t, "APP_NAME", result)
}

func TestCamelCaseReplacerFromMap(t *testing.T) {
	result := customReplacer{}.Replace("AWS.INSTANCEPREFIX")
	assert.Equal(t, "AWS_INSTANCE_PREFIX", result)
}
