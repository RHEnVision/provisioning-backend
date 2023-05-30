##@ OpenAPI

.PHONY: generate-spec
generate-spec: ## Generate OpenAPI spec
	$(GO) run ./cmd/spec $(GIT_TAG)
	echo $(GIT_TAG) > ./cmd/spec/VERSION

.PHONY: validate-spec
validate-spec: ## Compare OpenAPI spec with git
	$(GO) run ./cmd/spec $(shell head -n1 ./cmd/spec/VERSION)
	git diff --exit-code api/openapi.gen.json

