package services

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	clientshttp "github.com/RHEnVision/provisioning-backend/internal/clients/http"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/render"
)

// writeBasicError returns an error code without utilizing the Chi rendering stack. It can
// be used for fatal errors which happens during rendering pipeline (e.g. JSON errors).
func writeBasicError(w http.ResponseWriter, r *http.Request, err error) {
	if logger := ctxval.Logger(r.Context()); logger != nil {
		logger.Error().Msgf("unable to render error %v", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)

	wrappedMessage := ""
	if errors.Unwrap(err) != nil {
		wrappedMessage = errors.Unwrap(err).Error()
	}
	traceId := ctxval.TraceId(r.Context())
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

func renderNewErrorFromClientErr(w http.ResponseWriter, r *http.Request, err error) {
	var sourceError clientshttp.SourceError
	var ibError clientshttp.ImageBuilderError
	if errors.As(err, &sourceError) {
		renderError(w, r, payloads.SourcesError(r.Context(), err))
	} else if errors.As(err, &ibError) {
		renderError(w, r, payloads.NewImageBuilderError(r.Context(), err))
	} else if errors.Is(err, clients.UnknownAuthenticationTypeErr) {
		renderError(w, r, payloads.NewUnknownAuthenticationType(r.Context(), err))
	} else if errors.Is(err, clients.ClientErr) {
		renderError(w, r, payloads.ClientError(r.Context(), err))
	} else {
		renderError(w, r, payloads.GeneralError(r.Context(), "client side error", err))
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
