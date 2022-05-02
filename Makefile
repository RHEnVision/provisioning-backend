.PHONY: build
build:
	go build ./cmd/pbapi

.PHONY: run
run:
	go run ./cmd/pbapi

.PHONY: models
models: sqlboiler.toml
	sqlboiler sqlite3 --wipe -o internal/models

.PHONY: migrate
migrate:
	sqlite3 devel.db < cmd/pbapi-migrate/schema.sql

