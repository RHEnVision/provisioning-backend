##@ HTTP Clients

.PHONY: update-clients
update-clients: ## Update OpenAPI specs from upstream
	wget -O ./config/ib_api.yaml -qN https://raw.githubusercontent.com/osbuild/image-builder/main/internal/v1/api.yaml
	wget -O ./config/sources_api.json -qN https://raw.githubusercontent.com/RedHatInsights/sources-api-go/main/public/openapi-3-v3.1.json

generate-clients: internal/clients/http/image_builder/client.gen.go internal/clients/http/sources/client.gen.go ## Generate HTTP client stubs

internal/clients/http/sources/client.gen.go: config/sources_config.yml config/sources_api.json
	$(OAPICODEGEN) -config ./config/sources_config.yml ./config/sources_api.json

internal/clients/http/image_builder/client.gen.go: config/ib_config.yaml config/ib_api.yaml
	$(OAPICODEGEN) -config ./config/ib_config.yaml ./config/ib_api.yaml

.PHONY: validate-clients
validate-clients: generate-clients ## Compare generated client code with git
	git diff --exit-code internal/clients/*/client.gen.go
