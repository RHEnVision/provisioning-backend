// Package random provides a common seed function and custom random generation
// functions.

//nolint:gosec
package random

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"

	"go.opentelemetry.io/otel/trace"
)

// SeedGlobal can be used to initialize the global thread-safe pseudo-random
// generator from the standard library. This should be called from all init()
// functions for all binaries.
func SeedGlobal() {
	var rngSeed int64
	_ = binary.Read(crand.Reader, binary.LittleEndian, &rngSeed)
	rand.Seed(rngSeed)
}

// TraceID generates a random OpenTelemetry Trace ID.
func TraceID() trace.TraceID {
	tid := trace.TraceID{}
	_, _ = rand.Read(tid[:])
	return tid
}
