.PHONY: build
build: build-pbapi build-pbmigrate

.PHONY: build-pbapi
build-pbapi:
	CGO_ENABLED=0 go build -o pbapi ./cmd/pbapi

.PHONY: build-pbmigrate
build-pbmigrate:
	CGO_ENABLED=0 go build -o pbmigrate ./cmd/pbmigrate

.PHONY: strip
strip: build
	strip pbapi pbmigrate

.PHONY: build-podman
build-podman:
	podman build .

.PHONY: prep
prep:
	go mod download

.PHONY: run
run:
	go run ./cmd/pbapi

.PHONY: install-tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: lint
lint:
	go fmt ./...
	go vet ./...
	golangci-lint run

.PHONY: migrate
migrate: build-pbmigrate
	pbmigrate

.PHONY: generate-migration
MIGRATION_NAME?=unnamed
generate-migration:
	migrate create -ext sql -dir internal/db/migrations -seq -digits 3 $(MIGRATION_NAME)
