package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlank1(t *testing.T) {
	require.False(t, present(""))
}

func TestBlank2(t *testing.T) {
	require.False(t, present("", ""))
}

func TestBlankNon1(t *testing.T) {
	require.True(t, present("x"))
}

func TestBlankNon2(t *testing.T) {
	require.True(t, present("x", "x"))
}
