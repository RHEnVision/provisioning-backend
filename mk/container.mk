##@ Image building

CONTAINER_BUILD_OPTS ?= --build-arg=quay_expiration=2d
CONTAINER_IMAGE ?= provisioning-backend

.PHONY: build-podman
build-podman: ## Build container image using Podman
	# remote podman build has problem with -f option, so we link the file as workaround
	ln -sf build/Dockerfile Containerfile
	podman build $(CONTAINER_BUILD_OPTS) -t $(CONTAINER_IMAGE) .
