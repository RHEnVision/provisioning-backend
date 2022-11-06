# Setting up dev environment

This document elaborates the main [README](../README.md) into more details.

## Go

We do recommend to use the [official Go build](https://go.dev/doc/install) in the version that is the minimum required version as specified in the README. Avoid using distributions from Linux package management or MacOS Homebrew because they may carry additional patches and the update pace is different.

The ideal installation location is `~/go`, then just add `~/go/bin` to PATH and restart all terminals. If you choose a different location, make sure to also set GOROOT [environment variable](https://pkg.go.dev/cmd/go#hdr-Environment_variables). There are also other variables available to change specific directories, this can be useful for moving cache or temp directories outside of your home directory and out of your backups.

Tip: Go installation can be fully done in GoLang preferences, the default installation location is `~/go`.

A great tool is [.envrc](https://direnv.net) which allows automatic switching of Go versions (or any environmental variables or PATH):

```
$ cat .envrc
GOVER=1.18
export GOROOT="$(go$GOVER env GOROOT)"
PATH_add "$(go$GOVER env GOROOT)/bin"
export GOBIN="$(pwd)/bin"
PATH_add "$(pwd)/bin"
```

## Editor or IDE

Luckily, Go language contains the very cool `go fmt` utility which we take advantage. We stick with the Go language recommendations, on top of that we use `goimports` tool for import sorting to prevent git conflicts.

We do have [EditorConfig](../.editorconfig) file that should automatically configure most of the editors and IDEs in regard to tabs vs spaces and other small details.

Most of our team members use either Jetbrains GoLand, or VSCode. Here are the recommended settings for **GoLand**:

* Make sure to select the correct Go version in Preferences - Go - GOROOT.
* You need to set `goimports` import sort style in Preferences - Editor - Code Style - Go, otherwise our CI job (and your `make fmt` target) will complain about it all the time.
* GoLand comes with batteries included, no other changes or plugins are needed.

**VSCode** recommended settings:

* Install the official Go language plugin and HTTP Client plugin.
* For the HTTP plugin to work with the provided [HTTP files](../scripts/rest_examples), open up the user settings JSON file and enter the snippet from below (copy variables from `http-client.env.json`).

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

## Compilation and startup

Use `make` command to compile the main application, use `make run` or start it manually via `./pbapi`.

The application performs automatic migration of database tables and keeps the schema up-to-date. In addition, it maintains initial data ([seed data](../internal/db/seeds/dev_small.sql)) in the database. If you delete such data, it will attempt to create it again.

Notable records created via seed script:

* Account number 13 with organization id 000013. This account is the first account (ID=1) and it is very often used on many examples (including this document). For [example](../scripts/rest_examples/http-client.env.json), RH-Identity-Header is an HTTP header that MUST be present in ALL requests, it is a base64-encoded JSON string which includes account number.
* An example SSH public key.

## Backend services

The application integrates with multiple backend services which are required to be available for most HTTP request to complete successfully.

## Sources

[Sources](https://github.com/RedHatInsights/sources-api-go) is an authentication inventory. Since it only requires Go, Redis and Postgres, we created a shell script that automatically checks out sources from git, compiles it, installs and creates postgres database, seeds data and starts the Sources application.

Follow [instructions](../scripts/README.sources) to perform the setup. Note that configuration via `sources.local.conf` is **required** before the setup procedure. This has been written and tested for Fedora Linux, in other operating systems perform all the commands manually.

Tip: On MacOS, you can install Sources on a remote Fedora Linux (or a small VM) and configure the application to connect there, instead of localhost.

Tip: Alternatively, the application supports connecting to the stage environment through a HTTP proxy. See [configuration example](../config/api.env.example) for more details. Make sure to use account number from stage environment instead of the pre-seeded account number 000013.

## Image Builder

Because Image Builder is more complex for installation, we do not recommend installing it on your local machine right now. Configure connection through HTTP proxy to the stage environment in `config/api.env`. See [configuration example](../config/api.env.example) for an example, you will need to ask someone from the company for real URLs for the service and the proxy.

## Writing Go code

Ready to write some Go code? Read [contributing guide](../CONTRIBUTING.md).
