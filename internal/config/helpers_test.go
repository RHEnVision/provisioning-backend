package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvironmentPrefix2(t *testing.T) {
	require.Equal(t, "test-123-dev", EnvironmentPrefix("test", "123"))
}
