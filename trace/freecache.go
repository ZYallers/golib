package trace

import "github.com/coocood/freecache"

type freeCache struct {
	cache *freecache.Cache
}

func NewFreeCache(maxBytes int) *freeCache {
	return &freeCache{cache: freecache.NewCache(maxBytes)}
}

func (c *freeCache) Set(k, v []byte) {
	_ = c.cache.Set(k, v, 0)
}

func (c *freeCache) Get(k []byte) []byte {
	if v, err := c.cache.Get(k); err == freecache.ErrNotFound || err != nil {
		return nil
	} else {
		return v
	}
}

func (c *freeCache) Del(k []byte) {
	c.cache.Del(k)
}

func (c *freeCache) Exist(k []byte) bool {
	if _, err := c.cache.Get(k); err == freecache.ErrNotFound || err != nil {
		return false
	}
	return true
}

func (c *freeCache) Clear() {
	c.cache.Clear()
}
