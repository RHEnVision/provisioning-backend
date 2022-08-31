package api

import (
	_ "embed"
	"fmt"
	"net/http"
)

//go:embed openapi.gen.json
var embeddedJSONSpec []byte

func ServeOpenAPISpec(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(embeddedJSONSpec)
	if err != nil {
		w.WriteHeader(508)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"msg": "%s"}`, err.Error())))
	}
}
