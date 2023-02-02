# Configuring Azure deployment

The provisioning service uses Azure SDK to perform all provisioning tasks.
We follow the best possible security standard and practices.
In this article, we will describe how to create necessary service account setup and configure this application with it.

## Setup Azure service account

### Prepare the Azure account

This account is the one that our app will use to connect to the Azure service.

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

- Take [`lighthouse_template.json`](./lighthouse_template.json)
- Replace `<ProviderTenantID>` by value in `AZURE_TENANT_ID`
- Go to Azure AD -> Enterprise applications
  - Select the App that was automatically registered for the App registration (provisioning-service)
  - Copy Object ID
- Replace `<EnterpriseAppID>` by the value copied from above
- Replace `<EnterpriseAppName>` by the name of the enterprise application (provisioning-service)
- Put this json on publicly available URL (gist for example)

### Get lighthouse offering URL

- Obfuscate the URL by escaping it for URL possibly by https://www.urlencoder.org/
- Compose the url by replacing above url for `<obfuscatedURI>` in https://portal.azure.com/#create/Microsoft.Template/uri/<obfuscatedURI>

## Setup Tenant account

This is the account we want to deploy into.

- Log into the account you want to deploy instances into
- Open the URI prepared in the steps above and click `Review + create`
- Confirm with `Create`
- Wait for the deployment to succeed

You're all set! :)
