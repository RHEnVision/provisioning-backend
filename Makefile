# Copyright Red Hat

check: check-copyright

check-copyright:
	@build/check-copyright.sh

build: 
	go build -o provisioning-api main.go

run:
	go run main.go

vet:
	go vet ./...

staticcheck: 
	staticcheck ./...

lint: vet staticcheck

# bonfire-config-local:
# 	@cp default_config.yaml.local.example config.yaml
# 	@sed -i ${OS_SED} 's|REPO|$(PWD)|g' config.yaml

# bonfire-config-github:
# 	@cp default_config.yaml.github.example config.yaml
