package types

import (
	"bytes"
	"embed"
	"fmt"
	"path/filepath"
	"time"

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

	// load supported types
	supportedBuf, err := fsTypes.ReadFile("supported.yaml")
	if err != nil {
		panic(err)
	}
	err = typeInfo.RegisteredTypes.LoadSupported(supportedBuf)
	if err != nil {
		panic(err)
	}

	// load availability information
	err = typeInfo.RegionalAvailability.Load(fsTypes, availabilityPath)
	if err != nil {
		panic(err)
	}

	// calculate etag
	// concatenate all yamls into one 1 MiB buffer
	dirEntries, err := fsTypes.ReadDir(availabilityPath)
	if err != nil {
		panic(err)
	}
	availBuf := bytes.NewBuffer(make([]byte, 0, 2^20))
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}
		file := filepath.Join(availabilityPath, dirEntry.Name())
		buffer, errBuf := fsTypes.ReadFile(file)
		if errBuf != nil {
			panic(errBuf)
		}
		availBuf.Write(buffer)
	}
	etag, err = middleware.GenerateETagFromBuffer("azure-types", 30*time.Minute, typesBuf, supportedBuf, availBuf.Bytes())
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
	typeInfo.RegionalAvailability.Print(region, zone)
}

func InstanceTypesForZone(region, zone string, supported *bool) ([]*clients.InstanceType, error) {
	result, err := typeInfo.InstanceTypesForZone(region, zone, supported)
	if err != nil {
		return nil, fmt.Errorf("unable to list instance types for region and zone: %w", err)
	}
	return result, nil
}
