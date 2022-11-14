package services

import (
	"errors"
	"net/http"
	"strings"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/chi/v5"
)

func StatusService(w http.ResponseWriter, r *http.Request) {
	writeOk(w, r)
}

func ReadyService(w http.ResponseWriter, r *http.Request) {
	writeOk(w, r)
}

var UnknownReadinessServiceErr = errors.New("unknown service for readiness test")

func ReadyBackendService(w http.ResponseWriter, r *http.Request) {
	service := chi.URLParam(r, "SRV")
	switch strings.ToLower(service) {
	case "sources":
		client, err := clients.GetSourcesClient(r.Context())
		if err != nil {
			renderError(w, r, payloads.NewClientError(r.Context(), err))
			return
		}
		err = client.Ready(r.Context())
		if err != nil {
			writeServiceUnavailable(w, r)
			return
		}
	case "ib", "image_builder", "imagebuilder":
		client, err := clients.GetImageBuilderClient(r.Context())
		if err != nil {
			renderError(w, r, payloads.NewClientError(r.Context(), err))
			return
		}
		err = client.Ready(r.Context())
		if err != nil {
			writeServiceUnavailable(w, r)
			return
		}
	default:
		renderError(w, r, payloads.NewURLParsingError(r.Context(), "unable to parse SRV parameter", UnknownReadinessServiceErr))
		return
	}

	writeOk(w, r)
}
