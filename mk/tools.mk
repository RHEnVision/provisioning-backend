##@ Go commands

.PHONY: install-go
install-go: ## Install required Go version
	go install golang.org/dl/go$(GO_VERSION)@latest
	$(GO) download

# Update GitHub Workflows when changing this
.PHONY: install-tools
install-tools: ## Install required Go commands into ./bin
	GOBIN=$(GOBIN) $(GO) install golang.org/x/tools/cmd/goimports@v0.8.0
	GOBIN=$(GOBIN) $(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1
	GOBIN=$(GOBIN) $(GO) install github.com/jackc/tern/v2@v2.1.0
	GOBIN=$(GOBIN) $(GO) install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4
	GOBIN=$(GOBIN) $(GO) install mvdan.cc/gofumpt@v0.5.0

.PHONY: update-tools
update-tools: ## Update required Go commands
	$(GO) get -u golang.org/x/tools/cmd/goimports github.com/golangci/golangci-lint/cmd/golangci-lint github.com/jackc/tern github.com/deepmap/oapi-codegen/cmd/oapi-codegen mvdan.cc/gofumpt

.PHONY: generate-changelog
generate-changelog: ## Generate CHANGELOG.md from git history
	@echo "Required tool: pip3 install -e https://github.com/RHEnVision/changelog"
	python3 -m rhenvision_changelog .
