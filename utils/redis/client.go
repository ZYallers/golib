package redis

import (
	"fmt"
	"time"

	"github.com/ZYallers/golib/types"
	"github.com/go-redis/redis"
)

const (
	retryMaxTimes  = 3
	retrySleepTime = time.Second
)

type Redis struct {
	Client func() *redis.Client
}

func (r *Redis) NewRedis(collector *types.RedisCollector, client *types.RedisClient, optsFunc func() *redis.Options) (*redis.Client, error) {
	var newErr error
	for i := 0; i < retryMaxTimes; i++ {
		collector.Once(func() {
			if optsFunc == nil {
				optsFunc = func() *redis.Options { return &redis.Options{} }
			}
			opts := optsFunc()
			opts.Addr = client.Host + ":" + client.Port
			opts.Password = client.Pwd
			opts.DB = client.Db
			collector.Pointer = redis.NewClient(opts)
		})

		if collector.Pointer == nil {
			newErr = fmt.Errorf("new redis(%s:%s) is nil", client.Host, client.Port)
		} else {
			newErr = collector.Pointer.Ping().Err()
		}

		if newErr != nil {
			collector.Reset(func() { time.Sleep(retrySleepTime) })
		} else {
			break
		}
	}
	return collector.Pointer, newErr
}
