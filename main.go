// Copyright Red Hat
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/rhenvision/provisioning-backend/config"
)

// Handler function that responds with Hello World
func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world")
}

// statusOK returns a simple 200 status code
func statusOK(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// Serve OpenAPI spec json
func serveOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	cfg := config.Get()
	http.ServeFile(w, r, cfg.OpenAPIFilePath)
}

func initDependencies() {
	config.Init()
}

func main() {
	initDependencies()
	cfg := config.Get()
}
