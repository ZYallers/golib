package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestDefault(t *testing.T) {
	// Create a new logger
	logger := zap.NewExample()
	// Set the default logger
	SetDefault(logger)
	// Check that the default logger is the same as the one we created
	assert.Equal(t, logger, Default())
}

func TestDebug(t *testing.T) {
	// Create a new logger
	logger := zap.NewExample()

	// Set the default logger
	SetDefault(logger)

	// Call Debug with some fields
	Debug("test message", zap.String("field1", "value1"))
}

func TestInfo(t *testing.T) {
	// Create a new logger
	logger := zap.NewExample()

	// Set the default logger
	SetDefault(logger)

	// Call Info with some fields
	Info("test message", zap.String("field1", "value1"))
}

func TestWarn(t *testing.T) {
	// Create a new logger
	logger := zap.NewExample()

	// Set the default logger
	SetDefault(logger)

	// Call Warn with some fields
	Warn("test message", zap.String("field1", "value1"))
}

func TestError(t *testing.T) {
	// Create a new logger
	logger := zap.NewExample()

	// Set the default logger
	SetDefault(logger)

	// Call Error with some fields
	Error("test message", zap.String("field1", "value1"))
}

func TestDPanic(t *testing.T) {
	// Create a new logger
	logger := zap.NewExample()

	// Set the default logger
	SetDefault(logger)

	// Call DPanic with some fields
	DPanic("test message", zap.String("field1", "value1"))
}
