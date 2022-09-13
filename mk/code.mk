##@ Code quality

.PHONY: format
format: ## Format Go source code using `go fmt`
	go fmt ./...

.PHONY: imports
imports: ## Rearrange imports using `goimports`
	goimports -w .

.PHONY: lint
lint: ## Run Go language linter `golangci-lint`
	golangci-lint run

.PHONY: fmt
fmt: format imports lint

