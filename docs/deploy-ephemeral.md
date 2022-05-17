# How to deploy to ephemeral environment

## Assumptions

1. You've got access to Ephemeral environment.
2. You've got [bonfire setuped-up](https://consoledot.pages.redhat.com/docs/dev/getting-started/ephemeral/onboarding.html)
3. You've got `oc` cli tool
4. You've logged in to Ephemeral environment through `oc login` as described in [Onboarding to Ephemeral](https://consoledot.pages.redhat.com/docs/dev/getting-started/ephemeral/onboarding.html)
5. You've bonfire ideally in active venv
6. You've joined https://quay.io/organization/envision org with write access
7. You've run `podman login quay.io` to login into quay with your account

# Deploy to Ephemeral

1. Copy `deploy/bonfire.example.yaml` to `deploy/bonfire.yaml`
2. set `<path_to_service_dir>` to the path to the local path of this repo
3. run following

```
make build-podman
podman push provisioning-backend quay.io/envision/provisioning-backend:$(git rev-parse --short=7 HEAD)
bonfire deploy --source local --local-config-path ./deploy/bonfire.yaml provisioning
```

Note: you can also use bonfire option `--set-image-tag quay.io/envision/provisioning-backend=latest` to set different tag of your choice.
