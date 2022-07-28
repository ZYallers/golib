package redis

import (
	"github.com/ZYallers/golib/utils/json"
	"math/rand"
	"time"
)

// 数据不存在情况下，为防止缓存雪崩，随机返回一个30到60秒的有效时间
func (r *Redis) NoDataExpiration() time.Duration {
	// 将时间戳设置成种子数
	rand.Seed(time.Now().UnixNano())
	return time.Duration(30+rand.Intn(30)) * time.Second
}

// 从String类型的缓存中读取数据，如没则重新调用指定方法重新从数据库中读取并写入缓存
func (r *Redis) CacheWithString(key string, output interface{}, expire time.Duration, fn func() (interface{}, bool)) error {
	if val := r.Client().Get(key).Val(); val != "" {
		return json.Unmarshal([]byte(val), &output)
	}
	var (
		isNull bool
		data   interface{}
	)
	if data, isNull = fn(); isNull {
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
