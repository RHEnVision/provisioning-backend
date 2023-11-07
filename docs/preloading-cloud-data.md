# Preloading cloud data

Some data does not change often and for the best possible user experience, we preload it and embed the information into the application itself. This includes:

* List of region/location/zone names
* Instance type names
* Common instance type details (vCPUs, cores, memory, local drive)
* Specific type details (VM generation for Azure)
* Supported flag (when type meets Minimum RHEL Requirements criteria)

The data is stored as [YAML](./../internal/preload) and available through a `preload` Go package and REST API endpoints.

## Model

The terminology is different for each hyperscaler, however, the data model uses the AWS terminology: regions and (availability) zones:

* AWS EC2: zones (`us-east-1`) and availability zones (`a`), fully qualified name is `us-east-1a`
* Azure: locations (`northeurope`) and availability zones (`1`), fully qualified name is `northeurope_1`
* GCP: zones (`us-west`) and regions (`b`), fully qualified name is `us-west-b`

### Registered Instance Types

The `InstanceType` type contains common fields like the number of vCPUs and cores, total amount of memory, ephemeral storage, and supported flag. It also has an architecture, therefore, a single instance type cannot have multiple architectures. In this model, every instance type with a different architecture has a different instance type name. This applies to all major hyperscalers except Azure where Intel 64bit and 32bit are shared. However, 32bit is no longer supported by RHEL and therefore this is not an issue.

The only supported architectures at the moment are: x86_64 and arm64. All other instance types are ignored. InstanceType also contains optional details, currently for Azure the VM generation.

Functions in the `preload` package can be used to find particular `InstanceType` by name.

### Regional Type Availability

Depending on hyperscaler provider, different instance types can be available in different regions. For this reason, a function in the `preload` package will return a slice of instance types for every zone and region pair. The data is stored in individual files as [YAML](./../internal/preload) with the following naming convention:

* AWS EC2: `region.yaml` (all availability zones have the same instance types)
* Azure: `location_zone.yaml` (different instance types in zones)
* GCP: `zone.yaml` (all)

## Refreshing data

In order to update the embedded data with the latest available instance types, use provided `typesctl` utility and Makefile targets. First, configure connections:

```
# cat config/api.env
AWS_KEY=KEY
AWS_SECRET=SECRET
AZURE_TENANT_ID=UUID
AZURE_SUBSCRIPTION_ID=UUID
AZURE_CLIENT_ID=UUID
AZURE_CLIENT_SECRET=SECRET
GCP_JSON=xxxxxxx
GCP_PROJECT_ID=myproject-12345
```

Important: AWS requires special attention: global endpoint session setting must be set for an account.

1. Sign in as a root user or a user with permissions to perform IAM administration tasks. To change the compatibility of session tokens, you must have a policy that allows the iam:SetSecurityTokenServicePreferences action.
2. Open the IAM console. In the navigation pane, choose Account settings.
3. Under Security Token Service (STS) section Session Tokens from the STS endpoints. The Global endpoint indicates Valid only in AWS Regions enabled by default. Choose Change.
4. In the Change region compatibility dialog box, select All AWS Regions. Then choose Save changes.

All regions must be activated in order to refresh data. To activate or deactivate AWS STS in a Region that is enabled by default (console)

1. Sign in as a root user or a user with permissions to perform IAM administration tasks.
2. Open the IAM console and in the navigation pane choose Account settings.
3. In the Security Token Service (STS) section Endpoints, find the Region that you want to configure, and then choose Active or Inactive in the STS status column.
4. In the dialog box that opens, choose Activate or Deactivate.

For more info read:

* https://aws.amazon.com/premiumsupport/knowledge-center/iam-validate-access-credentials/
* https://docs.aws.amazon.com/general/latest/gr/rande-manage.html
* https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_enable-regions.html#sts-regions-manage-tokens


If these steps are not done, the operation will throw AuthFailure "AWS was not able to validate the provided access credentials". Our shared (team) account, unfortunately, will not work. The solution is to use a personal account. You can create a `typesctl.env` config file which will override `api.env`:

```
# cat config/typesctl.env
AWS_KEY=PERSONAL_ACCOUNT_KEY
AWS_SECRET=PERSONAL_ACCOUNT_SECRET
```

Make sure you either use the root account, or the role has the necessary permissions (`DescribeRegions` and others). You can find these permissions in the [AWS configuration document](./configure-amazon-role.md).

To refresh data:

```
make generate-azure-types
make generate-ec2-types
make generate-gcp-types
```

Or to do this all at once:

```
make generate-types
```

## Pushing data to git

Make sure to refresh the data in separate commits or PRs. These changesets can be long and hard to read, so make sure this is not part of other code changes.
