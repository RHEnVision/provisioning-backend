#!/bin/bash

# --------------------------------------------
# Options that must be configured by app owner
# --------------------------------------------
APP_NAME="provisioning"  # name of app-sre "application" folder this component lives in
COMPONENT_NAME="provisioning-backend"  # name of resourceTemplate component for deploy
IMAGE="quay.io/cloudservices/provisioning-backend"  # image location on quay
DOCKERFILE="build/Dockerfile"

IQE_PLUGINS="provisioning"  # name of the IQE plugin for this app.
IQE_MARKER_EXPRESSION="api and smoke"  # This is the value passed to pytest -m
IQE_FILTER_EXPRESSION=""  # This is the value passed to pytest -k
IQE_CJI_TIMEOUT="30m"  # This is the time to wait for smoke test to complete or fail

#EXTRA_DEPLOY_ARGS=

printenv

# Install bonfire repo/initialize
# https://raw.githubusercontent.com/RedHatInsights/bonfire/master/cicd/bootstrap.sh
# This script automates the install / config of bonfire
CICD_URL=https://raw.githubusercontent.com/RedHatInsights/bonfire/master/cicd
curl -s $CICD_URL/bootstrap.sh > .cicd_bootstrap.sh && source .cicd_bootstrap.sh

# This script is used to build the image that is used in the PR Check
source $CICD_ROOT/build.sh

# This script is used to deploy the ephemeral environment for smoke tests.
# The manual steps for this can be found in:
# https://internal.cloud.redhat.com/docs/devprod/ephemeral/02-deploying/
source $CICD_ROOT/deploy_ephemeral_env.sh


# ADD the image stubs
oc project $NAMESPACE

orgID="0369233"
accountID="3340851"

dbPod=$(oc get pods -o custom-columns=POD:.metadata.name --no-headers | grep 'image-builder-db')

# AWS stub
imageID="ad6d6a99-7a95-4b74-8760-cce28df35bda" # created on 2023-08-23
imageName="pipeline-aws"

composeRequestAWSJson='{"image_name": "'$imageName'", "distribution": "rhel-92", "customizations": {}, "image_requests": [{"image_type": "aws", "architecture": "x86_64", "upload_request": {"type": "aws", "options": {"share_with_accounts": ["093942615996"]}}}]}'
oc exec $dbPod -- psql -d image-builder -c "INSERT INTO public.composes (job_id, request, created_at, account_number, org_id, image_name, deleted) VALUES
('$imageID', '$composeRequestAWSJson', '$(date +"%Y-%m-%d %T")', '$orgID', '$accountID', '$imageName', false);"

# GCP stub
imageID="6c79ab2c-176d-4e24-80ca-bd9f3a019ca6" # created on 2023-08-23
imageName="pipeline-gcp"

composeRequestGCPJson='{"image_name": "'$imageName'", "distribution": "rhel-92", "customizations": {}, "image_requests": [{"image_type": "gcp", "architecture": "x86_64", "upload_request": {"type": "gcp", "options": {"share_with_accounts": ["user:oezr@redhat.com"]}}}]}'
oc exec $dbPod -- psql -d image-builder -c "INSERT INTO public.composes (job_id, request, created_at, account_number, org_id, image_name, deleted) VALUES
('$imageID', '$composeRequestGCPJson', '$(date +"%Y-%m-%d %T")', '$orgID', '$accountID', '$imageName', false);"

# Run smoke tests using a ClowdJobInvocation and iqe-tests
source $CICD_ROOT/cji_smoke_test.sh

# Post a comment with test run IDs to the PR
source $CICD_ROOT/post_test_results.sh

SONAR_PR_CHECK="true" # used by sonarqube to set PR check variables
source $WORKSPACE/.rhcicd/sonarqube.sh
