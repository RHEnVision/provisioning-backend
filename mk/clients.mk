##@ HTTP Clients

.PHONY: update-clients
update-clients: ## Update OpenAPI specs from upstream
	curl -s -o ./config/ib_api.yaml -z ./config/ib_api.yaml https://raw.githubusercontent.com/osbuild/image-builder/main/internal/v1/api.yaml
	curl -s -o ./config/sources_api.json -z ./config/sources_api.json https://raw.githubusercontent.com/RedHatInsights/sources-api-go/main/public/openapi-3-v3.1.json
	curl -s -o ./config/rbac_api.json -z ./config/rbac_api.json https://raw.githubusercontent.com/RedHatInsights/insights-rbac/master/docs/source/specs/openapi.json

generate-clients: internal/clients/http/image_builder/client.gen.go internal/clients/http/sources/client.gen.go internal/clients/http/rbac/client.gen.go ## Generate HTTP client stubs

internal/clients/http/sources/client.gen.go: config/sources_config.yml config/sources_api.json
	$(OAPICODEGEN) -config ./config/sources_config.yml ./config/sources_api.json

internal/clients/http/image_builder/client.gen.go: config/ib_config.yaml config/ib_api.yaml
	$(OAPICODEGEN) -config ./config/ib_config.yaml ./config/ib_api.yaml

internal/clients/http/rbac/client.gen.go: config/rbac_config.yml config/rbac_api.json
	$(OAPICODEGEN) -config ./config/rbac_config.yml ./config/rbac_api.json

.PHONY: validate-clients
validate-clients: generate-clients ## Compare generated client code with git
	git diff --exit-code internal/clients/*/client.gen.go
