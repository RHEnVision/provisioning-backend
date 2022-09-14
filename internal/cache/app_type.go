package cache

import (
	"context"
	"sync"
)

var (
	appTypeId      string
	appTypeIdMutex sync.Mutex
)

func FindAppTypeId(_ context.Context) (string, bool) {
	appTypeIdMutex.Lock()
	defer appTypeIdMutex.Unlock()

	return appTypeId, appTypeId != ""
}

func SetAppTypeId(_ context.Context, newAppTypeId string) {
	appTypeIdMutex.Lock()
	defer appTypeIdMutex.Unlock()

	appTypeId = newAppTypeId
}
