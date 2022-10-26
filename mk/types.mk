##@ Instance types

.PHONY: generate-azure-types
generate-azure-types: ## Generate instance types for Azure
	go run cmd/typesctl/main.go -provider azure -generate

.PHONY: generate-ec2-types
generate-ec2-types: ## Generate instance types for EC2
	go run cmd/typesctl/main.go -provider ec2 -generate

.PHONY: generate-gcp-types
generate-gcp-types: ## Generate instance types for GCP
	go run cmd/typesctl/main.go -provider gcp -generate

.PHONY: generate-types
generate-types: generate-ec2-types generate-azure-types generate-gcp-types ## Generate instance types for all providers
