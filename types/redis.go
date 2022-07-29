package types

import "github.com/go-redis/redis"

type RedisClient struct {
	Db              int
	Host, Port, Pwd string
}

type RedisCollector struct {
	Done    uint32
	Pointer *redis.Client
}
