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
	$(GO) install github.com/lzap/gobump@latest
	$(GO) run github.com/lzap/gobump@latest
	$(GO) mod tidy

# aliases
.PHONY: prep
prep: download-deps

.PHONY: tidy
tidy: tidy-deps

.PHONY: bump
bump: update-deps
