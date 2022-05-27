package services

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"

	"github.com/rs/zerolog"
)

func ContextLogger(r *http.Request) zerolog.Logger {
	return r.Context().Value(ctxval.LoggerCtxKey).(zerolog.Logger)
}
