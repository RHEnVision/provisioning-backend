package main

import (
	"fmt"
	"os"
	"runtime"

	// DAO implementation, must be initialized before any database packages
	_ "github.com/RHEnVision/provisioning-backend/internal/dao/pgx"

	// HTTP client implementations
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/azure"
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/ec2"
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/gcp"
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/image_builder"
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/sources"

	"github.com/RHEnVision/provisioning-backend/internal/random"
	"github.com/RHEnVision/provisioning-backend/internal/version"
)

func init() {
	random.SeedGlobal()
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	switch os.Args[1] {
	case "api":
		api()
	case "worker":
		worker()
	case "migrate":
		migrate()
	case "statuser":
		statuser()
	case "version":
		ver()
	default:
		usage()
	}
}

func usage() {
	fmt.Println("Usage: pbackend [migrate|api|worker|statuser|version]")
	os.Exit(1)
}

func ver() {
	fmt.Printf("Version %s, Go version %s, built %s\n", version.BuildCommit, runtime.Version(), version.BuildTime)
}
