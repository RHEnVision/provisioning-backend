package preload

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAzureETagValue(t *testing.T) {
	tag := AzureInstanceType.ETagValue()
	require.NotEmpty(t, tag.Name)
	require.NotEmpty(t, tag.Value)
}

func TestAzureFindInstanceType(t *testing.T) {
	it := AzureInstanceType.FindInstanceType("Standard_B2s")
	require.NotNil(t, it)
	require.Equal(t, "Standard_B2s", it.Name.String())
	require.True(t, it.Supported)
}

func TestAzureValidateRegion(t *testing.T) {
	require.True(t, AzureInstanceType.ValidateRegion("westeurope_1"))
	require.False(t, AzureInstanceType.ValidateRegion("centralprague_6"))
}
