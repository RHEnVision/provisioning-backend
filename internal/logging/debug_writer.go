package logging

import (
	"strings"

	"github.com/rs/zerolog"
)

type DebugWriter struct {
	Logger *zerolog.Logger
}

func (dw *DebugWriter) Write(p []byte) (n int, err error) {
	dw.Logger.Debug().Msg(strings.TrimRight(string(p), "\n"))
	return len(p), nil
}
