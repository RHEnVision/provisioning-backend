# Make documentation
```

Usage:
  make <target>

HTTP Clients
  update-clients        Update HTTP client stubs from upstream git repos
  validate-clients      Compare generated client code with git

Code quality
  fmt                   Format the project using `go fmt`
  lint                  Run Go language linter `golangci-lint`

Building
  build                 Build all binaries
  pbapi                 Build backend API service
  pbworker              Build worker service
  pbmigrate             Build migration command
  strip                 Strip debug information
  run-go                Run backend API using `go run`
  run                   Build and run backend API
  clean                 Clean build artifacts

Image building
  build-podman          Build container image using Podman

Database migrations
  migrate               Run database migration
  purgedb               Delete database (dangerous!)
  generate-migration    Generate new migration file, use MIGRATION_NAME=name

Go modules
  tidy-deps             Cleanup Go modules
  download-deps         Download Go modules
  update-deps           Update Go modules to latest versions
  help                  Print out the help content

OpenAPI
  generate-spec         Generate OpenAPI spec
  validate-spec         Compare OpenAPI spec with git

Testing
  test                  Run unit tests
  integration-test      Run integration tests (require database)

Go commands
  install-tools         Install required Go commands

Instance types
  generate-azure-types  Generate instance types for Azure
  generate-types        Generate instance types for all providers
```
