# Configuring Amazon AWS

The provisioning service uses Amazon's AWS IAM API to perform all provisioning tasks in order to achieve the best possible security standard and practices. In this article, we will describe how to create IAM User, Role and Policy and configure the application with it.

## Configuring provisioning role

Two accounts are needed for the setup:

* **Service account** which is configured in the application (key/secret), this is an account of the service provider (e.g. Red Hat for cloud.redhat.com). It must not be root AWS account, the recommended practice is to create a dedicated user only with the minimum set of permissions.
* **Tenant account** is the account into which the service account switches via AssumeRole API operation. It needs an AWS IAM Role with set of permissions to perform operations like list instance types, create new instance, start and stop. The role is identified by AWS ARN string which is fetched from the Sources microservice, Authorization endpoint.

You will need to know account numbers of both accounts. For development and testing purposes, both can be actually a single account.

### Service account policy

Create a new policy for the Service account:

* Navigate to Identity and Access Management (IAM) on AWS.
* Click on Policies - Create policy.
* Click on JSON tab and paste the code below.
* Click on your login name in the top-right corner and copy the Account ID to the clipboard. Then edit the "Resource" string and replace `123456789` with your account ID.
* Click Next, Next.
* Give the policy a name: `redhat-provisioning-service`
* Click on Create Policy
* Select **Access key - Programmatic access** for access type.

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "RedHatProvisioningFederation",
            "Effect": "Allow",
            "Action": [
                "sts:AssumeRole",
                "sts:GetFederationToken"
            ],
            "Resource": "arn:aws:iam::123456789:role/*"
        }
    ]
}
```

### Service account user

Create a new user for the Service account:

* Navigate to Identity and Access Management (IAM) on AWS.
* Click on Users - Add users.
* Enter name: `redhat-provisioning-user`.
* Select **Access key - Programmatic access** for access type.
* On the next screen, select Attach existing policies directly and find policy named **redhat-provisioning-service**.
* Click on Next and Review.
* Confirm by clicking on Create user.
* On the next page, make sure to copy Access key and secret key and paste them both into the application configuration (e.g. `local.yaml` or K8s configuration).

### Tenant account policy

Create a new policy for the Tenant account:

* Navigate to Identity and Access Management (IAM) on AWS.
* Click on Policies - Create policy.
* Click on JSON tab and paste the code below.
* Click Next, Next.
* Give the policy a name: `redhat-provisioning-policy-1`
* Click on Create Policy
* Select **Access key - Programmatic access** for access type.

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

### Tenant account role

* Navigate to Identity and Access Management (IAM) on AWS.
* Click on Roles - Create role.
* Select **AWS account** as the entity type.
* Select trusted AWS account, click Another AWS account and enter account number of the Service account. When using only single account (for development/testing), select This account instead.
* Click Next.
* Find 'redhat-provisioning-policy-1' policy and select it, click Next.
* Enter role name: `redhat-provisioning-role-1`.
* Click on Create role.
* Copy the ARN string of the Role and store it on a safe place: `arn:aws:iam::123456789:role/redhat-provisioning-role-1`

There can be multiple tenant accounts defined, therefore it is good to give them numbers.

## Configuring Sources microservice

### Configuring locally (development setup)

Configure the ARN string in the script/sources.local.conf and run the script/sources.setup.sh and then the script/sources.seed.sh to create authorization entry, it will have database ID 1. Use this Source ID when performing provisioning operations through the app.

### Configuring on prod/stage

Go to Sources App, create new Provisioning Source, select Manually and enter the AWS ARN string. _This needs to be elaborated more, we do not have the feature fully implemented._

## Configuring AWS CloudWatch

The service can be configured to send all its logs via AWS CloudWatch API. This feature must be enabled and AWS credentials (region, key, secret) must be configured. It is not good practice to use root accounts directly, therefore in this section we will describe how to configure a dedicated user for that.

* Navigate to Identity and Access Management (IAM) on AWS.
* Click on Users - Add User.
* Enter a username: `redhat-cloudwatch-user`
* Select **Access key - Programmatic access** for access type.
* On the next screen, select Attach existing policies directly and find policy named **AmazonAPIGatewayPushToCloudWatchLogs**.
* Click on Next and Review.
* Confirm by clicking on Create user.
* On the next page, make sure to copy Access key and secret key and paste them both into the application configuration (e.g. `local.yaml` or K8s configuration).


