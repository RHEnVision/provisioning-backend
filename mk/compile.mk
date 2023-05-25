##@ Building

SRC_GO := $(shell find . -name \*.go -print)
SRC_SQL := $(shell find . -name \*.sql -print)
SRC_YAML := $(shell find . -name \*.yaml -print)

build: pbackend ## Build all binaries

all-deps: $(SRC_GO) $(SRC_SQL) $(SRC_YAML)

pbackend: check-go all-deps ## Build backend
	CGO_ENABLED=0 $(GO) build -o pbackend ./cmd/pbackend

.PHONY: strip
strip: build ## Strip debug information
	strip pbackend

.PHONY: run-api
run-api: check-go ## Run backend API using `go run`
	$(GO) run ./cmd/pbackend api

.PHONY: run-worker
run-worker: check-go ## Run backend API using `go run`
	$(GO) run ./cmd/pbackend worker

.PHONY: run-statuser
run-statuser: check-go ## Run backend API using `go run`
	$(GO) run ./cmd/pbackend statuser

CMD?=version
.PHONY: run
run: pbackend ## Build and run backend API
	./pbackend $(CMD)

.PHONY: clean
clean: ## Clean build artifacts and cache
	-rm pb*
	$(GO) clean -cache -modcache -testcache -fuzzcache
	GOROOT=$(GOROOT) GOCACHE=$(GOCACHE) $(GOLINT) cache clean
