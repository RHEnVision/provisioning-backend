package services

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
)

// writeBasicError returns an error code without utilizing the Chi rendering stack. It can
// be used for fatal errors which happens during rendering pipeline (e.g. JSON errors).
func writeBasicError(w http.ResponseWriter, r *http.Request, err error) {
	if logger := zerolog.Ctx(r.Context()); logger != nil {
		logger.Error().Err(err).Msg("Unable to render error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)

	wrappedMessage := ""
	if errors.Unwrap(err) != nil {
		wrappedMessage = errors.Unwrap(err).Error()
	}
	traceId := logging.TraceId(r.Context())
	writeErrorBody(w, r, err.Error(), traceId, wrappedMessage)
}

func writeErrorBody(w http.ResponseWriter, _ *http.Request, msg, traceId, err string) {
	_, _ = w.Write([]byte(fmt.Sprintf(`{"msg": "%s", "trace_id": "%s", "error": "%s"}`, msg, traceId, err)))
}

func renderError(w http.ResponseWriter, r *http.Request, renderer render.Renderer) {
	errRender := render.Render(w, r, renderer)
	if errRender != nil {
		writeBasicError(w, r, errRender)
	}
}

func renderNotFoundOrDAOError(w http.ResponseWriter, r *http.Request, err error, resource string) {
	if errors.Is(err, dao.ErrNoRows) {
		renderError(w, r, payloads.NewNotFoundError(r.Context(), resource, err))
	} else {
		renderError(w, r, payloads.NewDAOError(r.Context(), resource, err))
	}
}

func writeEmptyResponse(w http.ResponseWriter, _ *http.Request, code int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", "0")
	w.WriteHeader(code)
}

func writeOk(w http.ResponseWriter, r *http.Request) {
	writeEmptyResponse(w, r, http.StatusOK)
}

func writeNotFound(w http.ResponseWriter, r *http.Request) {
	writeEmptyResponse(w, r, http.StatusNotFound)
}

func writeBadRequest(w http.ResponseWriter, r *http.Request) {
	writeEmptyResponse(w, r, http.StatusBadRequest)
}

func writeServiceUnavailable(w http.ResponseWriter, r *http.Request) {
	writeEmptyResponse(w, r, http.StatusServiceUnavailable)
}

func writeNoContent(w http.ResponseWriter, r *http.Request) {
	writeEmptyResponse(w, r, http.StatusNoContent)
}

func writeUnauthorized(w http.ResponseWriter, r *http.Request) {
	writeEmptyResponse(w, r, http.StatusUnauthorized)
}
