##@ Database migrations

MIGRATION_NAME ?= unnamed

.PHONY: migrate
migrate: ## Run database migration
	$(GO) run ./cmd/pbmigrate

.PHONY: purgedb
purgedb: ## Delete database (dangerous!)
	$(GO) run ./cmd/pbmigrate purgedb

.PHONY: generate-migration
MIGRATION_NAME?=unnamed
generate-migration: ## Generate new migration file, use MIGRATION_NAME=name
	$(TERN) new -m internal/db/migrations $(MIGRATION_NAME)

