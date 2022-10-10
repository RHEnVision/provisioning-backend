package clients

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var smallType = InstanceType{
	Name:         "small",
	VCPUs:        1,
	Cores:        1,
	MemoryMiB:    1500,
	Architecture: ArchitectureTypeX86_64,
}

func TestAddSmall(t *testing.T) {
	rit := NewRegionalInstanceTypes()
	rit.Add("region", "zone", smallType)
	require.Equal(t, "\nRegion 'region' availability zone 'zone': small\n", rit.Sprint("region", "zone"))
}

func TestAddSmallTwice(t *testing.T) {
	rit := NewRegionalInstanceTypes()
	rit.Add("region", "zone", smallType)
	rit.Add("region", "zone", smallType)
	require.Equal(t, "\nRegion 'region' availability zone 'zone': small\n", rit.Sprint("region", "zone"))
}

func TestAddSmallNoZone(t *testing.T) {
	rit := NewRegionalInstanceTypes()
	rit.Add("region", "", smallType)
	rit.Add("region", "", smallType)
	require.Equal(t, "\nRegion 'region' availability zone '': small\n", rit.Sprint("region", ""))
}

func TestAddSmallDifferentZones(t *testing.T) {
	rit := NewRegionalInstanceTypes()
	rit.Add("region", "zone1", smallType)
	rit.Add("region", "zone2", smallType)
	require.Equal(t, "\nRegion 'region' availability zone 'zone1': small\n", rit.Sprint("region", "zone1"))
	require.Equal(t, "\nRegion 'region' availability zone 'zone2': small\n", rit.Sprint("region", "zone2"))
}
