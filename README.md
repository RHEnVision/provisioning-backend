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
* Gorm database models (structs) do belong into `internal/models`.
* Do not put any code logic into Gorm structs (models).
* All database operations (CRUD, eager loading) lives in `internal/dao` (Data Access Objects).
* Gorm models must be not exposed directly into JSON API, use `internal/payloads` package to wrap them.
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
