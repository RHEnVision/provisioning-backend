##@ HTTP Clients

.PHONY: update-clients
update-clients: ## Update HTTP client stubs from upstream git repos
	wget -O ./configs/ib_api.yaml -qN https://raw.githubusercontent.com/osbuild/image-builder/main/internal/v1/api.yaml
	wget -O ./configs/sources_api.json -qN https://raw.githubusercontent.com/RedHatInsights/sources-api-go/main/public/openapi-3-v3.1.json

generate-clients: internal/clients/http/image_builder/client.gen.go internal/clients/http/sources/client.gen.go

internal/clients/http/sources/client.gen.go: configs/sources_config.yml configs/sources_api.json
	oapi-codegen -config ./configs/sources_config.yml ./configs/sources_api.json

internal/clients/http/image_builder/client.gen.go: configs/ib_config.yaml configs/ib_api.yaml
	oapi-codegen -config ./configs/ib_config.yaml ./configs/ib_api.yaml

.PHONY: validate-clients
validate-clients: generate-clients ## Compare generated client code with git
	git diff --exit-code internal/clients/*/client.gen.go

