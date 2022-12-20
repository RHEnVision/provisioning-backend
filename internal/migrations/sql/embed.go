package sql

import "embed"

//go:embed *.sql
var EmbeddedSQLMigrations embed.FS
