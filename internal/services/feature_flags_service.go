package services

import (
	"net/http"
	"strings"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/go-chi/chi/v5"
)

func FeatureFlagService(w http.ResponseWriter, r *http.Request) {
	flag := chi.URLParam(r, "FLAG")

	if !strings.HasPrefix(flag, config.Unleash.Prefix) {
		writeBadRequest(w, r)
		return
	}

	if config.FeatureEnabled(r.Context(), flag) {
		writeOk(w, r)
	} else {
		writeUnauthorized(w, r)
	}
}
