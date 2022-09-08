##@ Instance types

.PHONY: generate-azure-types
generate-azure-types: ## Generate instance types for Azure
	go run cmd/typesctl/main.go -provider azure -generate

.PHONY: generate-types
generate-types: generate-azure-types ## Generate instance types for all providers
