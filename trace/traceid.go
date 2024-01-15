package trace

import (
	"errors"
	"sync"

	"github.com/ZYallers/golib/goid"
	"github.com/coocood/freecache"
	"github.com/google/uuid"
)

const (
	defaultCacheSize = 100 * 1024 * 1024 // 50MB, In bytes, where 1024 * 1024 represents a single Megabyte.
	IdKey            = "Trace-Id"        // trace header key
	expireSeconds    = 60                // cache expire seconds
)

var (
	cacheSize int
	singleton sync.Once
	cache     *freecache.Cache
)

func GetCache() *freecache.Cache {
	return cache
}

func SetCacheSize(size int) {
	if cache != nil {
		panic(errors.New("cache resource has been created"))
	}
	cacheSize = size
}

func ready() {
	singleton.Do(func() {
		size := defaultCacheSize
		if cacheSize > 0 {
			size = cacheSize
		}
		cache = freecache.NewCache(size)
	})
}

func GenTraceId() string {
	return uuid.NewString()
}

func GetGoIdTraceId() string {
	return GetTraceId(goid.GetString())
}

func GetTraceId(key string) string {
	ready()
	value, _ := cache.Get([]byte(key))
	return string(value)
}

func SetTraceId(key string, value string) {
	if value == "" {
		value = GenTraceId()
	}
	ready()
	_ = cache.Set([]byte(key), []byte(value), expireSeconds)
}

func DelTraceId(key string) bool {
	ready()
	return cache.Del([]byte(key))
}

func CountTrace() int64 {
	ready()
	return cache.EntryCount()
}

func ClearTrace() {
	ready()
	cache.Clear()
}
