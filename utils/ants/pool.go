package ants

import (
	"reflect"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/panjf2000/ants/v2"
)

const (
	defaultPoolSize       = 100
	defaultSubmitTimeout  = 30 * time.Second
	defaultSubmitInterval = 10 * time.Millisecond
	defaultExpiryDuration = 10 * time.Second
)

var (
	pool   *ants.Pool
	config *PoolConfig
)

func NewPool(options ...PoolOption) {
	config = &PoolConfig{
		PoolSize:       defaultPoolSize,
		ExpiryDuration: defaultExpiryDuration,
		SubmitTimeout:  defaultSubmitTimeout,
		SubmitInterval: defaultSubmitInterval,
		Logger:         &poolLogger{},
		PanicHandler: func(r interface{}) {
			config.Logger.Printf("worker exits from panic: %v\n%s", r, debug.Stack())
		},
	}
	for _, option := range options {
		option(config)
	}
	pool, _ = ants.NewPool(config.PoolSize,
		ants.WithNonblocking(true),
		ants.WithLogger(config.Logger),
		ants.WithExpiryDuration(config.ExpiryDuration),
		ants.WithPanicHandler(config.PanicHandler),
	)
}

func Pool() *ants.Pool { return pool }

func Config() *PoolConfig { return config }

func Go(task func()) {
	go func(f func()) {
		defer func() {
			if r := recover(); r != nil {
				config.Logger.Printf("task %s exits from panic: %v\n%s", taskName(task), r, debug.Stack())
			}
		}()
		toCh := time.After(config.SubmitTimeout)
		for {
			select {
			case <-toCh:
				config.Logger.Printf("task %s submit timeout %s", taskName(task), config.SubmitTimeout)
				return
			default:
				if err := pool.Submit(f); err != nil {
					if config.SubmitInterval > 0 {
						time.Sleep(config.SubmitInterval)
					}
				} else {
					return
				}
			}
		}
	}(task)
}

func taskName(task func()) string {
	funcValue := reflect.ValueOf(task)
	funcName := runtime.FuncForPC(funcValue.Pointer()).Name()
	return funcName
}
