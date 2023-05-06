##@ Go modules

.PHONY: tidy-deps
tidy-deps: ## Cleanup Go modules
	$(GO) mod tidy

.PHONY: download-deps
download-deps: ## Download Go modules
	$(GO) mod download

.PHONY: list-mods
list-mods: ## List application modules
	$(GO) list ./...

.PHONY: list-deps
list-deps: ## List dependencies and their versions
	$(GO) list -m -u all

.PHONY: update-deps
update-deps: ## Update Go modules to latest versions
	$(GO) get -u ./...
	# Needs Go 1.19: https://github.com/jackc/puddle/issues/26
	$(GO) get github.com/jackc/puddle/v2@v2.0.0
	# Compile errors with Go 1.18
	$(GO) get go.opentelemetry.io/contrib@v1.15.0
	$(GO) get go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.40.0
	$(GO) get go.opentelemetry.io/otel@v1.14.0
	$(GO) get go.opentelemetry.io/otel/exporters/jaeger@v1.14.0
	$(GO) get go.opentelemetry.io/otel/metric@v0.37.0
	$(GO) get go.opentelemetry.io/otel/sdk@v1.14.0
	$(GO) get go.opentelemetry.io/otel/trace@v1.14.0
	$(GO) mod tidy

# alias for download-deps
.PHONY: prep
prep: download-deps
