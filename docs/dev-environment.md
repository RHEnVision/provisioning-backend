# Setting up dev environment

This document elaborates the main [README](../README.md) into more details.

## Go

The project Makefile contains a target which downloads required version of Go and required tools. The only prerequisite is any version of Go present on the system. We do recommend to use the [official Go build](https://go.dev/doc/install), just make sure the command `go` is on the `PATH`. To install the required version of go:

    make install-go

This will download `go1.YY.ZZ` into `GOROOT` directory (`~/go` when not specified) and then downloads the Go SDK and unpacks is into the subdirectory. See the official Go documentation on how this works, you can change the GOROOT [environment variable](https://pkg.go.dev/cmd/go#hdr-Environment_variables) in order to use a different directory. There are also other variables available to change specific directories, this can be useful for moving cache or temp directories outside of your home directory and out of your backups.

To install required utilities (linter, OpenAPI generator):

        make install-tools

Utilities are installed into `./bin` subdirectory of the project folder. The specific Go version and the tools from the `bin` directory are used in all targets in the `Makefile`. When using IDE, make sure to use this version of Go and these tools too, or use `make` to perform compilation or run tests.

## Editor or IDE

Luckily, Go language contains the very cool `go fmt` utility which we take advantage. We stick with the Go language recommendations, on top of that we use `goimports` tool for import sorting to prevent git conflicts.

We do have [EditorConfig](../.editorconfig) file that should automatically configure most of the editors and IDEs in regard to tabs vs spaces and other small details.

Most of our team members use either Jetbrains GoLand, or VSCode. Here are the recommended settings for **GoLand**:

* Make sure to select the correct Go version in Preferences - Go - GOROOT.
* You need to set `goimports` import sort style in Preferences - Editor - Code Style - Go, otherwise our CI job (and your `make fmt` target) will complain about it all the time.
* GoLand comes with batteries included, no other changes or plugins are needed.
* Make sure to use the correct version of Go and utilities from the `bin` project subdirectory.

**VSCode** recommended settings:

* Install the official Go language plugin and HTTP Client plugin.
* For the HTTP plugin to work with the provided [HTTP files](../scripts/rest_examples), open up the user settings JSON file and enter the snippet from below (copy variables from `http-client.env.json`).
* Make sure to use the correct version of Go and utilities from the `bin` project subdirectory.

```json
{
  "rest-client.environmentVariables": {
    "dev": {
      "hostname": "x",
      "port": "x"
    }
  }
}
```

## Code checking and linting

We use several formatting and linting tools (like go fmt, goimports, golangci-lint) to enforce code style and quality. There are multiple make targets available and an alias which runs them all: `make fmt`.

Additionally, migration files (`internal/db/migrations/*.sql`) must not be modified, you are only allowed to add new files into this directory. One of the validators will enforce this, however, it is possible to skip it by creating an empty commit on top of the change in case it's necessary (e.g. reformat the SQL code).

## Additional Go versions

Go language fully supports multi-version installation. Just follow [the official documentation](https://go.dev/doc/manage-install) by using the `go install` command.

Tip: GoLang fully supports multi-version installation and you can switch versions in preferences.

## Make

A make utility is used, we test on GNU Make which is available on all supported platforms. Use `make help` or read the [on-line make help](../docs/make.md) for additional help on what's available.

## Additional utilities

There are few utilities that you will need like code linter, `goimports` or migration tool for creating new migrations. Install them with `make install-tools`.

## Compilation and startup

Use `make` command to compile the main application, use `make run` or start it manually via `./pbapi`.

The application performs automatic migration of database tables and keeps the schema up-to-date. In addition, it maintains initial data ([seed data](../internal/db/seeds/dev_small.sql)) in the database. If you delete such data, it will attempt to create it again.

Migration files reside in the [migrations/sql](../internal/migrations/sql) directory, for each migration a new file prefixed with sequence integer must be present. These are "up" migrations, we do not have any "down" migrations at the moment.

It is possible to create a Go function that will be executed before a SQL migration, in that case create a function in [migrations/code](../internal/migrations/code) directory and update map in [migrations/callbacks.go](../internal/migrations/callbacks.go) with the sequence number of a SQL file. The code will be executed BEFORE the SQL in a transaction. If the function returns an error, the program panics and SQL migration does not start executing. In case only code migration is needed, create an empty SQL file with a number and use the number to create a function.

Notable records created via seed script:

* Account number 13 with organization id 000013. This account is the first account (ID=1) and it is very often used on many examples (including this document). For [example](../scripts/rest_examples/http-client.env.json), RH-Identity-Header is an HTTP header that MUST be present in ALL requests, it is a base64-encoded JSON string which includes account number.
* An example SSH public key.

## Worker

Worker processes (`pbworker`) are responsible for running background jobs. There must be one or more processes running in order to pick up background jobs (e.g. launch reservations). There are multiple configuration options available via `WORKER_QUEUE`:

* `redis` - uses queue via Redis
* `memory` - in-memory worker (default option)

The default behavior is the in-memory worker, which spawns a single goroutine within the main application which picks up all jobs sequentially. This is only meant for development setups so that no extra worker process is required when testing background jobs.

In stage/prod, we currently use `redis`.

## Statuser

Statuser process (`pbstatuser`) is a custom executable that runs in a single instance responsible for performing sources availability checks. These are requested over HTTP from the Sources app (see below), messages are enqueued in Kafka where the statuser instance picks them up in batches, performs checking, and sends the results back to Kafka to Sources.

## Backend services

The application integrates with multiple backend services:

* Postgres
* Kafka
* Redis
* RBAC Service
* Sources Service
* Notifications Service

All backend services can be started easily via [provisioning-compose](https://github.com/RHEnVision/provisioning-compose) on a local machine or remotely.

### Image Builder on Stage

Because Image Builder is more complex for installation, we do not recommend installing it on your local machine right now. Configure connection through HTTP proxy to the stage environment in `config/api.env`. See [configuration example](../config/api.env.example) for an example, you will need to ask someone from the company for real URLs for the service and the proxy.

### Notifications on Stage

When you just want to verify a notification kafka's messages, you can use `send-notification.http` to send a message directly to stage env, please notice that a cookie session is required, [click here](https://internal.console.stage.redhat.com/api/turnpike/session/) to generate one. 

## Writing Go code

Ready to write some Go code? Read [contributing guide](../CONTRIBUTING.md).
