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

## Code style

We run `go fmt`, `goimports`, `go vet` and `golangci-lint` lint suite via GitHub Actions. To run them locally do:

```
make install-tools
make fmt lint
```

Make sure to set your editor to use [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports) formatting style of code.

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

Check out [CONTRIBUTING.md](CONTRIBUTING.md)

## License

GNU GPL 3.0, see LICENSE
