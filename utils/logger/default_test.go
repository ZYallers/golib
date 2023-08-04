package logger

import (
	"fmt"
	"go.uber.org/zap"
	"testing"
)

func TestDefault(t *testing.T) {
	SetLoggerDir(".")

	//logger := zap.NewExample() // Create a new logger
	logger := Use("golib")
	SetDefault(logger) // Set the default logger
	logger.Info("info message")

	dfLogger := Default()
	dfLogger.Info("default message")

	loggerPointerAddr, defaultPointerAddr := fmt.Sprintf("%p", logger), fmt.Sprintf("%p", dfLogger)
	t.Log(loggerPointerAddr, defaultPointerAddr)

	// Check that the default logger is the same as the one we created
	if loggerPointerAddr != defaultPointerAddr {
		t.Error("the default value is different from the one created")
	}
}

func TestDebug(t *testing.T) {
	SetLoggerDir(".")
	SetDefault(Use("golib"))
	Debug("debug message", zap.String("field1", "value1"))
}

func TestInfo(t *testing.T) {
	SetLoggerDir(".")
	SetDefault(Use("golib"))
	Info("info message", zap.String("field1", "value1"))
}

func TestWarn(t *testing.T) {
	SetLoggerDir(".")
	SetDefault(Use("golib"))
	Warn("warn message", zap.String("field1", "value1"))
}

func TestError(t *testing.T) {
	SetLoggerDir(".")
	SetDefault(Use("golib"))
	Error("warn message", zap.String("field1", "value1"))
}

func TestDPanic(t *testing.T) {
	SetLoggerDir(".")
	SetDefault(Use("golib"))
	DPanic("dpanic message", zap.String("field1", "value1"))
}

func TestPanic(t *testing.T) {
	SetLoggerDir(".")
	SetDefault(Use("golib"))
	Panic("panic message", zap.String("field1", "value1"))
}

func TestFatal(t *testing.T) {
	SetLoggerDir(".")
	SetDefault(Use("golib"))
	Fatal("fatal message", zap.String("field1", "value1"))
}

func TestWith(t *testing.T) {
	SetLoggerDir(".")
	SetDefault(Use("golib"))
	With(zap.Bool("bool", true)).Debug("debug message with")
}

func TestWithOptions(t *testing.T) {
	SetLoggerDir(".")
	SetDefault(Use("golib"))
	WithOptions(zap.Fields(zap.String("option", "value"))).Debug("debug message with options")
}
