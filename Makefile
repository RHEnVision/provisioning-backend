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
