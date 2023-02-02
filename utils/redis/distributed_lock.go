package redis

import (
	"errors"
	"fmt"
	"github.com/ZYallers/golib/funcs/maths"
	"github.com/go-redis/redis"
	"io"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

// DistributedLock redis分布式锁
type DistributedLock struct {
	Key               string        // redis锁key
	Expiration        time.Duration // redis锁key有效时间
	NodeLockTimeout   time.Duration // 获取节点锁超时时间
	RedisLockTimeout  time.Duration // 获取redis锁超时时间
	NodeLockInterval  time.Duration // 循环获取节点锁间隔时间
	RedisLockInterval time.Duration // 循环获取redis锁间隔时间
	WatchdogRate      time.Duration // 看门狗检测频率，注意：此时间必须小于RedisLockTimeout，否则看门狗不会执行

	redisClient    *redis.Client    // redis实例
	ioWriter       io.Writer        // 日志输出实例
	watchdogQuitCh chan interface{} // 看门狗退出通道
	redisLockKey   string           // redis锁key的完整字符
	idfa           string           // 客户端标示
	debug          bool             // 是否开启调试模式
}

const (
	redisLockKeyFormat = "lock@distributed:%s:string"
)

var (
	nodeLockInt            int64 // 节点锁数字
	ErrRedisLockKeyEmpty   = errors.New("redis lock key cannot be an empty string")
	ErrRedisClientNil      = errors.New("redis client cannot be nil")
	ErrGetNodeLockTimeout  = errors.New("timeout get node lock")
	ErrGetRedisLockTimeout = errors.New("timeout get redis lock")
)

func NewDistributedLock(opts ...DistributedLockOption) (*DistributedLock, error) {
	hostname, _ := os.Hostname()
	dl := &DistributedLock{
		ioWriter:          os.Stdout,
		Expiration:        10 * time.Second,
		NodeLockTimeout:   10 * time.Second,
		RedisLockTimeout:  10 * time.Second,
		NodeLockInterval:  1 * time.Millisecond,
		RedisLockInterval: 1 * time.Millisecond,
		WatchdogRate:      2 * time.Second,
		idfa:              hostname + "." + randIntString(),
	}
	for _, opt := range opts {
		opt(dl)
	}
	if dl.redisLockKey == "" {
		return nil, ErrRedisLockKeyEmpty
	}
	if dl.redisClient == nil {
		return nil, ErrRedisClientNil
	}
	return dl, nil
}

func (dl *DistributedLock) Exec(fn func(*DistributedLock) (interface{}, error)) (interface{}, error) {
	nodeLockTimeoutCh := time.After(dl.NodeLockTimeout)
	for {
		select {
		case <-nodeLockTimeoutCh:
			dl.println("get node lock timeout")
			return nil, ErrGetNodeLockTimeout
		default:
			if !dl.nodeLockLock() {
				time.Sleep(dl.NodeLockInterval)
				//dl.println("waiting get node lock...")
				continue
			} else { // 获取到节点锁
				dl.println("got node lock")
				redisLockTimeoutCh := time.After(dl.RedisLockTimeout)
				for {
					select {
					case <-redisLockTimeoutCh:
						dl.nodeLockRelease()
						dl.println("get redis lock timeout")
						return nil, ErrGetRedisLockTimeout
					default:
						if randStr, ok := dl.redisLockLock(); !ok { // 获取redis锁失败
							time.Sleep(dl.RedisLockInterval)
							//dl.println("waiting get redis lock...")
							continue
						} else { // 获取到redis锁
							dl.println("got redis lock")
							dl.watchdogQuitCh = make(chan interface{}, 0)
							go dl.watchdog()
							dl.println("exec func starting")
							res, err := fn(dl)
							dl.println("exec func finished")
							dl.execAfter(randStr)
							return res, err
						}
					}
				}
			}
		}
	}
}

func (dl *DistributedLock) execAfter(value string) {
	if dl.redisClient.Get(dl.redisLockKey).Val() == value {
		dl.redisLockRelease()
		dl.watchdogClose()
		dl.nodeLockRelease()
	}
}

func (dl *DistributedLock) RedisClient() *redis.Client {
	return dl.redisClient
}

func (dl *DistributedLock) RedisLockKey() string {
	return dl.redisLockKey
}

func (dl *DistributedLock) nodeLockLock() bool {
	return atomic.CompareAndSwapInt64(&nodeLockInt, 0, 1)
}

func (dl *DistributedLock) nodeLockRelease() {
	atomic.StoreInt64(&nodeLockInt, 0)
	dl.println("release node lock")
}

func (dl *DistributedLock) redisLockLock() (string, bool) {
	s := randIntString()
	return s, dl.redisClient.SetNX(dl.redisLockKey, s, dl.Expiration).Val()
}

func (dl *DistributedLock) redisLockRelease() {
	dl.redisClient.Del(dl.redisLockKey)
	dl.println("release redis lock")
}

func (dl *DistributedLock) watchdog() {
	defer func() { recover() }()
	redisLockTimeout := time.After(dl.RedisLockTimeout)
	for {
		// 每隔 WatchdogRate 检查一下当前客户端是否持有redis锁，如果依然持有，那么就延长锁的过期时间
		time.Sleep(dl.WatchdogRate)
		select {
		case <-redisLockTimeout:
			dl.println("timeout watch dog")
			return
		case <-dl.watchdogQuitCh:
			dl.println("close watch dog")
			return
		default:
			if d := dl.redisClient.TTL(dl.redisLockKey).Val(); d > 0 {
				d2 := d + dl.WatchdogRate
				dl.redisClient.Expire(dl.redisLockKey, d2)
				dl.println(fmt.Sprintf("watch client holds redis lock, extend expiration %s to %s", d, d2))
			} else {
				return
			}
		}
	}
}

func (dl *DistributedLock) watchdogClose() {
	defer func() { recover() }()
	close(dl.watchdogQuitCh)
}

func (dl *DistributedLock) println(s string) {
	if dl.debug {
		ts := time.Now().Format("2006/01/02 15:04:05.00000")
		_, _ = fmt.Fprintf(dl.ioWriter, "[%s] <%s> %s\n", ts, dl.idfa, s)
	}
}

type DistributedLockOption func(dl *DistributedLock)

func WithKey(key string) DistributedLockOption {
	return func(dl *DistributedLock) {
		dl.Key = key
		dl.redisLockKey = fmt.Sprintf(redisLockKeyFormat, dl.Key)
	}
}

func WithExpiration(t time.Duration) DistributedLockOption {
	return func(dl *DistributedLock) {
		dl.Expiration = t
	}
}

func WithNodeLockTimeout(t time.Duration) DistributedLockOption {
	return func(dl *DistributedLock) {
		dl.NodeLockTimeout = t
	}
}

func WithRedisLockTimeout(t time.Duration) DistributedLockOption {
	return func(dl *DistributedLock) {
		dl.RedisLockTimeout = t
	}
}

func WithWatchdogRate(t time.Duration) DistributedLockOption {
	return func(dl *DistributedLock) {
		dl.WatchdogRate = t
	}
}

func WithNodeLockInterval(t time.Duration) DistributedLockOption {
	return func(dl *DistributedLock) {
		dl.NodeLockInterval = t
	}
}

func WithRedisLockInterval(t time.Duration) DistributedLockOption {
	return func(dl *DistributedLock) {
		dl.RedisLockInterval = t
	}
}

func WithIoWriter(w io.Writer) DistributedLockOption {
	return func(dl *DistributedLock) {
		dl.ioWriter = w
	}
}

func WithRedisClient(cli *redis.Client) DistributedLockOption {
	return func(dl *DistributedLock) {
		dl.redisClient = cli
	}
}

func WithIdfa(name string) DistributedLockOption {
	return func(dl *DistributedLock) {
		hostname, _ := os.Hostname()
		dl.idfa = hostname + "." + name
	}
}

func WithDebug() DistributedLockOption {
	return func(dl *DistributedLock) {
		dl.debug = true
	}
}

func randIntString() string {
	return strconv.Itoa(maths.RandIntn(99999999))
}
