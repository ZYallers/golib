package trace

import "github.com/VictoriaMetrics/fastcache"

type fastCache struct {
	cache *fastcache.Cache
}

func NewFastCache(maxBytes int) *fastCache {
	return &fastCache{cache: fastcache.New(maxBytes)}
}

func (c *fastCache) Set(k, v []byte) {
	c.cache.Set(k, v)
}

func (c *fastCache) Get(k []byte) []byte {
	if v, exist := c.cache.HasGet(nil, k); exist {
		return v
	}
	return nil
}

func (c *fastCache) Del(k []byte) {
	c.cache.Del(k)
}

func (c *fastCache) Exist(k []byte) bool {
	return c.cache.Has(k)
}

func (c *fastCache) Clear() {
	c.cache.Reset()
}
