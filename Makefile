.PHONY: build
build:
	go build ./cmd/pbapi

.PHONY: run
run:
	go run ./cmd/pbapi
