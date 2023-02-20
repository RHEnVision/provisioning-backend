package preload

import (
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/middleware"
)

type instanceType struct {
	filename string
	path     string
	etagName string
	tag      *middleware.ETag
	typeInfo clients.InstanceTypeInfo
}

func (p *instanceType) Load() error {
	// load instance types
	typesBuf, err := fsTypes.ReadFile(p.filename)
	if err != nil {
		return fmt.Errorf("unable to read instance types %s: %w", p.filename, err)
	}

	err = p.typeInfo.RegisteredTypes.Load(typesBuf)
	if err != nil {
		return fmt.Errorf("unable to load instance types %s: %w", p.filename, err)
	}

	// load availability information
	err = p.typeInfo.RegionalAvailability.Load(fsTypes, p.path)
	if err != nil {
		return fmt.Errorf("unable to load regional info %s: %w", p.path, err)
	}

	availBuf := clients.ConcatBuffers(fsTypes, p.path)
	p.tag, err = middleware.GenerateETagFromBuffer(p.etagName, middleware.InstanceTypeExpiration, typesBuf, availBuf)
	if err != nil {
		return fmt.Errorf("unable to generate etag %s: %w", p.etagName, err)
	}
	return nil
}

// ETagValue returns HTTP ETag information. It is calculated as a hash from source YAML files.
func (p *instanceType) ETagValue() *middleware.ETag {
	return p.tag
}

// PrintRegisteredTypes prints relevant data to standard output.
func (p *instanceType) PrintRegisteredTypes(typeName string) {
	p.typeInfo.RegisteredTypes.Print(typeName)
}

// PrintRegionalAvailability prints relevant data to standard output.
func (p *instanceType) PrintRegionalAvailability(region, zone string) {
	str := p.typeInfo.RegionalAvailability.Sprint(region, zone)
	fmt.Println(str)
}

// InstanceTypesForZone returns instance type info for particular zone. Can list supported, unsupported
// or all types when nil is passed.
func (p *instanceType) InstanceTypesForZone(region, zone string, supported *bool) ([]*clients.InstanceType, error) {
	result, err := p.typeInfo.InstanceTypesForZone(region, zone, supported)
	if err != nil {
		return nil, fmt.Errorf("unable to list instance types for region and zone: %w", err)
	}
	return result, nil
}

// FindInstanceType looks up instance type by name.
func (p *instanceType) FindInstanceType(name clients.InstanceTypeName) *clients.InstanceType {
	return p.typeInfo.RegisteredTypes.Get(name)
}
