package redis

import (
	"sync"
	"sync/atomic"

	"github.com/go-redis/redis"
)

type Collector struct {
	done    uint32
	m       sync.Mutex
	Pointer *redis.Client
}

func (c *Collector) Once(f func()) {
	if atomic.LoadUint32(&c.done) == 0 {
		c.doSlow(f)
	}
}

func (c *Collector) Reset(f func()) {
	c.m.Lock()
	defer c.m.Unlock()
	if c.done == 1 {
		defer atomic.StoreUint32(&c.done, 0)
		f()
	}
}

func (c *Collector) doSlow(f func()) {
	c.m.Lock()
	defer c.m.Unlock()
	if c.done == 0 {
		defer atomic.StoreUint32(&c.done, 1)
		f()
	}
}
