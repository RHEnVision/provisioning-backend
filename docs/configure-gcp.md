# Configuring Google Cloud Plattform GCP

The provisioning service uses Google cloud compute engine to perform all provisioning tasks.
We try to follow the best possible security standard and practices.
In this article, we will describe how to create necessary service account setup and configure this application with it.

## Configuring provisioning role

There are two accounts used in the provisioning flow:

* **Service account** which is configured using GOOGLE_APPLICATION_CREDENTIALS environment variable which holds a credentials JSON file.
    This is an account of the service provider (e.g. Red Hat for cloud.redhat.com).
    It needs a Role with a set of permissions that enable the application to perform operations. 
    These are list machine types, import pubkey, create new instance, start and stop.

* **Tenant account** is the account into which the service account switches via GCP's IAM permissions to API operation.
    It needs to specify the permissions it gives to the service account when adding it as an IAM.

You will need to know the e-mail of the Service account.
For development and testing purposes, you can use the same account for both.

### Service account

This is the account that this application needs the credentials and will use to connect to GCP services.

Create two projects from two different google accounts, one served as the service account's project and one as the tenant's project. 

#### Service account Configuration

Create a new Service account:

1. Navigate to the Service account project
2. Navigate to IAM & Admin on GCP.
3. Click Service Accounts
4. Click Create Service Account.
5. Fill in the Service Account details.
6. Click CREATE AND CONTINUE.
7. Choose the permissions the provisioning app needs (TBC):

Roles:
   - Service Account User
   - Compute Admin (We are allowing here more permissions than needed)

8. Click CONTINUE
9.  Click DONE 
10. Copy the service account's e-mail


### Tenant account

This is a setup in the account in which the service shall deploy the actual instances.

#### Tenant Configuration

1. Navigate to the Tenant project
2. Navigate to IAM & Admin on GCP.
3. Click IAM.
4. Click +ADD.
5. Fill in the Service Account e-mail under New principals.
6. Choose the permissions the provisioning app needs (TBC): 
Roles:
   - Service Account User
   - Compute Admin (We are allowing here more permissions than needed)
7. Click SAVE.

#### Authenticating as the service account

1. In the Google Cloud console, go to the Service accounts page.
2. Select the service account project 
3. On the Service accounts page, click the e-mail address of the service account that you want to create a key for.
4. Click the Keys tab.
5. Click the Add key drop-down menu, then select Create new key.
6. Select JSON as the Key type and click Create.
Clicking Create downloads a service account key file. 


The downloaded key has the following format, where PRIVATE_KEY is the private portion of the public/private key pair:
```json
{
  "type": "service_account",
  "project_id": "PROJECT_ID",
  "private_key_id": "KEY_ID",
  "private_key": "-----BEGIN PRIVATE KEY-----\nPRIVATE_KEY\n-----END PRIVATE KEY-----\n",
  "client_e-mail": "SERVICE_ACCOUNT_e-mail",
  "client_id": "CLIENT_ID",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://accounts.google.com/o/oauth2/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/SERVICE_ACCOUNT_e-mail"
}
```

7. Set the GCP_JSON environment variable to the downloaded JSON in base64:
  ```shell
   cat creds.json | base64
  ```
8. Paste the service account project id under GCP_PROJECT_ID 
9. Paste the base64 json in api.env under GCP_JSON variable (notice there are no new lines)

10. Paste the Tenant's project id under PROJECT_ID variable in sources.local.conf

11. Create an image in Image Builder for GCP https://console.stage.redhat.com/api/image-builder/v1. 
Share that image with the **service account** you have created (Copy the service account's email from the IAM console and paste it in Image builder wizard). 

## Configuring Sources microservice

In real life, the Tenant project id will be stored in a secret place, which is an external microservice.
This microservice also provides a guide on how to set up the Tenant account.
This app then fetches the project id and performs operations in this project.
If the setup is not correct as the above steps, this service fails, but does not provide corrective measures.
