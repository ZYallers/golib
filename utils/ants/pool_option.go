package ants

import "time"

type PoolConfig struct {
	PoolSize       int
	ExpiryDuration time.Duration
	SubmitTimeout  time.Duration
	SubmitInterval time.Duration
	Logger         PoolLogger
	PanicHandler   func(interface{})
}

type PoolOption func(c *PoolConfig)

func WithPoolSize(size int) PoolOption {
	return func(c *PoolConfig) {
		c.PoolSize = size
	}
}

func WithLogger(logger PoolLogger) PoolOption {
	return func(c *PoolConfig) {
		c.Logger = logger
	}
}

func WithExpiryDuration(t time.Duration) PoolOption {
	return func(c *PoolConfig) {
		c.ExpiryDuration = t
	}
}

func WithSubmitTimeout(t time.Duration) PoolOption {
	return func(c *PoolConfig) {
		c.SubmitTimeout = t
	}
}

func WithSubmitInterval(t time.Duration) PoolOption {
	return func(c *PoolConfig) {
		c.SubmitInterval = t
	}
}

func WithPanicHandler(h func(interface{})) PoolOption {
	return func(c *PoolConfig) {
		c.PanicHandler = h
	}
}
