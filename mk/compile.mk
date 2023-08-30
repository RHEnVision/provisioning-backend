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
run-worker: check-go ## Run worker using `go run`
	$(GO) run ./cmd/pbackend worker

.PHONY: run-statuser
run-statuser: check-go ## Run statuser using `go run`
	$(GO) run ./cmd/pbackend statuser

.PHONY: run-stats
run-stats: check-go ## Run stats using `go run`
	$(GO) run ./cmd/pbackend stats

.PHONY: cd-api
cd-api: check-go ## Run backend API using compile daemon
	$(CDAEMON) -build='go build -buildvcs=false -o ./pbackend-cd-api ./cmd/pbackend' -command='./pbackend-cd-api api'

.PHONY: cd-worker
cd-worker: check-go ## Run worker using compile daemon
	$(CDAEMON) -build='go build -buildvcs=false -o ./pbackend-cd-worker ./cmd/pbackend' -command='./pbackend-cd-worker worker'

.PHONY: cd-statuser
cd-statuser: check-go ## Run statuser using compile daemon
	$(CDAEMON) -build='go build -buildvcs=false -o ./pbackend-cd-statuser ./cmd/pbackend' -command='./pbackend-cd-statuser statuser'

.PHONY: cd-stats
cd-stats: check-go ## Run stats using compile daemon
	$(CDAEMON) -build='go build -buildvcs=false -o ./pbackend-cd-stats ./cmd/pbackend' -command='./pbackend-cd-stats stats'

CMD?=version
.PHONY: run
run: pbackend ## Build and run backend API
	./pbackend $(CMD)

.PHONY: clean
clean: ## Clean build artifacts and cache
	-rm pb*
	$(GO) clean -cache -modcache -testcache -fuzzcache
	GOROOT=$(GOROOT) GOCACHE=$(GOCACHE) $(GOLINT) cache clean
