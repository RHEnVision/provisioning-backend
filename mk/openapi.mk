##@ OpenAPI

.PHONY: generate-spec
generate-spec: ## Generate OpenAPI spec
	go run ./cmd/spec

.PHONY: validate-spec
validate-spec: generate-spec ## Compare OpenAPI spec with git
	git diff --exit-code api/openapi.gen.json

