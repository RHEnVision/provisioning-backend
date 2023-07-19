// Package random provides a common seed function and custom random generation
// functions.

//nolint:gosec
package random

import (
	crand "crypto/rand"
	"encoding/binary"
	mrand "math/rand"

	"go.opentelemetry.io/otel/trace"
	erand "golang.org/x/exp/rand"
)

// SeedGlobal can be used to initialize the global thread-safe pseudo-random
// generator from the standard library. This should be called from all init()
// functions for all binaries.
//
// It seeds both math/rand and e/exp/rand packages just in case IDE imports
// one or the other by accident.
func SeedGlobal() {
	var seedInt64 int64
	_ = binary.Read(crand.Reader, binary.LittleEndian, &seedInt64)
	mrand.Seed(seedInt64)

	var seedUint64 uint64
	_ = binary.Read(crand.Reader, binary.LittleEndian, &seedUint64)
	erand.Seed(seedUint64)
}

// TraceID generates a random OpenTelemetry Trace ID.
func TraceID() trace.TraceID {
	tid := trace.TraceID{}
	_, _ = mrand.Read(tid[:])
	return tid
}

// Float32 returns mathematical random float number in the (0.0, 1.0> interval.
func Float32() float32 {
	return mrand.Float32()
}
