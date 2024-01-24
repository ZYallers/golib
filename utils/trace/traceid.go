package trace

import (
	"github.com/google/uuid"
)

const IdKey = "Trace-Id"

var cache Cache

func init() {
	cache = NewFastCache(32 * 1024 * 1024)
}

func SetCache(c Cache) {
	cache.Clear()
	cache = c
}

func GetCache() Cache {
	return cache
}

func NewTraceId() string {
	return uuid.NewString()
}

func HasTraceId(key string) bool {
	return cache.Exist([]byte(key))
}

func GetTraceId(key string) string {
	return string(cache.Get([]byte(key)))
}

func SetTraceId(key string, value string) {
	cache.Set([]byte(key), []byte(value))
}

func DelTraceId(key string) {
	cache.Del([]byte(key))
}
