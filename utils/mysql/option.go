package mysql

import (
	blogger "github.com/ZYallers/golib/utils/logger"
	"time"

	"gorm.io/gorm/logger"
)

type Option func(c *modelConfig)

type modelConfig struct {
	maxIdleConns    int
	maxOpenConns    int
	connMaxIdleTime time.Duration
	connMaxLifeTime time.Duration
	logger          logger.Interface
}

func (*Model) WithLogger(args ...interface{}) Option {
	return func(c *modelConfig) {
		var sender blogger.GormLogSender
		filename, logLevel, slowThreshold := "", DefaultLogLevel, DefaultSlowThreshold
		for _, arg := range args {
			switch v := arg.(type) {
			case logger.Interface:
				c.logger = v
				return
			case string:
				filename = v
			case logger.LogLevel:
				logLevel = v
			case time.Duration:
				slowThreshold = v
			case blogger.GormLogSender:
				sender = v
			}
		}
		c.logger = blogger.NewGormLogger(filename, slowThreshold, logLevel, sender)
	}
}

func (m *Model) WithDebugLogger() Option {
	return func(c *modelConfig) { c.logger = m.DebugLogger() }
}

func (*Model) WithMaxIdleConns(i int) Option {
	return func(c *modelConfig) { c.maxIdleConns = i }
}

func (*Model) WithMaxOpenConns(i int) Option {
	return func(c *modelConfig) { c.maxOpenConns = i }
}

func (*Model) WithConnMaxIdleTime(t time.Duration) Option {
	return func(c *modelConfig) { c.connMaxIdleTime = t }
}

func (*Model) WithConnMaxLifeTime(t time.Duration) Option {
	return func(c *modelConfig) { c.connMaxLifeTime = t }
}
