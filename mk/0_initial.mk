#
# Initial file (included as the first)
#

PROJECT_DIR:=$(shell dirname $(abspath $(firstword $(MAKEFILE_LIST))))

# Update GitHub Workflows when changing this
GO_VERSION?=1.19.6
GOBIN:=$(PROJECT_DIR)/bin
GO?=$(shell go env GOPATH)/bin/go$(GO_VERSION)

GOLINT?=$(GOBIN)/golangci-lint
GOFUMPT?=$(GOBIN)/gofumpt
GOIMPORTS?=$(GOBIN)/goimports
OAPICODEGEN?=$(GOBIN)/oapi-codegen
TERN?=$(GOBIN)/tern

.PHONY: check-go
check-go:
	@test -x $(GO) || test "$(GO)" = go || (echo "Go $(GO_VERSION) not installed, run: make install-go install-tools" && exit 1)

.PHONY: check-system-go
check-system-go:
	@go version | grep $(GO_VERSION) >/dev/null || (echo "System Go version does not match required Go: $(GO_VERSION)" && exit 1)
