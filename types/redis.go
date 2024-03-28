package types

import (
	"sync"
	"sync/atomic"

	"github.com/go-redis/redis"
)

type RedisClient struct {
	Host string
	Port string
	Pwd  string
	Db   int
}

type RedisCollector struct {
	done    uint32
	m       sync.Mutex
	Pointer *redis.Client
}

func (r *RedisCollector) Once(f func()) {
	if atomic.LoadUint32(&r.done) == 0 {
		r.doSlow(f)
	}
}

func (r *RedisCollector) Reset(f func()) {
	r.m.Lock()
	defer r.m.Unlock()
	if r.done == 1 {
		defer atomic.StoreUint32(&r.done, 0)
		f()
	}
}

func (r *RedisCollector) doSlow(f func()) {
	r.m.Lock()
	defer r.m.Unlock()
	if r.done == 0 {
		defer atomic.StoreUint32(&r.done, 1)
		f()
	}
}
