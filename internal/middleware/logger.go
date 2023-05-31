package middleware

import (
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

var panicStatus = http.StatusInternalServerError

func LoggerMiddleware(rootLogger *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			bytesIn, _ := strconv.Atoi(r.Header.Get("Content-Length"))
			traceId := logging.TraceId(r.Context())
			edgeId := logging.EdgeRequestId(r.Context())
			corrId := logging.CorrelationId(r.Context())
			lctx := rootLogger.With().
				Timestamp().
				Str("trace_id", traceId).
				Str("remote_ip", r.RemoteAddr).
				Str("url", r.URL.Path).
				Str("method", r.Method).
				Int("bytes_in", bytesIn)
			if edgeId != "" {
				lctx = lctx.Str("request_id", edgeId)
			}
			if corrId != "" {
				lctx = lctx.Str("correlation_id", corrId)
			}
			logger := lctx.Logger()
			t1 := time.Now()

			lHeaders := logger.With().Fields(r.Header).Logger()
			lHeaders.Debug().Msgf("Started %s request %s", r.Method, r.URL.Path)

			// store logger under zerolog context key
			ctx := logger.WithContext(r.Context())

			defer func() {
				duration := time.Since(t1)
				log := zerolog.Ctx(ctx).With().
					Dur("latency_ms", duration).
					Int("bytes_out", ww.BytesWritten()).
					Logger()

				// prevent the application from exiting
				if rec := recover(); rec != nil {
					log.Error().
						Bool("panic", true).
						Int("status", panicStatus).
						Msgf("Unhandled panic: %s\n%s", rec, debug.Stack())
					http.Error(ww, http.StatusText(panicStatus), panicStatus)
				}

				log.Info().
					Int("status", ww.Status()).
					Msgf("Completed %s request %s in %s with %d",
						r.Method, r.URL.Path, duration.Round(time.Millisecond).String(), ww.Status())
			}()

			next.ServeHTTP(ww, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
