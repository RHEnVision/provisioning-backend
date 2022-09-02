package clients

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type InstanceTypeName string
type ArchitectureType string

const (
	ArchitectureTypeI386       ArchitectureType = "i386"
	ArchitectureTypeX8664      ArchitectureType = "x86_64"
	ArchitectureTypeArm64      ArchitectureType = "arm64"
	ArchitectureTypeAppleX8664 ArchitectureType = "apple-x86_64"
	ArchitectureTypeAppleArm64 ArchitectureType = "apple-arm64"
)

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
	Supported bool `json:"supported,omitempty" yaml:"supported"`

	// Instance type's Architecture: i386, arm64, x86_64
	Architecture ArchitectureType `json:"architecture,omitempty" yaml:"arch"`

	// Extra information for Azure, nil for other types
	AzureDetail *InstanceTypeDetailAzure `json:"azure,omitempty" yaml:"azure"`
}

// InstanceTypeDetailAzure contains specific details for Azure.
type InstanceTypeDetailAzure struct {
	GenV1 bool `json:"gen_v1" yaml:"gen_v1"`
	GenV2 bool `json:"gen_v2" yaml:"gen_v2"`
}

func (it *InstanceTypeName) String() string {
	return string(*it)
}

func (at *ArchitectureType) String() string {
	return string(*at)
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

// RegisteredInstanceTypes holds all details about instance types.
type RegisteredInstanceTypes struct {
	types map[InstanceTypeName]*InstanceType `yaml:"registered_types"`
}

func NewRegisteredInstanceTypes() *RegisteredInstanceTypes {
	return &RegisteredInstanceTypes{
		types: make(map[InstanceTypeName]*InstanceType),
	}
}

func (rit *RegisteredInstanceTypes) Register(it InstanceType) {
	rit.types[it.Name] = &it
}

func (rit *RegisteredInstanceTypes) Load(buffer []byte) error {
	err := yaml.Unmarshal(buffer, &rit.types)
	if err != nil {
		return fmt.Errorf("unable to unmarshal registered instance types: %w", err)
	}

	return nil
}

func (rit *RegisteredInstanceTypes) LoadSupported(buffer []byte) error {
	supportedList := make([]InstanceTypeName, 0)
	err := yaml.Unmarshal(buffer, &supportedList)
	if err != nil {
		return fmt.Errorf("unable to unmarshal supported instance types: %w", err)
	}
	supported := make(map[InstanceTypeName]bool, len(supportedList))
	for _, name := range supportedList {
		supported[name] = true
	}

	for itn, it := range rit.types {
		if _, ok := supported[itn]; ok {
			it.Supported = true
		} else {
			it.Supported = false
		}
	}

	return nil
}

func (rit *RegisteredInstanceTypes) Save(filename string) error {
	return compareAndMarshal(filename, rit.types)
}

func (rit *RegisteredInstanceTypes) Print(typeName string) {
	if typeName != "" {
		it, ok := rit.types[InstanceTypeName(typeName)]
		if ok {
			fmt.Println(it.String())
		} else {
			fmt.Println("Not found")
		}
	} else {
		for _, v := range rit.types {
			fmt.Println(v.String())
		}
		fmt.Printf("Total: %d", len(rit.types))
	}
}

// RegionalTypeAvailability type is used to capture available instance types per
// region and zone.
type RegionalTypeAvailability struct {
	types map[string][]InstanceTypeName
}

const regionSeparator = "_"

func NewRegionalInstanceTypes() *RegionalTypeAvailability {
	return &RegionalTypeAvailability{
		types: make(map[string][]InstanceTypeName),
	}
}

func (rit *RegionalTypeAvailability) Add(region, zone string, it InstanceType) {
	raz := region + regionSeparator + zone
	if _, ok := rit.types[raz]; !ok {
		rit.types[raz] = make([]InstanceTypeName, 0)
	}
	rit.types[raz] = append(rit.types[raz], it.Name)
}

func (rit *RegionalTypeAvailability) Save(directory string) error {
	for key, value := range rit.types {
		filename := filepath.Join(directory, key+".yaml")
		err := compareAndMarshal(filename, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (rit *RegionalTypeAvailability) Load(fsTypes embed.FS, path string) error {
	rit.types = make(map[string][]InstanceTypeName)

	dirEntries, err := fsTypes.ReadDir(path)
	if err != nil {
		return fmt.Errorf("unable to read availability dir: %w", err)
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}
		file := filepath.Join(path, dirEntry.Name())
		buffer, err := fsTypes.ReadFile(file)
		if err != nil {
			return fmt.Errorf("unable to read availability file %s: %w", file, err)
		}
		key := strings.TrimSuffix(dirEntry.Name(), ".yaml")
		var value []InstanceTypeName
		err = yaml.Unmarshal(buffer, &value)
		if err != nil {
			return fmt.Errorf("unable to unmarshal availability file %s: %w", file, err)
		}
		rit.types[key] = value
	}

	return nil
}

var RegionAndZoneSplitErr = errors.New("unable to split region and zone for")

func splitRegionZone(str string) (string, string, error) {
	result := strings.Split(str, regionSeparator)
	if len(result) != 2 {
		return "", "", fmt.Errorf("%w: %s", RegionAndZoneSplitErr, str)
	}
	return result[0], result[1], nil
}

func (rit *RegionalTypeAvailability) Print(fRegion, fZone string) {
	for raz, names := range rit.types {
		region, zone, err := splitRegionZone(raz)
		if err != nil {
			panic(err)
		}
		if (fRegion == "" && fZone == "") ||
			(fRegion == region && fZone == "") ||
			(fRegion == region && fZone == zone) ||
			(fRegion == "all" && fZone == "") {
			fmt.Printf("Region '%s' availability zone '%s':\n", region, zone)
			sb := strings.Builder{}
			for _, name := range names {
				sb.WriteString(name.String())
				sb.WriteString(", ")
			}
			fmt.Println(sb.String())
			fmt.Println("")
		}
	}
}

func compareAndMarshal(filename string, obj any) error {
	newBuffer, err := yaml.Marshal(obj)
	if err != nil {
		return fmt.Errorf("unable to marshal instance types: %w", err)
	}

	oldBuffer, err := os.ReadFile(filename)
	if os.IsNotExist(err) {
		oldBuffer = make([]byte, 0, 0)
	} else if err != nil {
		return fmt.Errorf("unable to read instance types: %w", err)
	}

	if !bytes.Equal(newBuffer, oldBuffer) {
		/* #nosec */
		err = os.WriteFile(filename, newBuffer, 0644)
		if err != nil {
			return fmt.Errorf("unable to save instance types: %w", err)
		}
	}

	return nil
}
