package clients

import (
	"fmt"

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

func (rit *RegisteredInstanceTypes) Register(it InstanceType) {
	// Set supported flag (it should be 1536 but let's make it safe)
	if it.MemoryMiB >= 1500 {
		it.Supported = true
	} else {
		it.Supported = false
	}

	rit.types[it.Name] = &it
}

func (rit *RegisteredInstanceTypes) Get(name InstanceTypeName) *InstanceType {
	return rit.types[name]
}

func (rit *RegisteredInstanceTypes) Load(buffer []byte) error {
	err := yaml.Unmarshal(buffer, &rit.types)
	if err != nil {
		return fmt.Errorf("unable to unmarshal registered instance types: %w", err)
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
		fmt.Printf("Total: %d\n", len(rit.types))
	}
}
