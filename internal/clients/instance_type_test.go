package clients

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstanceType_SetMemoryFromGiB(t *testing.T) {
	it := InstanceType{}
	it.SetMemoryFromGiB(64)
	assert.Equal(t, int64(65_536), it.MemoryMiB)
}

func TestInstanceType_SetMemoryFromKiB(t *testing.T) {
	it := InstanceType{}
	it.SetMemoryFromKiB(67_108_864)
	assert.Equal(t, int64(65_536), it.MemoryMiB)
}

func TestInstanceType_SetMemoryFromBytes(t *testing.T) {
	it := InstanceType{}
	it.SetMemoryFromBytes(67_108_864)
	assert.Equal(t, int64(64), it.MemoryMiB)
}

func TestInstanceType_SetEphemeralStorageFromMB(t *testing.T) {
	it := InstanceType{}
	it.SetEphemeralStorageFromMB(320_000)
	assert.Equal(t, int64(320), it.EphemeralStorageGB)
}
