package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	defaultLogger, stdLogger *zap.Logger
)

type (
	Field  = zap.Field
	Level  = zapcore.Level
	Logger = zap.Logger
	Option = zap.Option
)

func Default() *Logger {
	return defaultLogger
}

func SetDefault(logger *Logger) {
	defaultLogger = logger
	stdLogger = defaultLogger.WithOptions(zap.AddCallerSkip(1))
	zap.ReplaceGlobals(defaultLogger)
	zap.RedirectStdLog(defaultLogger)
}

func Debug(msg string, fields ...Field) {
	stdLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	stdLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	stdLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...Field) {
	stdLogger.Error(msg, fields...)
}

func DPanic(msg string, fields ...Field) {
	stdLogger.DPanic(msg, fields...)
}

func Panic(msg string, fields ...Field) {
	stdLogger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	stdLogger.Fatal(msg, fields...)
}

func With(fields ...Field) *Logger {
	return stdLogger.With(fields...)
}

func WithOptions(opts ...Option) *Logger {
	return stdLogger.WithOptions(opts...)
}
