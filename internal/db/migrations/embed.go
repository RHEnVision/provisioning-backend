package migrations

import "embed"

//go:embed *.sql
var EmbeddedMigrations embed.FS
