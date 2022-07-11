TEST_TAGS?=test

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

.PHONY: clean
clean:
	rm pbapi pbmigrate

.PHONY: build-podman
build-podman:
	# remote podman build has problem with -f option, so we link the file as workaround
	ln -f build/Dockerfile Containerfile
	podman build --build-arg quay_expiration=2d -t provisioning-backend .

.PHONY: prep
prep:
	go mod download

.PHONY: run
run:
	go run ./cmd/pbapi

.PHONY: update-clients
update-clients:
	wget -O ./configs/ib_api.yaml -qN https://raw.githubusercontent.com/osbuild/image-builder/main/internal/v1/api.yaml
	wget -O ./configs/sources_api.json -qN https://raw.githubusercontent.com/RedHatInsights/sources-api-go/main/public/openapi-3-v3.1.json

generate-clients: internal/clients/image_builder/client.gen.go internal/clients/sources/client.gen.go

internal/clients/sources/client.gen.go: configs/sources_config.yml configs/sources_api.json
	oapi-codegen -config ./configs/sources_config.yml ./configs/sources_api.json

internal/clients/image_builder/client.gen.go: configs/ib_config.yaml configs/ib_api.yaml
	oapi-codegen -config ./configs/ib_config.yaml ./configs/ib_api.yaml

.PHONY: validate-clients
validate-clients: generate-clients
	git diff --exit-code internal/clients/*/client.gen.go

.PHONY: install-tools
install-tools:
	go install golang.org/x/tools/cmd/goimports@latest
	# pin for a bug: https://github.com/golangci/golangci-lint/issues/2851
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.2
	go install github.com/jackc/tern@latest
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

.PHONY: fmt
fmt:
	go fmt ./...
	goimports -w .

.PHONY: lint
lint:
	golangci-lint run

.PHONY: migrate
migrate:
	go run ./cmd/pbmigrate

.PHONY: purgedb
purgedb:
	go run ./cmd/pbmigrate purgedb

.PHONY: test
test:
	go test -tags=$(TEST_TAGS) ./...

.PHONY: generate-migration
MIGRATION_NAME?=unnamed
generate-migration:
	migrate create -ext sql -dir internal/db/migrations -seq -digits 3 $(MIGRATION_NAME)

.PHONY: generate-spec
generate-spec:
	go run ./cmd/spec

.PHONY: validate-spec
validate-spec: generate-spec
	git diff --exit-code api/openapi.gen.json

.PHONY: validate
validate: validate-spec validate-clients
