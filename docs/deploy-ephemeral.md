# How to deploy to ephemeral environment

## Assumptions

1. You've got access to Ephemeral environment.
2. You've got [bonfire setuped-up](https://consoledot.pages.redhat.com/docs/dev/getting-started/ephemeral/install-bonfire.html)
3. You've bonfire in active venv
4. You've got `oc` cli tool
5. You've logged in to Ephemeral environment through `oc login` as described in [Installing and using Bonfire](https://consoledot.pages.redhat.com/docs/dev/getting-started/ephemeral/install-bonfire.html)


## Deploy API to Ephemeral

Deploy to ephemeral with dependencies also.

```
bonfire deploy provisioning sources
```

Note: you can also use bonfire option `--set-image-tag quay.io/cloudservices/provisioning-backend=pr-<#PR>-<shorthash>` to set different tag of your choice for api image.

## Deploy with to Ephemeral

```
bonfire deploy --frontends true provisioning sources
```

Note: you can also use bonfire option `--set-image-tag quay.io/cloudservices/provisioning-frontend=pr-<#PR>-<shorthash>` to set different tag of your choice for frontend image.
