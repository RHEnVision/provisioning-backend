##@ Database migrations

.PHONY: migrate
migrate: ## Run database migration
	go run ./cmd/pbmigrate

.PHONY: purgedb
purgedb: ## Delete database (dangerous!)
	go run ./cmd/pbmigrate purgedb

.PHONY: generate-migration
MIGRATION_NAME?=unnamed
generate-migration: ## Generate new migration file, use MIGRATION_NAME=name
	migrate create -ext sql -dir internal/db/migrations -seq -digits 3 $(MIGRATION_NAME)

