package mysql

import (
	"log"
	"os"
	"time"

	blogger "github.com/ZYallers/golib/utils/logger"
	"gorm.io/gorm/logger"
)

const (
	DefaultLogLevel      = logger.Warn
	DefaultSlowThreshold = 500 * time.Millisecond
)

func (m *Model) DebugLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{SlowThreshold: 200 * time.Millisecond, LogLevel: logger.Info, Colorful: true},
	)
}

func (m *Model) NewLogger(filename string, args ...interface{}) logger.Interface {
	var sender blogger.GormLogSender
	logLevel, slowThreshold := DefaultLogLevel, DefaultSlowThreshold
	for _, arg := range args {
		switch v := arg.(type) {
		case logger.LogLevel:
			logLevel = v
		case time.Duration:
			slowThreshold = v
		case blogger.GormLogSender:
			sender = v
		}
	}
	return blogger.NewGormLogger(filename, slowThreshold, logLevel, sender)
}
