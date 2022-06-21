package api

import (
	_ "embed"
	"fmt"
	"net/http"
)

//go:embed openapi.gen.json
var embeddedJSONSpec []byte

func ServeOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.WriteHeader(http.StatusCreated)
	_, err := w.Write(embeddedJSONSpec)
	if err != nil {
		w.WriteHeader(508)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"msg": "%s"}`, err.Error())))
	}
}
