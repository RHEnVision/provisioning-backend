##@ Go commands

.PHONY: install-tools
install-tools: ## Install required Go commands
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1
	go install github.com/jackc/tern@latest
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
	go install mvdan.cc/gofumpt@latest

.PHONY: update-tools
update-tools: ## Update required Go commands
	go get -u golang.org/x/tools/cmd/goimports github.com/golangci/golangci-lint/cmd/golangci-lint github.com/jackc/tern github.com/deepmap/oapi-codegen/cmd/oapi-codegen mvdan.cc/gofumpt

.PHONY: generate-changelog
generate-changelog: ## Generate CHANGELOG.md from git history
	@echo "Required tool: pip3 install -e https://github.com/RHEnVision/changelog"
	python3 -m rhenvision_changelog .
