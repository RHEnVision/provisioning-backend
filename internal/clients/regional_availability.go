package clients

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type sortableInstanceTypeName []InstanceTypeName

func (a sortableInstanceTypeName) Len() int {
	return len(a)
}

func (a sortableInstanceTypeName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a sortableInstanceTypeName) Less(i, j int) bool {
	return a[i] < a[j]
}

// RegionalTypeAvailability type is used to capture available instance types per
// region and zone.
type RegionalTypeAvailability struct {
	types map[string]sortableInstanceTypeName
}

const regionSeparator = "_"

func NewRegionalInstanceTypes() *RegionalTypeAvailability {
	return &RegionalTypeAvailability{
		types: make(map[string]sortableInstanceTypeName),
	}
}

var UnknownRegionZoneCombinationErr error = errors.New("unknown region and zone combination")

func key(region, zone string) string {
	if zone == "" {
		return region
	}
	return region + regionSeparator + zone
}

func (rit *RegionalTypeAvailability) NamesForZone(region, zone string) ([]InstanceTypeName, error) {
	result, ok := rit.types[key(region, zone)]
	if !ok {
		return nil, UnknownRegionZoneCombinationErr
	}
	return result, nil
}

func (rit *RegionalTypeAvailability) Add(region, zone string, it InstanceType) {
	raz := key(region, zone)
	if _, ok := rit.types[raz]; !ok {
		rit.types[raz] = make([]InstanceTypeName, 0)
	}
	rit.types[raz] = append(rit.types[raz], it.Name)
}

func (rit *RegionalTypeAvailability) Save(directory string) error {
	for key, value := range rit.types {
		sort.Sort(value)
		filename := filepath.Join(directory, key+".yaml")
		err := compareAndMarshal(filename, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (rit *RegionalTypeAvailability) Load(fsTypes embed.FS, path string) error {
	rit.types = make(map[string]sortableInstanceTypeName)

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

	if len(result) == 2 {
		return result[0], result[1], nil
	} else if len(result) == 1 {
		return result[0], "", nil
	} else {
		return "", "", fmt.Errorf("%w: %s", RegionAndZoneSplitErr, str)
	}
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

func ConcatBuffers(fsTypes embed.FS, path string) []byte {
	result := bytes.NewBuffer(make([]byte, 0))
	dirEntries, err := fsTypes.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}
		file := filepath.Join(path, dirEntry.Name())
		buffer, errBuf := fsTypes.ReadFile(file)
		if errBuf != nil {
			panic(errBuf)
		}
		result.Write(buffer)
	}
	return result.Bytes()
}
