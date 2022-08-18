package middleware

import (
	"fmt"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/version"
)

func VersionMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if version.BuildCommit != "" && version.BuildTime != "" {
			w.Header().Set("X-RH-Version", fmt.Sprintf("%s (%s)", version.BuildCommit, version.BuildTime))
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
