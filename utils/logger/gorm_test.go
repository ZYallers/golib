package logger

import (
	"gorm.io/gorm/logger"
	"testing"
	"time"
)

func TestNewGormLogger(t *testing.T) {
	NewGormLogger("test", 3*time.Second, logger.Info, &gormLogger{})
}
