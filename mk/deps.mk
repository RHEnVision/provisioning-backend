##@ Go modules

.PHONY: tidy-deps
tidy-deps: ## Cleanup Go modules
	go mod tidy

.PHONY: download-deps
download-deps: ## Download Go modules
	go mod download

.PHONY: list-mods
list-mods: ## List application modules
	go list ./...

.PHONY: list-deps
list-deps: ## List dependencies and their versions
	go list -m -u all

.PHONY: update-deps
update-deps: ## Update Go modules to latest versions
	go get -u ./...
	# Needs Go 1.19: https://github.com/jackc/puddle/issues/26
	go get github.com/jackc/puddle/v2@v2.0.0
	# We rely on a feature in unreleased version
	go get github.com/jackc/tern/v2@v2.0.0-beta.3
	go mod tidy

# alias for download-deps
.PHONY: prep
prep: download-deps
