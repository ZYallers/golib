package redis

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ZYallers/golib/utils/json"
	"github.com/go-redis/redis"
)

// NoDataExpiration 数据不存在情况下，为防止缓存雪崩，随机返回一个30到60秒的有效时间
func (r *Redis) NoDataExpiration() time.Duration {
	rand.Seed(time.Now().UnixNano()) // 将时间戳设置成种子数
	return time.Duration(30+rand.Intn(30)) * time.Second
}

// CacheWithString 从String类型的缓存中读取数据，如没则重新调用指定方法重新从数据库中读取并写入缓存
func (r *Redis) CacheWithString(key string, output interface{}, expire time.Duration, getDataFunc func() (data interface{}, isNull bool)) error {
	if val := r.Client().Get(key).Val(); val != "" {
		return json.Unmarshal([]byte(val), &output)
	}
	var isNull bool
	var data interface{}
	if data, isNull = getDataFunc(); isNull {
		expire = r.NoDataExpiration()
	}
	var value string
	bte, err := json.Marshal(data)
	if err != nil {
		value = "null"
	} else {
		value = string(bte)
		_ = json.Unmarshal(bte, &output)
	}
	return r.Client().Set(key, value, expire).Err()
}

// DeductStock 减库存，通过lua脚本执行以防超卖
func (r *Redis) DeductStock(key string, quantity int) (bool, error) {
	var script = redis.NewScript(`
        if tonumber(redis.call("get", KEYS[1])) > 0 then
			redis.call("decrby", KEYS[1], tonumber(ARGV[1]))
			return 1
        else
			return 0
        end
    `)
	result, err := script.Run(r.Client(), []string{key}, quantity).Result()
	if err != nil && err != redis.Nil {
		return false, fmt.Errorf("script run error: %v", err)
	}
	if i, ok := result.(int64); ok && i == 0 {
		return false, ErrOutOfStock
	}
	return true, nil
}

const FrequencyLimitKey = "rds@freq:limit:%v"

// FrequencyLimit 频率限制, 返回ture代表频率过高
func (r *Redis) FrequencyLimit(key interface{}, second uint8) bool {
	b := r.Client().SetNX(fmt.Sprintf(FrequencyLimitKey, key), "1", time.Duration(second)*time.Second).Val()
	return !b
}
