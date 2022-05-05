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

## License

GNU GPL 3.0, see LICENSE
