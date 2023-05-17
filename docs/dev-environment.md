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

## Postgres

Installation and configuration of Postgres is not covered neither in this article nor in the README. Full administrator access into an empty database is assumed as the application performs DDL commands during start. Or create a user with create table privileges.

Tip: On MacOS, you can install Postgres on a remote Linux (or a small VM) and configure the application to connect there, instead of localhost.

## Kafka

In order to work on Kafka integrated services (statuser, sources), Kafka local deployment is needed. We do simply use the official Kafka binary that can be simply extracted and started.

The [scripts](../scripts) directory contains README with further instructions and scripts which can download, extract, configure and start Kafka for local development.

## Compilation and startup

Use `make` command to compile the main application, use `make run` or start it manually via `./pbapi`.

The application performs automatic migration of database tables and keeps the schema up-to-date. In addition, it maintains initial data ([seed data](../internal/db/seeds/dev_small.sql)) in the database. If you delete such data, it will attempt to create it again.

Migration files reside in the [migrations/sql](../internal/migrations/sql) directory, for each migration a new file prefixed with sequence integer must be present. These are "up" migrations, we do not have any "down" migrations at the moment.

It is possible to create a Go function that will be executed before a SQL migration, in that case create a function in [migrations/code](../internal/migrations/code) directory and update map in [migrations/callbacks.go](../internal/migrations/callbacks.go) with the sequence number of a SQL file. The code will be executed BEFORE the SQL in a transaction. If the function returns an error, the program panics and SQL migration does not start executing. In case only code migration is needed, create an empty SQL file with a number and use the number to create a function.

Notable records created via seed script:

* Account number 13 with organization id 000013. This account is the first account (ID=1) and it is very often used on many examples (including this document). For [example](../scripts/rest_examples/http-client.env.json), RH-Identity-Header is an HTTP header that MUST be present in ALL requests, it is a base64-encoded JSON string which includes account number.
* An example SSH public key.

## Backend services

The application integrates with multiple backend services:

## Worker

Worker processes (`pbworker`) are responsible for running background jobs. There must be one or more processes running in order to pick up background jobs (e.g. launch reservations). There are multiple configuration options available via `WORKER_QUEUE`:

* `redis` - uses queue via Redis
* `memory` - in-memory worker (default option)

The default behavior is the in-memory worker, which spawns a single goroutine within the main application which picks up all jobs sequentially. This is only meant for development setups so that no extra worker process is required when testing background jobs.

In stage/prod, we currently use `redis`.

## Statuser

Statuser process (`pbstatuser`) is a custom executable that runs in a single instance responsible for performing sources availability checks. These are requested over HTTP from the Sources app (see below), messages are enqueued in Kafka where the statuser instance picks them up in batches, performs checking, and sends the results back to Kafka to Sources.

## Sources

[Sources](https://github.com/RedHatInsights/sources-api-go) is an authentication inventory. Since it only requires Go, Redis and Postgres, we created a shell script that automatically checks out sources from git, compiles it, installs and creates postgres database, seeds data and starts the Sources application.

Follow [instructions](../scripts/README.sources) to perform the setup. Note that configuration via `sources.local.conf` is **required** before the setup procedure. This has been written and tested for Fedora Linux, in other operating systems perform all the commands manually.

Tip: On MacOS, you can install Sources on a remote Fedora Linux (or a small VM) and configure the application to connect there, instead of localhost.

Tip: Alternatively, the application supports connecting to the stage environment through a HTTP proxy. See [configuration example](../config/api.env.example) for more details. Make sure to use account number from stage environment instead of the pre-seeded account number 000013.

## Image Builder

Because Image Builder is more complex for installation, we do not recommend installing it on your local machine right now. Configure connection through HTTP proxy to the stage environment in `config/api.env`. See [configuration example](../config/api.env.example) for an example, you will need to ask someone from the company for real URLs for the service and the proxy.

## Containerized environment

A `docker-compose.yml` file helps rolling up the provisioning application, including frontend and other services such as sources locally for development purpose or demo with no extra setup.

Please notice that the compose file use a dedicated dev Dockerfile.dev both for backend and frontend, it uses [CompileDaemon](github.com/githubnemo/CompileDaemon) for live reloading, it watches for changes and re-build using `go build` command when a change occurs, no need to build the container after code changes.

### Install
A docker or podman (with [podman-compose](https://github.com/containers/podman-compose)) is needed, the folder structure should be:
```
.
├── provisioning-backend
├── provisioning-frontend
├── sources-api-go
└── image-builder-frontend
```

Edit [app.env](/config/api.env.example) to fit containerized services (i.e db, redis, kafka), these are not exposed directly to your localhost ports.

Run 
```sh
$ COMPOSE_PROFILES=migrate docker compose up 
```
This command also migrates data to postgres db, using the `migrate` profile.

Alternatively do:
```sh
$ COMPOSE_PROFILES=migrate podman-compose up 
```


### Profiles
A compose profile allows you to run a subset of containers. When no profile is given, 
the provisioning backend, postgres and redis will run by default.

Currently there are a few profiles:
- migrate: migrate provisioning backend, terminates after migration
- kafka: run kafka with zookeeper, register topics
- frontend: run local provisioning frontend
- sources: run local sources with postgres db, on first use notice that you will need to run `/script/sources.seed.sh` for seeding your local sources data.

For example, in order to run sources, kafka and frontend profiles, run
```sh
# using docker
$ COMPOSE_PROFILES=frontend,kafka,sources docker compose up 
```

## Writing Go code

Ready to write some Go code? Read [contributing guide](../CONTRIBUTING.md).
