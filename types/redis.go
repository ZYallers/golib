package types

import "github.com/go-redis/redis"

type RedisClient struct {
	Host, Port, Pwd string
	Db              int
}

type RedisCollector struct {
	Done    uint32
	Pointer *redis.Client
}
