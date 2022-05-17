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

Configuration is done via configuration files in `configs/` directory, see [defaults.yaml](configs/defaults.yaml) file for list of options with documentation. Environmental variables override the values from configuration files and always take precedence.
To use environment variables for nested configurations just join the nested keys by `_` e.g. `database.host` becomes `DATABASE_HOST`.
When `local.yaml` file is found in `configs` directory, it overrides `defaults.yaml`
When running in Clowder env, defaults are overriden by config consumed through [consoleDot shared library](https://github.com/RedHatInsights/app-common-go/)

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

We currently do not allow down migrations, so delete the `down.sql` file and do not commit it into git (it will fail build).

To apply migrations, build and run `pbmigrate` binary. This is exactly what application performs during startup. The `pbmigrate` utility also supports seeding initial data. There are files available in `internal/db/seeds` with various seed configurations. Feel free to create your own configuration which you may or may not want to commit into git. When you set `DB_SEED_SCRIPT` configuration variable, the migration tool will execute all statements from that file. By default, the variable is empty, meaning no data will be seeded.

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
* All database operations (CRUD, eager loading) lives in `internal/dao` (Data Access Objects) with a common API interface.
* The actual implementation lives in `internal/dao/sqlx` package, all operations are passed with Context and errors are wrapped.
* Database models must be not exposed directly into JSON API, use `internal/payloads` package to wrap them.
* Business logic (the actual code) does belong into `internal/services` package, each API call should have a dedicated file.
* HTTP routes go into `internal/routes` package.
* HTTP middleware go into `internal/middleware` package.
* Monitoring metrics are in `internal/metrics` package.
* Use the standard library context package for context operations. Context keys are defined in `internal/ctxval` as well as accessor functions.
* Database connection is at `internal/db`, cloud connections are in `internal/clouds`.
* Do not introduce `utils` or `tools` common packages.
* Keep the line of sight (happy code path).
* PostgreSQL version we currently use in production is v14+ so take advantage of all modern SQL features.
* Use `BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY` for primary keys. This allows us to use 1-100 in our seed data (see above).
* Use `TEXT` for string columns, do not apply any limits in the DB: https://wiki.postgresql.org/wiki/Don%27t_Do_This#Text_storage

Keep security in mind: https://github.com/RedHatInsights/secure-coding-checklist

## License

GNU GPL 3.0, see LICENSE
