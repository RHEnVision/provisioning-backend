package types

import (
	"embed"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
)

//go:embed *.yaml availability/*.yaml
var fsTypes embed.FS

var registeredTypes clients.RegisteredInstanceTypes
var regionalAvailability clients.RegionalTypeAvailability

func init() {
	// load instance types
	buffer, err := fsTypes.ReadFile("types.yaml")
	if err != nil {
		panic(err)
	}
	err = registeredTypes.Load(buffer)
	if err != nil {
		panic(err)
	}

	// load supported types
	buffer, err = fsTypes.ReadFile("supported.yaml")
	if err != nil {
		panic(err)
	}
	err = registeredTypes.LoadSupported(buffer)
	if err != nil {
		panic(err)
	}

	// load availability information
	err = regionalAvailability.Load(fsTypes, "availability")
	if err != nil {
		panic(err)
	}
}

func PrintRegisteredTypes(typeName string) {
	registeredTypes.Print(typeName)
}

func PrintRegionalAvailability(region, zone string) {
	regionalAvailability.Print(region, zone)
}
