package kafka

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateHeadersEmpty(t *testing.T) {
	h := GenericHeaders()
	require.Empty(t, h)
}

func TestGenerateHeadersOne(t *testing.T) {
	h := GenericHeaders("a", "b")
	require.Len(t, h, 1)
	require.Equal(t, GenericHeader{
		Key:   "a",
		Value: "b",
	}, h[0])
}

func TestGenerateHeadersTwo(t *testing.T) {
	h := GenericHeaders("a", "b", "c", "d")
	require.Len(t, h, 2)
	require.Equal(t, GenericHeader{
		Key:   "a",
		Value: "b",
	}, h[0])
	require.Equal(t, GenericHeader{
		Key:   "c",
		Value: "d",
	}, h[1])
}

// nolint:staticcheck
func TestGenerateHeadersOdd(t *testing.T) {
	require.Panicsf(t, func() {
		_ = GenericHeaders("")
	}, "generic headers: odd amount of arguments")
}
