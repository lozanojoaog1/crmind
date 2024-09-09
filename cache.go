package main

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	dataCache *cache.Cache
)

func initCache() {
	dataCache = cache.New(5*time.Minute, 10*time.Minute)
}

func getCachedData(key string) (interface{}, bool) {
	return dataCache.Get(key)
}

func setCachedData(key string, data interface{}, duration time.Duration) {
	dataCache.Set(key, data, duration)
}
