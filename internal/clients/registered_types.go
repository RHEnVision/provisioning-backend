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
		fmt.Printf("Total: %d\n", len(rit.types))
	}
}
