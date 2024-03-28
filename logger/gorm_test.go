package logger

import (
	"testing"
	"time"

	"gorm.io/gorm/logger"
)

func TestNewGormLogger(t *testing.T) {
	NewGormLogger("test", 3*time.Second, logger.Info, nil)
}
