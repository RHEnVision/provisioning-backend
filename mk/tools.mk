##@ Go commands

.PHONY: install-tools
install-tools: ## Install required Go commands
	go install golang.org/x/tools/cmd/goimports@latest
	# pin for a bug: https://github.com/golangci/golangci-lint/issues/2851
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.2
	go install github.com/jackc/tern@latest
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

