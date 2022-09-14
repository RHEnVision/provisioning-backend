##@ Go commands

.PHONY: install-tools
install-tools: ## Install required Go commands
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/jackc/tern@latest
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
	go install mvdan.cc/gofumpt@latest

