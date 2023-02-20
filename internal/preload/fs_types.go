package preload

import (
	"embed"
)

//go:embed *.yaml *_availability/*.yaml
var fsTypes embed.FS
