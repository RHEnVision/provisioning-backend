##@ Building

SRC_GO := $(shell find . -name \*.go -print)
SRC_SQL := $(shell find . -name \*.sql -print)
SRC_YAML := $(shell find . -name \*.yaml -print)

PACKAGE_BASE = github.com/RHEnVision/provisioning-backend/internal
LDFLAGS = "-X $(PACKAGE_BASE)/version.BuildCommit=$(shell git rev-parse --short HEAD) -X $(PACKAGE_BASE)/version.BuildTime=$(shell date +'%Y-%m-%d_%T')"

build: pbapi pbmigrate pbworker ## Build all binaries

pbapi: $(SRC_GO) $(SRC_SQL) $(SRC_YAML) ## Build backend API service
	CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -o pbapi ./cmd/pbapi

pbworker: $(SRC_GO) $(SRC_SQL) $(SRC_YAML) ## Build worker service
	CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -o pbworker ./cmd/pbworker

pbmigrate: $(SRC_GO) $(SRC_SQL) $(SRC_YAML) ## Build migration command
	CGO_ENABLED=0 go build -o pbmigrate ./cmd/pbmigrate

.PHONY: strip
strip: build ## Strip debug information
	strip pbapi pbworker pbmigrate

.PHONY: run-go
run-go: ## Run backend API using `go run`
	go run ./cmd/pbapi

.PHONY: run
run: pbapi ## Build and run backend API
	./pbapi

.PHONY: clean
clean: ## Clean build artifacts
	rm pbapi pbmigrate pbworker
