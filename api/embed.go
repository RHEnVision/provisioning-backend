package api

import (
	_ "embed"
	"fmt"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/middleware"
)

//go:embed openapi.gen.json
var embeddedJSONSpec []byte

var etag *middleware.ETag

func init() {
	var err error
	etag, err = middleware.GenerateETagFromBuffer("json-spec", middleware.OpenAPIExpiration, embeddedJSONSpec)
	if err != nil {
		panic(err)
	}
}

// ETagValue returns etag generated from the input file content with 30 minute expiration.
func ETagValue() *middleware.ETag {
	return etag
}

func ServeOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(embeddedJSONSpec)
	if err != nil {
		w.WriteHeader(508)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"msg": "%s"}`, err.Error())))
	}
}
