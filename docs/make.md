# Make documentation
```

Usage:
  make <target>

Database migrations
  migrate               Run database migration
  purgedb               Delete database (dangerous!)
  generate-migration    Generate new migration file, use MIGRATION_NAME=name

Code quality
  format                Format Go source code using `go fmt`
  imports               Rearrange imports using `goimports`
  lint                  Run Go language linter `golangci-lint`
  check-migrations      Check migration files for changes
  check-commits         Check commit format

Building
  build                 Build all binaries
  pbackend              Build backend
  strip                 Strip debug information
  run-api               Run backend API using `go run`
  run-worker            Run backend API using `go run`
  run-statuser          Run backend API using `go run`
  run                   Build and run backend API
  clean                 Clean build artifacts and cache

Help
  help                  Print out the help content
  generate-help-doc     Generate 'make help' markdown in docs/
  validate-help-doc     Compare example configuration
  generate-example-config  Generate example configuration
  validate-example-config  Compare example configuration

Dashboard
  generate-dashboard    Generate dashboard
  validate-dashboard    Compare dashboard configmaps with git

Go modules
  tidy-deps             Cleanup Go modules
  download-deps         Download Go modules
  list-mods             List application modules
  list-deps             List dependencies and their versions
  update-deps           Update Go modules to latest versions

Go commands
  install-go            Install required Go version
  install-tools         Install required Go commands into ./bin
  update-tools          Update required Go commands
  generate-changelog    Generate CHANGELOG.md from git history

Testing
  test                  Run unit tests
  integration-test      Run integration tests (require database)

OpenAPI
  generate-spec         Generate OpenAPI spec
  validate-spec         Compare OpenAPI spec with git

Instance types
  generate-azure-types  Generate instance types for Azure
  generate-ec2-types    Generate instance types for EC2
  generate-gcp-types    Generate instance types for GCP
  generate-types        Generate instance types for all providers

HTTP Clients
  update-clients        Update OpenAPI specs from upstream
  generate-clients      Generate HTTP client stubs
  validate-clients      Compare generated client code with git

Image building
  build-podman          Build container image using Podman
```
