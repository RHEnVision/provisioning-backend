# Configuring Azure deployment

The provisioning service uses Azure SDK to perform all provisioning tasks.
We follow the best possible security standard and practices.
In this article, we will describe how to create necessary service account setup and configure this application with it.

## Setup Azure service account

Lighthouse offering is concept how one Azure account offers (offering tenant) to another (target tenant, also customer) given solution.
Main concept to understand is [cross-tenant management](https://learn.microsoft.com/en-us/azure/lighthouse/concepts/cross-tenant-management-experience).
The offering can deploy anything that Azure Resource Manager allows.
In our use case we use only permission delegation.

The template we will prepare provides roles in target Tenant to the Principal from offering tenant.
The template assigns following roles to the offering Principal:

| Role name                                                                                                                                                                     | UUID                                 | Motivation                                                                                                                                                                                                               |
|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [Reader](https://learn.microsoft.com/en-us/azure/role-based-access-control/built-in-roles#reader)                                                                             | acdd72a7-3385-48ef-bd42-f606fba81ae7 | Allows us to fetch information about resources                                                                                                                                                                           |
| [Virtual Machine Contributor](https://learn.microsoft.com/en-us/azure/role-based-access-control/built-in-roles#virtual-machine-contributor)                                   | 9980e02c-c2be-4d73-94e8-173b1dc7cf3c | Allows to deploy virtual machines                                                                                                                                                                                        |
| [Contributor](https://learn.microsoft.com/en-us/azure/role-based-access-control/built-in-roles#contributor)                                                                   | b24988ac-6180-42a0-ab88-20f7382dd24c | Unfortunate wokaround for creating supporting resources, hopefully better solution will be found and this role will not be necessary in the future                                                                       |
| [Registration assignment Delete Role](https://learn.microsoft.com/en-us/azure/role-based-access-control/built-in-roles#managed-services-registration-assignment-delete-role)  | 91c1777a-f3dc-4fae-b103-61d183457e46 | Enables to remove the delegation when no longer needed. It is [recommended best practice](https://learn.microsoft.com/en-us/azure/lighthouse/concepts/tenants-users-roles#best-practices-for-defining-users-and-roles).  |


### Prepare Azure offering tenant

Here we will prepare the principal in offering tenant.
This Principal will get the permissions once Customer tenant accepts the offering.

- Go to Azure Active Directory
- Copy Tenant ID into config/api.env `AZURE_TENANT_ID`
- Go to App Registrations
- Select `New registration`
  - Name it example: provisioning-service
  - Select Single tenant account type
  - No Redirect URI is not required
- Go to details of this new registration
- Copy `Application (client) ID` into config/api.env `AZURE_CLIENT_ID`
- Go to `Certificates & secrets`
- Select `New client secret`
- Copy the Value into `AZURE_CLIENT_SECRET`

### Prepare the Lighthouse offering

Following steps get us the offering template:

- Start with the template [`lighthouse_template.json`](./lighthouse.tmpl.json)
- Replace `{{.TenantID}}` by value in `AZURE_TENANT_ID`
- Set a Principal ID - there are two options
  1. Use the Enterprise App as a principal.
     1. Go to Azure AD -> Enterprise applications
     2. Select the App that was automatically registered for the App registration (provisioning-service)
     3. Copy Object ID
     4. Replace `{{.PrincipalID}}` by the value copied from above
     5. Replace `{{.PrincipalName}}` by the name of the enterprise application (provisioning-service)
  2. Create a Security group and add the Enterprise App as its member
     1. Go to Azure AD -> Groups
     2. Create a group and name it
     3. Go to members and add our Enterprise App
     4. Go to the Group overview and copy its Object ID
     5. Replace `{{.PrincipalID}}` by the value copied from above
     6. Replace `{{.PrincipalName}}` by the name of the Group
- Set default offering name and description
  - Tenant account can change this while accepting the offering
  - Replace `{{.OfferingDefaultName}}` by a default offering name
  - Replace `{{.OfferingDefaultDescription}}` by a default offering description
- Put this json on publicly available URL (gist for example)

### Get lighthouse offering URL

- Obfuscate the URL by escaping it for URL possibly by https://www.urlencoder.org/
- Compose the url by replacing above url for `<obfuscatedURI>` in https://portal.azure.com/#create/Microsoft.Template/uri/<obfuscatedURI>

## Setup Tenant (Customer) account

This is the account we want to deploy into.

- Log into the account you want to deploy instances into
- Open the URI prepared in the steps above and click `Review + create`
- Confirm with `Create`
- Wait for the deployment to succeed

You're all set! :)
