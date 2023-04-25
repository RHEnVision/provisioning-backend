#!/bin/bash
set -e

IMAGE="${IMAGE:-quay.io/cloudservices/provisioning-backend}"
IMAGE_TAG="$(git rev-parse --short=7 HEAD)"
SMOKE_TEST_TAG="latest"

if [[ -z "$QUAY_USER" || -z "$QUAY_TOKEN" ]]; then
    echo "QUAY_USER and QUAY_TOKEN must be set"
    exit 1
fi

if [[ -z "$RH_REGISTRY_USER" || -z "$RH_REGISTRY_TOKEN" ]]; then
    echo "RH_REGISTRY_USER and RH_REGISTRY_TOKEN  must be set"
    exit 1
fi

# Login to quay
podman login \
    -u ${QUAY_USER} \
    -p ${QUAY_TOKEN} \
    quay.io

# Login to registry.redhat
podman login \
    -u ${RH_REGISTRY_USER} \
    -p ${RH_REGISTRY_TOKEN} \
    registry.redhat.io

# build and push
make build-podman \
    CONTAINER_IMAGE="${IMAGE}:${IMAGE_TAG}"

# push to logged in registries and tag for SHA
podman tag "${IMAGE}:${IMAGE_TAG}" "${IMAGE}:${SMOKE_TEST_TAG}"
podman push "${IMAGE}:${IMAGE_TAG}"
podman push "${IMAGE}:${SMOKE_TEST_TAG}"

# SONAR_PR_CHECK="false" # used by sonarqube to not set PR check variables
# source $WORKSPACE/.rhcicd/sonarqube.sh
