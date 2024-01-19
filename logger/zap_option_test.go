package logger

import (
	"testing"
	"time"
)

func TestWithMaxSize(t *testing.T) {
	Use("test", WithMaxSize(20)).Debug("message")
}

func TestWithMaxAge(t *testing.T) {
	Use("test", WithMaxAge(int(time.Hour))).Debug("message")
}

func TestWithMaxBackups(t *testing.T) {
	Use("test", WithMaxBackups(20)).Debug("message")
}

func TestWithLocalTime(t *testing.T) {
	Use("test", WithLocalTime(false)).Debug("message")
}

func TestWithCompress(t *testing.T) {
	Use("test", WithCompress(true)).Debug("message")
}
