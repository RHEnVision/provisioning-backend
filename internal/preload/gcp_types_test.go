package preload

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGCPETagValue(t *testing.T) {
	tag := GCPInstanceType.ETagValue()
	require.NotEmpty(t, tag.Name)
	require.NotEmpty(t, tag.Value)
}

func TestGCPFindInstanceType(t *testing.T) {
	it := GCPInstanceType.FindInstanceType("e2-standard-2")
	require.NotNil(t, it)
	require.Equal(t, "e2-standard-2", it.Name.String())
	require.True(t, it.Supported)
}

func TestGCPValidateRegion(t *testing.T) {
	require.True(t, GCPInstanceType.ValidateRegion("europe-west1-b"))
	require.False(t, GCPInstanceType.ValidateRegion("velky-tynec7-b"))
}
