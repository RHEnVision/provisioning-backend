package types

import (
	"embed"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/middleware"
)

//go:embed *.yaml availability/*.yaml
var fsTypes embed.FS

const availabilityPath = "availability"

var etag *middleware.ETag

var typeInfo clients.InstanceTypeInfo

func init() {
	// load instance types
	typesBuf, err := fsTypes.ReadFile("types.yaml")
	if err != nil {
		panic(err)
	}
	err = typeInfo.RegisteredTypes.Load(typesBuf)
	if err != nil {
		panic(err)
	}

	// load availability information
	err = typeInfo.RegionalAvailability.Load(fsTypes, availabilityPath)
	if err != nil {
		panic(err)
	}

	availBuf := clients.ConcatBuffers(fsTypes, availabilityPath)
	etag, err = middleware.GenerateETagFromBuffer("gcp-types", middleware.InstanceTypeExpiration, typesBuf, availBuf)
	if err != nil {
		panic(err)
	}
}

func ETagValue() *middleware.ETag {
	return etag
}

func PrintRegisteredTypes(typeName string) {
	typeInfo.RegisteredTypes.Print(typeName)
}

func PrintRegionalAvailability(region, zone string) {
	str := typeInfo.RegionalAvailability.Sprint(region, zone)
	fmt.Println(str)
}

func InstanceTypesForZone(region, zone string, supported *bool) ([]*clients.InstanceType, error) {
	result, err := typeInfo.InstanceTypesForZone(region, zone, supported)
	if err != nil {
		return nil, fmt.Errorf("unable to list instance types for region and zone: %w", err)
	}
	return result, nil
}
