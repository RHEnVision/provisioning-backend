package cache

import "github.com/hashicorp/golang-lru"

const lruSize = 1

const appTypeKey = "app_type"

var cache *lru.Cache

func init() {
	var err error
	if cache, err = lru.New(lruSize); err != nil {
		panic(err)
	}
}

func AppTypeId() (string, bool) {
	appTypeId, ok := cache.Get(appTypeKey)
	if ok {
		return appTypeId.(string), ok
	}
	return "", ok
}

func SetAppTypeId(appTypeId string) {
	cache.Add(appTypeKey, appTypeId)
}
