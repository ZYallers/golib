package redis

import (
	"errors"
	"fmt"
	"github.com/ZYallers/golib/funcs/arrays"
	strings2 "github.com/ZYallers/golib/funcs/strings"
	"github.com/ZYallers/golib/types"
	"github.com/ZYallers/golib/utils/json"
	redis2 "github.com/go-redis/redis"
	"math/rand"
	"strings"
	"sync/atomic"
	"time"
)

type Redis struct {
	Client func() *redis2.Client
}

const (
	retryMaxTimes   = 3
	hashAllFieldKey = "all"
)

func (r *Redis) NewRedis(rdc *types.RedisCollector, client *types.RedisClient) (*redis2.Client, error) {
	var err error
	for i := 1; i <= retryMaxTimes; i++ {
		if atomic.LoadUint32(&rdc.Done) == 0 {
			atomic.StoreUint32(&rdc.Done, 1)
			rdc.Pointer, err = r.newClient(client)
		}
		if err == nil {
			if rdc.Pointer == nil {
				err = fmt.Errorf("new redis(%s:%s) is nil", client.Host, client.Port)
			} else {
				err = rdc.Pointer.Ping().Err()
			}
		}
		if err != nil {
			atomic.StoreUint32(&rdc.Done, 0)
			if i < retryMaxTimes {
				time.Sleep(time.Millisecond * time.Duration(i*100))
				continue
			} else {
				return nil, fmt.Errorf("new redis(%s:%s) error: %v", client.Host, client.Port, err)
			}
		}
		break
	}
	return rdc.Pointer, nil
}

func (r *Redis) newClient(client *types.RedisClient) (*redis2.Client, error) {
	if client == nil {
		return nil, errors.New("redis client is nil")
	}
	rds := redis2.NewClient(&redis2.Options{
		Addr:     client.Host + ":" + client.Port,
		Password: client.Pwd,
		DB:       client.Db,
	})
	return rds, nil
}

// 数据不存在情况下，为防止缓存雪崩，随机返回一个30到60秒的有效时间
func (r *Redis) NoDataExpiration() time.Duration {
	// 将时间戳设置成种子数
	rand.Seed(time.Now().UnixNano())
	return time.Duration(30+rand.Intn(30)) * time.Second
}

// 从String类型的缓存中读取数据，如没则重新调用指定方法重新从数据库中读取并写入缓存
func (r *Redis) CacheWithString(key string, output interface{}, expiration time.Duration, fn func() (interface{}, bool)) error {
	if val := r.Client().Get(key).Val(); val != "" {
		return json.Unmarshal(strings2.String2Bytes(val), &output)
	}

	var (
		isNull bool
		data   interface{}
	)

	if data, isNull = fn(); isNull {
		expiration = r.NoDataExpiration()
	}

	var value string
	bte, err := json.Marshal(data)
	if err != nil {
		value = "null"
	} else {
		value = strings2.Bytes2String(bte)
		_ = json.Unmarshal(bte, &output)
	}
	return r.Client().Set(key, value, expiration).Err()
}

// 根据key删除对应缓存
func (r *Redis) DeleteCache(key ...string) (int64, error) {
	return r.Client().Del(key...).Result()
}

func (r *Redis) HGetAll(key string) (result []interface{}) {
	all := r.Client().HGet(key, hashAllFieldKey).Val()
	if all == "" {
		return
	}
	keys := arrays.RemoveDuplicateWithString(strings.Split(all, ","))
	if len(keys) == 0 {
		return
	}
	result = r.Client().HMGet(key, keys...).Val()
	return
}

func (r *Redis) HMSet(key string, data map[string]interface{}) error {
	fields := make([]string, 0)
	fieldValues := make(map[string]interface{}, 0)
	for k, v := range data {
		if k == "" || v == nil {
			continue
		}
		if b, err := json.Marshal(v); err == nil {
			fieldValues[k] = strings2.Bytes2String(b)
			fields = append(fields, k)
		}
	}

	if len(fields) == 0 {
		return errors.New("the data that can be saved is empty")
	}

	if val := r.Client().HGet(key, hashAllFieldKey).Val(); val != "" {
		fields = append(fields, strings.Split(val, ",")...)
	}

	var allFieldValue string
	if len(fields) > 0 {
		allFieldValue = strings.Join(arrays.RemoveDuplicateWithString(fields), ",")
	}
	fieldValues[hashAllFieldKey] = allFieldValue
	return r.Client().HMSet(key, fieldValues).Err()
}

func (r *Redis) HMDelete(key string, fields ...string) error {
	newFields := make([]string, 0)
	if val := r.Client().HGet(key, hashAllFieldKey).Val(); val != "" {
		newFields = append(newFields, strings.Split(val, ",")...)
	}
	if len(newFields) > 0 {
		for _, field := range fields {
			newFields = arrays.RemoveWithString(newFields, field)
		}
	}

	var allFieldValue string
	if len(newFields) > 0 {
		allFieldValue = strings.Join(newFields, ",")
	}

	pl := r.Client().Pipeline()
	pl.HDel(key, fields...)
	pl.HSet(key, hashAllFieldKey, allFieldValue)
	_, err := pl.Exec()
	return err
}
