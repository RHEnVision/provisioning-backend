# provisioning-backend

Provisioning backend service for cloud.redhat.com.

Requirements: Go 1.16+

## Components

* pbapi - API backend service
* pbmigrate - database migration tool with embedded SQL scripts

## Building

```
make build
```

## Configuration

Configuration is done via configuration files in `configs/` directory, see [defaults.yaml](configs/default.yaml) file for list of options with documentation. Environmental variables override the values from configuration files and always take precedence.
To use environment variables for nested configurations just join the nested keys by `_` e.g. `database.host` becomes `DATABASE_HOST`.
When `local.yaml` file is found in `configs` directory, it overrides `defaults.yaml`
When running in Clowder env, defaults are overriden by config consumed through [consoleDot shared library](github.com/RedHatInsights/app-common-go/)

## Development setup

To run all the components from this repository, you will need:

* Go compiler
* PostgreSQL server with UUID module
* GNU Makefile

```
dnf install postgresql-server postgresql-contrib
make run
```

## Code lint

We run `go fmt`, `go vet` and `golangci-lint` lint suite via GitHub Actions. To run them locally do:

```
make install-tools
make lint
```

## Migrations

Migrations can be found in `internal/db/migrations` in SQL format, the only supported database platform is PostgreSQL. To create a new migration:

```
make generate-migration MIGRATION_NAME=add_new_column
```

To execute migration either use [migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) command line utility from the golang-migrate library which supports more advanced operations like down, or run `pbmigrate` to perform up migration the same way as the application does when it boots up.

## Building container

```
podman build -t pb .
podman run --name pb1 --rm -ti -p 8000:8000 -p 5000:5000 pb
curl http://localhost:8000
curl http://localhost:5000/metrics
```

## Contributing

Here are few points before you start contributing:

* Binaries go into the `cmd/name` package.
* All code go into the `internal/` package.
* Binaries must not take any arguments.
* All configuration is done via environment variables.
* Database database models (structs) do belong into `internal/models`.
* Do not put any code logic into database models.
* All database operations (CRUD, eager loading) lives in `internal/dao` (Data Access Objects).
* Database models must be not exposed directly into JSON API, use `internal/payloads` package to wrap them.
* Business logic (the actual code) does belong into `internal/services` package, each API call should have a dedicated file.
* HTTP routes go into `internal/routes` package.
* HTTP middleware go into `internal/middleware` package.
* Monitoring metrics are in `internal/metrics' package.`
* Use the standard library context package for context operations. Context keys must be defined in `internal/ctxval`.
* Database connection is at `internal/db`, cloud connections are in 'internal/clouds'.
* Do not introduce `utils` or `tools` common packages, they tend to grow.
* Keep the line of sight (happy code path).

## License

GNU GPL 3.0, see LICENSE
