package clients

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type InstanceTypeInfo struct {
	RegisteredTypes      RegisteredInstanceTypes
	RegionalAvailability RegionalTypeAvailability
}

func (iii *InstanceTypeInfo) InstanceTypesForZone(region, zone string, supported *bool) ([]*InstanceType, error) {
	result := make([]*InstanceType, 0, 64)
	names, err := iii.RegionalAvailability.NamesForZone(region, zone)
	if err != nil {
		return nil, err
	}
	for _, name := range names {
		rt := iii.RegisteredTypes.Get(name)
		if supported != nil && *supported != rt.Supported {
			continue
		}
		result = append(result, rt)
	}
	return result, nil
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
