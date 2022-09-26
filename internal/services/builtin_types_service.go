package services

import (
	"net/http"
	"strings"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/render"
)

type InstanceTypesForZoneFunc func(region, zone string, supported *bool) ([]*clients.InstanceType, error)

func ListBuiltinInstanceTypes(typeFunc InstanceTypesForZoneFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		start := time.Now()
		instances, err := typeFunc(region, zone, supported)
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
}
