package clients

import (
	"fmt"
	"reflect"

	"gopkg.in/yaml.v3"
)

// RegisteredInstanceTypes holds all details about instance types.
type RegisteredInstanceTypes struct {
	types map[InstanceTypeName]*InstanceType `yaml:"registered_types"`
}

func NewRegisteredInstanceTypes() *RegisteredInstanceTypes {
	return &RegisteredInstanceTypes{
		types: make(map[InstanceTypeName]*InstanceType),
	}
}

// Register puts instance type into the list and sets the supported flag. Currently, only
// instances with more than 1.5 GB (not GiB) are considered as supported.
//
// The function prints a warning to standard input if a type was already registered but has
// a different fields. Some hyperscalers (e.g. Azure) can have different attributes for
// the same types in different zones (e.g. ephemeral storage size). Unless there is a bigger
// difference, this isn't a problem. This helps to track these during generation.
func (rit *RegisteredInstanceTypes) Register(it InstanceType) {
	// Set supported flag (it should be 1536 but let's make it safe)
	if it.MemoryMiB >= 1500 {
		it.Supported = true
	} else {
		it.Supported = false
	}

	// Do not allow registering different types under same name.
	if existing, ok := rit.types[it.Name]; ok {
		if !reflect.DeepEqual(*existing, it) {
			fmt.Printf("WARNING: registering %s instance type that has different attributes:\n existing: %+v\n new: %+v\n", it.Name, *existing, it)
		}
	}

	rit.types[it.Name] = &it
}

// Get returns instance type by name or nil when such type does not exist.
func (rit *RegisteredInstanceTypes) Get(name InstanceTypeName) *InstanceType {
	return rit.types[name]
}

// Load existing instances from YAML buffer
func (rit *RegisteredInstanceTypes) Load(buffer []byte) error {
	err := yaml.Unmarshal(buffer, &rit.types)
	if err != nil {
		return fmt.Errorf("unable to unmarshal registered instance types: %w", err)
	}

	return nil
}

// Save instance list to YAML
func (rit *RegisteredInstanceTypes) Save(filename string) error {
	return compareAndMarshal(filename, rit.types)
}

// Print is useful for debugging
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
		fmt.Printf("Total: %d\n", len(rit.types))
	}
}
