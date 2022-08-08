package types

import (
	"github.com/go-redis/redis"
	"sync"
)

type RedisClient struct {
	Db              int
	Host, Port, Pwd string
}

type RedisCollector struct {
	Done    uint32
	M       sync.Mutex
	Pointer *redis.Client
}
