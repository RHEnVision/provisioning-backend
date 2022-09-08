##@ Code quality

.PHONY: fmt
fmt: ## Format the project using `go fmt`
	go fmt ./...
	goimports -w .

.PHONY: lint
lint: ## Run Go language linter `golangci-lint`
	golangci-lint run

