##@ Code quality

.PHONY: format-fmt
format-fmt:
	$(GO) fmt ./...

.PHONY: format-fumpt
format-fumpt:
	$(GOFUMPT) -l -w .

.PHONY: format
format: format-fmt format-fumpt ## Format Go source code using `go fmt`

.PHONY: imports
imports: ## Rearrange imports using `goimports`
	$(GOIMPORTS) -w .

.PHONY: lint
lint: ## Run Go language linter `golangci-lint`
	GOROOT=$(GOROOT) PATH=$(GOROOT)/bin:${PATH} GOCACHE=$(GOCACHE) $(GOLINT) run

.PHONY: check-migrations
check-migrations: ## Check migration files for changes
	@scripts/check_migrations.sh

.PHONY: check-commits
check-commits: ## Check commit format
	python -m rhenvision_changelog commit-check

.PHONY: fmt ## Alias to perform all code formatting and linting
fmt: format imports lint

.PHONY: check-fmt ## Reformat the code and check git diff
check-fmt: format imports
	git diff --exit-code

.PHONY: check ## Alias to perform commit message and migration checking
check: check-commits check-migrations
