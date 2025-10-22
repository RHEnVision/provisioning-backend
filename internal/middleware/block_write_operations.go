package middleware

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/config"
)

func BlockWriteOperationsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isProd := config.InProdClowder()
		isStage := config.InStageClowder()

		if (isProd || isStage) && (r.Method == "POST" || r.Method == "PUT") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotImplemented)
			_, _ = w.Write([]byte(`{"msg": "Write operations are not available in this environment, service is decommissioned"}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}
