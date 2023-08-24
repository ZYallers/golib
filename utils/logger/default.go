package logger

import (
	"go.uber.org/zap"
)

var defaultLogger *zap.Logger

func SetDefault(log *zap.Logger) {
	defaultLogger = clone(log)
}

func RedirectStdLog(log *zap.Logger) {
	stdLogger := log.WithOptions(zap.AddCallerSkip(1))
	zap.ReplaceGlobals(stdLogger)
	zap.RedirectStdLog(stdLogger)
}

func Default() *zap.Logger                       { return defaultLogger }
func With(fields ...zap.Field) *zap.Logger       { return defaultLogger.With(fields...) }
func WithOptions(opts ...zap.Option) *zap.Logger { return defaultLogger.WithOptions(opts...) }
func Debug(msg string, fields ...zap.Field)      { defaultLogger.Debug(msg, fields...) }
func Info(msg string, fields ...zap.Field)       { defaultLogger.Info(msg, fields...) }
func Warn(msg string, fields ...zap.Field)       { defaultLogger.Warn(msg, fields...) }
func Error(msg string, fields ...zap.Field)      { defaultLogger.Error(msg, fields...) }
func DPanic(msg string, fields ...zap.Field)     { defaultLogger.DPanic(msg, fields...) }
func Panic(msg string, fields ...zap.Field)      { defaultLogger.Panic(msg, fields...) }
func Fatal(msg string, fields ...zap.Field)      { defaultLogger.Fatal(msg, fields...) }

func clone(log *zap.Logger) *zap.Logger {
	cp := *log
	return &cp
}
