package preload

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEC2ETagValue(t *testing.T) {
	tag := EC2InstanceType.ETagValue()
	require.NotEmpty(t, tag.Name)
	require.NotEmpty(t, tag.Value)
}

func TestEC2FindInstanceType(t *testing.T) {
	it := EC2InstanceType.FindInstanceType("m1.small")
	require.NotNil(t, it)
	require.Equal(t, "m1.small", it.Name.String())
	require.True(t, it.Supported)
}

func TestEC2ValidateRegion(t *testing.T) {
	require.True(t, EC2InstanceType.ValidateRegion("us-east-1"))
	require.False(t, EC2InstanceType.ValidateRegion("cz-olomouc-2"))
}
