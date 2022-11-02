# Development setup scripts

This directory contains scripts which helps with initial setup for development setup which is container-less. Meaning, it will install software in your OS (PostgreSQL, Kafka) and checkout services from git (sources).

These scripts are tested on Fedora Server or Workstation, latest stable release. You can use them as a step by step tutorial in case you are on a different OS.

## Kafka service

Kafka is needed for inter-service streaming and also for sources availability checks.

### Configuration

Review ./kafka.conf configuration. If you would like to do any changes, create ./kafka.local.conf which overrides the main configuration file and it is ignored by git. Currently, no changes are needed - you can even keep the unique cluster ID since for development purposes, only single-process instance is deployed.

### Set up

Run `./kafka.setup.sh` to download, extract and configure Kafka as a single-process deployment.

### Start up

Run `./kafka.start.sh` to start Kafka.

### Clean up

Run `./kafka.clean.sh` to start Kafka. This will also delete all data stored in `/tmp`!

## Sources service

Sources are needed to fetch authorization information for cloud providers.

### Configuration

Review ./sources.conf configuration. If you would like to do any changes, create ./sources.local.conf which overrides the main configuration file and it is ignored by git.

The minimum configuration values required are:

* ARN_ROLE: Amazon AWS role ARN string that will be used to seed the sources. A new Source, Application and Authentication records will be created with the credential. The Account/Tenant will be created with account_id/org_id 13/000013, therefore the same account number must be used in provisioning application. The resulting Source record will have ID 1 and name "Amazon provisioning". Example format: arn:aws:iam::123456789:role/redhat-provisioning-1
* SUBSCRIPTION_ID: Azure subscription ID that will be used to seed the sources.
* PROJECT_ID: GCP project ID that will be used to seed the sources.

If you run the seed script for the first time, the authorization records will have database (primary key) IDs of 1, 2 and 3 for AWS EC2, Amazon and GCP respectively.

### Set up

Install database, redis, init database, create sources user and database. Can be executed multiple times. You will be asked for "sudo" password.

    ./sources.setup.sh

Checkout, compile and start sources backend. If you want to update the app,
just do "git pull" and start again.

### Start up

    ./sources.start.sh

### Seed

Populate database with some data, with the sources app running do:

    ./sources.seed.sh

### Clean up

When you want to start over, to delete the database and user and git checkouts:

    ./sources.clean.sh
