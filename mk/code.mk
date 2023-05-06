##@ Code quality

.PHONY: format-fmt
format-fmt:
	$(GOFMT) -l -w -s .

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
	$(GOLINT) run

.PHONY: check-migrations
check-migrations: ## Check migration files for changes
	@scripts/check_migrations.sh

.PHONY: check-commits
check-commits: ## Check commit format
	python -m rhenvision_changelog commit-check

.PHONY: fmt ## Alias to perform all code formatting and linting
fmt: format imports lint

.PHONY: check ## Alias to perform all checking (commits, migrations)
check: check-commits check-migrations

