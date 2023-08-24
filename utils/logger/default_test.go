package logger

import (
	"go.uber.org/zap"
	"log"
	"testing"
	"time"
)

func TestSetDefault(t *testing.T) {
	SetLoggerDir(".")
	logger := Use("test")
	SetDefault(logger)
	Debug("message")
	time.Sleep(3 * time.Second)
	SetDefault(Use("test2"))
	Debug("message2")
}

func TestRedirectStdLog(t *testing.T) {
	logger := Use("test2")
	RedirectStdLog(logger)
	log.Println("RedirectStdLog")
}

func TestDefault(t *testing.T) {
	SetLoggerDir(".")
	logger := Use("test")
	SetDefault(logger)
	Default().Debug("message")
}

func TestDebug(t *testing.T) {
	SetLoggerDir(".")
	SetDefault(Use("test"))
	Debug("debug message", zap.String("field1", "value1"))
}

func TestInfo(t *testing.T) {
	SetLoggerDir(".")
	SetDefault(Use("test"))
	Info("info message", zap.String("field1", "value1"))
}

func TestWarn(t *testing.T) {
	SetLoggerDir(".")
	SetDefault(Use("test"))
	Warn("warn message", zap.String("field1", "value1"))
}

func TestError(t *testing.T) {
	SetLoggerDir(".")
	SetDefault(Use("test"))
	Error("warn message", zap.String("field1", "value1"))
}

func TestDPanic(t *testing.T) {
	SetLoggerDir(".")
	SetDefault(Use("test"))
	DPanic("dpanic message", zap.String("field1", "value1"))
}

func TestPanic(t *testing.T) {
	SetLoggerDir(".")
	SetDefault(Use("test"))
	Panic("panic message", zap.String("field1", "value1"))
}

func TestFatal(t *testing.T) {
	SetLoggerDir(".")
	SetDefault(Use("test"))
	Fatal("fatal message", zap.String("field1", "value1"))
}

func TestWith(t *testing.T) {
	SetLoggerDir(".")
	SetDefault(Use("test"))
	With(zap.Bool("bool", true)).Debug("debug message with")
}

func TestWithOptions(t *testing.T) {
	SetLoggerDir(".")
	SetDefault(Use("test"))
	WithOptions(zap.Fields(zap.String("option", "value"))).Debug("debug message with options")
}
