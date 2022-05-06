package services

import (
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"net/http"

	"github.com/rs/zerolog"
)

func ContextLogger(r *http.Request) zerolog.Logger {
	return r.Context().Value(ctxval.LoggerCtxKey).(zerolog.Logger)
}
