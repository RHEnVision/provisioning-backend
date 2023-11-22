package sources

import (
	"context"
	"errors"
	"sync"
)

var ErrCacheMiss = errors.New("source constants requested, before populating")

type sourcesConstantsType struct {
	AppTypeId   string
	SourceTypes map[string]string
}

var (
	sourcesConstantsValues     sourcesConstantsType
	sourcesConstantsCacheMutex sync.Mutex
)

// FindAppTypeId returns "application id" special identifier, or returns ErrNotFound.
func getCachedConstants(_ context.Context) (sourcesConstantsType, error) {
	sourcesConstantsCacheMutex.Lock()
	defer sourcesConstantsCacheMutex.Unlock()

	if sourcesConstantsValues.AppTypeId == "" {
		return sourcesConstantsValues, ErrCacheMiss
	}
	return sourcesConstantsValues, nil
}

// SetAppTypeId sets "application id" special identifier.
func setCachedConstants(_ context.Context, value sourcesConstantsType) {
	sourcesConstantsCacheMutex.Lock()
	defer sourcesConstantsCacheMutex.Unlock()

	sourcesConstantsValues = value
}
