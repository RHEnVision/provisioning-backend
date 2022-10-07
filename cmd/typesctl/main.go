package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/RHEnVision/provisioning-backend/cmd/typesctl/providers"
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/azure"
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/ec2"
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/gcp"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
)

func main() {
	config.Initialize("configs/api.env", "configs/typesctl.env")
	logging.InitializeStdout()

	validProviders := make([]string, 0)
	for p := range providers.TypeProviders {
		validProviders = append(validProviders, p)
	}
	helpProviders := strings.Join(validProviders, ",")

	providerFlag := flag.String("provider", "", fmt.Sprintf("provider: [%s] (required)", helpProviders))
	printAllFlag := flag.Bool("all", false, "print everything (long output)")
	printTypeFlag := flag.String("type", "", "print specific instance type detail (or 'all')")
	printRegionFlag := flag.String("region", "", "print instance type names for a region (or 'all')")
	printZoneFlag := flag.String("zone", "", "print instance type names for a zone (region is needed too)")
	generateFlag := flag.Bool("generate", false, "generate new type information")
	flag.Parse()

	provider, ok := providers.TypeProviders[strings.ToLower(*providerFlag)]
	if !ok {
		fmt.Println("Unknown or unspecified provider, use -provider")
		flag.Usage()
		return
	}

	if *printAllFlag {
		provider.PrintRegisteredTypes("")
		provider.PrintRegionalAvailability("", "")
	} else if *printTypeFlag == "all" {
		provider.PrintRegisteredTypes("")
	} else if *printTypeFlag != "" {
		provider.PrintRegisteredTypes(*printTypeFlag)
	} else if (*printRegionFlag != "" && *printZoneFlag != "") ||
		(*printRegionFlag != "" && *printZoneFlag == "") ||
		(*printRegionFlag == "all" && *printZoneFlag == "") {
		provider.PrintRegionalAvailability(*printRegionFlag, *printZoneFlag)
	} else if *generateFlag {
		err := provider.GenerateTypes()
		if err != nil {
			panic(err)
		}
	} else {
		flag.Usage()
	}
}
