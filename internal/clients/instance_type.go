package clients

import (
	"strconv"
	"strings"
)

type InstanceTypeName string

const (
	_ = 1 << (10 * iota)
	kiB
	miB
)

// InstanceType defines a model for an instance type that corresponds to one in a cloud provider.
type InstanceType struct {
	// The name of the instance type
	Name InstanceTypeName `json:"name,omitempty" yaml:"name"`

	// Virtual CPU (maps to hypervisor hyper-thread)
	VCPUs int32 `json:"vcpus,omitempty" yaml:"vcpus"`

	// Core (physical or virtual core)
	Cores int32 `json:"cores,omitempty" yaml:"cores"`

	// The size of the memory, in MiB.
	MemoryMiB int64 `json:"memory_mib,omitempty" yaml:"memory_mib"`

	// The total size of ephemeral disks, in GB. Is set to 0 if local disk(s) are not available.
	EphemeralStorageGB int64 `json:"storage_gb" yaml:"storage_gb"`

	// Does the instance type supports RHEL
	Supported bool `json:"supported" yaml:"supported"`

	// Instance type's Architecture: i386, arm64, x86_64
	Architecture ArchitectureType `json:"architecture,omitempty" yaml:"arch"`

	// Extra information for Azure, nil for other types
	AzureDetail *InstanceTypeDetailAzure `json:"azure,omitempty" yaml:"azure,omitempty"`
}

// InstanceTypeDetailAzure contains specific details for Azure.
type InstanceTypeDetailAzure struct {
	GenV1 bool `json:"gen_v1" yaml:"gen_v1"`
	GenV2 bool `json:"gen_v2" yaml:"gen_v2"`
}

func (it *InstanceTypeName) String() string {
	return string(*it)
}

func (it *InstanceType) String() string {
	sb := strings.Builder{}
	sb.WriteString(it.Name.String())
	sb.WriteString(" | Arch: ")
	sb.WriteString(it.Architecture.String())
	sb.WriteString(" | vCPUs: ")
	sb.WriteString(strconv.Itoa(int(it.VCPUs)))
	sb.WriteString(" | Cores: ")
	sb.WriteString(strconv.Itoa(int(it.Cores)))
	sb.WriteString(" | Memory: ")
	sb.WriteString(strconv.Itoa(int(it.MemoryMiB)))
	sb.WriteString(" MiB | Disk: ")
	sb.WriteString(strconv.Itoa(int(it.EphemeralStorageGB)))
	sb.WriteString(" GB")
	if it.AzureDetail != nil {
		sb.WriteString(" | Azure Gen:")
		if it.AzureDetail.GenV1 {
			sb.WriteString(" V1")
		}
		if it.AzureDetail.GenV2 {
			sb.WriteString(" V2")
		}
	}
	sb.WriteString(" | Supported:")
	if it.Supported {
		sb.WriteString(" Yes")
	} else {
		sb.WriteString(" No")
	}
	return sb.String()
}

func (it *InstanceType) SetMemoryFromGiB(memGib int64) {
	it.MemoryMiB = memGib * kiB
}

func (it *InstanceType) SetMemoryFromKiB(memKib int64) {
	it.MemoryMiB = memKib / kiB
}

func (it *InstanceType) SetMemoryFromBytes(memKib int64) {
	it.MemoryMiB = memKib / miB
}

func (it *InstanceType) SetEphemeralStorageFromMB(storageMb int64) {
	it.EphemeralStorageGB = storageMb / 1000
}
