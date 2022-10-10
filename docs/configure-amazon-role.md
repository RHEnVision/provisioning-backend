# Configuring Amazon AWS

The provisioning service uses Amazon's AWS IAM API to perform all provisioning tasks.
We try to follow the best possible security standard and practices.
In this article, we will describe how to create necessary IAM setup and configure this application with it.

## Configuring provisioning role

There are two accounts used in the provisioning flow:

* **Service account** which is configured in the application (key/secret), this is an account of the service provider (e.g. Red Hat for cloud.redhat.com). 
    It must not be root AWS account, the recommended practice is to create a dedicated user only with the minimum set of permissions.
* **Tenant account** is the account into which the service account switches via AssumeRole API operation.
    It needs an AWS IAM Role with set of permissions to perform operations this application is supposed to perform.
    These are list instance types, create new instance, start and stop. The role is identified by AWS ARN string.

You will need to know account number of the Service account.
For development and testing purposes, you can use the same account for both.

### Service account

This is the account that this application needs the credentials and will use to connect to AWS services.

#### Service account policy

Create new Policy to be used by the service account:

1. Navigate to Identity and Access Management (IAM) on AWS.
2. Click on Policies - Create policy.
3. Click on JSON tab and paste the code below.
4. Click Next, Next.
5. Give the policy a name, e.g.: `redhat-provisioning-policy`
6. Click on Create Policy

The policy JSON to use in step above:
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "AssumeRole",
            "Effect": "Allow",
            "Action": [
                "sts:AssumeRole",
                "sts:GetFederationToken"
            ],
            "Resource": "arn:aws:iam::*:role/*"
        },
        {
            "Sid": "DescribeGenericData",
            "Effect": "Allow",
            "Action": [
                "ec2:DescribeInstanceTypeOfferings",
                "ec2:DescribeRegions",
                "ec2:DescribeInstanceTypes"
            ],
            "Resource": "*"
        }
    ]
}
```

#### Service account user

Create a new user for the Service account:

* Navigate to Identity and Access Management (IAM) on AWS.
* Click on Users - Add users.
* Enter name: `redhat-provisioning-user`.
* Select **Access key - Programmatic access** for access type.
* On the next screen, select Attach existing policies directly and find policy named **redhat-provisioning-policy**.
* Click on Next and Review.
* Confirm by clicking on Create user.
* On the next page, make sure to copy Access key and secret key and paste them both into the application configuration (e.g. `config/api.env` or K8s configuration).

#### Service account regions and STS endpoints

The application needs to be able to connect to all regions, by default only selected regions are enabled:

* Navigate to Account on the top-right side under account name.
* Scroll down to AWS Regions
* Enable regions (Bahrain, Jakarta, Cape Town, UAE, Milan and others)

To be able to use global STS endpoint, enable STS version 1 tokens for all the regions:

* Navigate to Identity and Access Management (IAM) on AWS.
* Click on Account Settings
* In the STS section, click on Edit for the Global endpoint section.
* Set to "Valid in all AWS Regions"

For more info visit the [IAM User Guide](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_enable-regions.html#sts-regions-manage-tokens) and [Managing AWS Regions](https://docs.aws.amazon.com/general/latest/gr/rande-manage.html).

### Tenant account

This is a setup in the account in which the service shall deploy the actual instances.

#### Tenant policy

Create a new policy for the Tenant account:

* Navigate to Identity and Access Management (IAM) on AWS.
* Click on Policies - Create policy.
* Click on JSON tab and paste the code below.
* Click Next, Next.
* Give the policy a name: `redhat-provisioning-policy-1`
* Click on Create Policy

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "RedHatProvisioning",
            "Effect": "Allow",
            "Action": [
                "ec2:CreateKeyPair",
                "ec2:CreateTags",
                "ec2:DeleteKeyPair",
                "ec2:DeleteTags",
                "ec2:DescribeAvailabilityZones",
                "ec2:DescribeImages",
                "ec2:DescribeInstanceTypes",
                "ec2:DescribeInstances",
                "ec2:DescribeKeyPairs",
                "ec2:DescribeRegions",
                "ec2:DescribeSecurityGroups",
                "ec2:DescribeSnapshotAttribute",
                "ec2:DescribeTags",
                "ec2:ImportKeyPair",
                "ec2:RunInstances",
                "ec2:StartInstances"
            ],
            "Resource": "*"
        }
    ]
}
```

#### Tenant account role

* Navigate to Identity and Access Management (IAM) on AWS.
* Click on Roles - Create role.
* Select **AWS account** as the entity type.
* Select trusted AWS account, click Another AWS account and enter account number of the Service account.
  * If using only single account (for development/testing), select This account instead.
* Click Next.
* Find `redhat-provisioning-policy-1` policy and select it, click Next.
* Enter role name: `redhat-provisioning-role-1`.
* Click on Create role.
* Copy the ARN string of the Role and store it on a safe place: `arn:aws:iam::123456789:role/redhat-provisioning-role-1`

There can be multiple tenant accounts defined, therefore it is good to give them numbers.

## Configuring Sources microservice

In real life the Tenant account ARN will be stored in secret place, which is external microservice.
This microservice also provides guide on how to set up the Tenant account.
This app then fetches the ARN and assumes the role by this ARN.
If the setup is not correct as to above steps, this service fails, but does not provide corrective measures.

### Configuring locally (development setup)

Configure the ARN string in the script/sources.local.conf and run the script/sources.setup.sh and then the script/sources.seed.sh to create authorization entry, it will have database ID 1.
Use this Source ID when performing provisioning operations through the app.

For more details on Sources setup in development follow [Dev environment guide](dev-environment.md#Sources)

### Configuring on ConsoleDot prod/stage

Go to Sources App, create new Provisioning Source, select Manually and enter the AWS ARN string. _This needs to be elaborated more, we do not have the feature fully implemented._

## Configuring AWS CloudWatch

*This is bit unrelated, but It is another use of AWS services by our app, so it's together for the time being.*

The service can be configured to send all its logs via AWS CloudWatch API. This feature must be enabled and AWS credentials (region, key, secret) must be configured. It is not good practice to use root accounts directly, therefore in this section we will describe how to configure a dedicated user for that.

* Navigate to Identity and Access Management (IAM) on AWS.
* Click on Users - Add User.
* Enter a username: `redhat-cloudwatch-user`
* Select **Access key - Programmatic access** for access type.
* On the next screen, select Attach existing policies directly and find policy named **AmazonAPIGatewayPushToCloudWatchLogs**.
* Click on Next and Review.
* Confirm by clicking on Create user.
* On the next page, make sure to copy Access key and secret key and paste them both into the application configuration (e.g. `config/api.env` or K8s configuration).


