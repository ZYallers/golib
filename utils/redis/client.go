package redis

import (
	"fmt"
	"github.com/ZYallers/golib/types"
	"github.com/go-redis/redis"
	"time"
)

const (
	retryMaxTimes  = 3
	retrySleepTime = time.Second
)

type Redis struct {
	Client func() *redis.Client
}

func (r *Redis) NewRedis(rdc *types.RedisCollector, cli *types.RedisClient, f func() *redis.Options) (*redis.Client, error) {
	var newErr error
	for i := 0; i < retryMaxTimes; i++ {
		rdc.Once(func() {
			opts := &redis.Options{}
			if f != nil {
				opts = f()
			}
			opts.Addr = cli.Host + ":" + cli.Port
			opts.Password = cli.Pwd
			opts.DB = cli.Db
			rdc.Pointer = redis.NewClient(opts)
		})

		if rdc.Pointer == nil {
			newErr = fmt.Errorf("new redis(%s:%s) is nil", cli.Host, cli.Port)
		} else {
			newErr = rdc.Pointer.Ping().Err()
		}

		if newErr != nil {
			rdc.Reset(func() { time.Sleep(retrySleepTime) })
		} else {
			break
		}
	}
	return rdc.Pointer, newErr
}
