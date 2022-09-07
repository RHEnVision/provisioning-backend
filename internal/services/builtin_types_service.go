package services

import (
	"net/http"
	"strings"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/clients/http/azure/types"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/render"
)

func ListAzureBuiltinInstanceTypes(w http.ResponseWriter, r *http.Request) {
	region := strings.ToLower(r.URL.Query().Get("region"))
	zone := strings.ToLower(r.URL.Query().Get("zone"))
	supported, err := ParseBool(r.URL.Query().Get("supported"))
	if err != nil {
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), err))
		return
	}

	if region == "" {
		renderError(w, r, payloads.NewMissingRequestParameterError(r.Context(), "region"))
		return
	}
	if zone == "" {
		renderError(w, r, payloads.NewMissingRequestParameterError(r.Context(), "zone"))
		return
	}

	start := time.Now()
	instances, err := types.InstanceTypesForZone(region, zone, supported)
	logger := ctxval.Logger(r.Context())
	logger.Trace().TimeDiff("duration", time.Now(), start).Msg("Listed instance types")
	if err != nil {
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), err))
		return
	}

	if err := render.RenderList(w, r, payloads.NewListInstanceTypeResponse(instances)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "list instance types", err))
		return
	}
}
