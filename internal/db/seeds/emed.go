package seeds

import "embed"

//go:embed *.sql
var EmbeddedSeeds embed.FS
