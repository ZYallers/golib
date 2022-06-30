package redis

import (
	"errors"
	"fmt"
	"github.com/ZYallers/golib/types"
	redis2 "github.com/go-redis/redis"
	"sync/atomic"
	"time"
)

const (
	retryMaxTimes  = 3
	retrySleepTime = 100 * time.Millisecond
)

type Redis struct {
	Client func() *redis2.Client
}

func (r *Redis) NewRedis(rdc *types.RedisCollector, client *types.RedisClient, options func() *redis2.Options) (*redis2.Client, error) {
	if client == nil {
		return nil, errors.New("redis client is nil")
	}
	var err error
	for i := 1; i <= retryMaxTimes; i++ {
		if atomic.LoadUint32(&rdc.Done) == 0 {
			atomic.StoreUint32(&rdc.Done, 1)
			opts := &redis2.Options{}
			if options != nil {
				opts = options()
			}
			opts.Addr = client.Host + ":" + client.Port
			opts.Password = client.Pwd
			opts.DB = client.Db
			rdc.Pointer = redis2.NewClient(opts)
		}
		if rdc.Pointer == nil {
			err = fmt.Errorf("new redis(%s:%s) is nil", client.Host, client.Port)
		} else {
			err = rdc.Pointer.Ping().Err()
		}
		if err != nil {
			atomic.StoreUint32(&rdc.Done, 0)
			if i < retryMaxTimes {
				time.Sleep(retrySleepTime)
				continue
			} else {
				return nil, fmt.Errorf("new redis(%s:%s) error: %v", client.Host, client.Port, err)
			}
		}
		break
	}
	return rdc.Pointer, nil
}
