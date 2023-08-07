package logger

import (
	"go.uber.org/zap"
)

var defaultLogger, stdLogger *zap.Logger

func SetDefault(logger *zap.Logger) {
	cp := *logger
	defaultLogger = &cp
	stdLogger = defaultLogger.WithOptions(zap.AddCallerSkip(1))
	zap.ReplaceGlobals(defaultLogger)
	zap.RedirectStdLog(defaultLogger)
}

func Default() *zap.Logger                       { return defaultLogger }
func With(fields ...zap.Field) *zap.Logger       { return stdLogger.With(fields...) }
func WithOptions(opts ...zap.Option) *zap.Logger { return stdLogger.WithOptions(opts...) }
func Debug(msg string, fields ...zap.Field)      { stdLogger.Debug(msg, fields...) }
func Info(msg string, fields ...zap.Field)       { stdLogger.Info(msg, fields...) }
func Warn(msg string, fields ...zap.Field)       { stdLogger.Warn(msg, fields...) }
func Error(msg string, fields ...zap.Field)      { stdLogger.Error(msg, fields...) }
func DPanic(msg string, fields ...zap.Field)     { stdLogger.DPanic(msg, fields...) }
func Panic(msg string, fields ...zap.Field)      { stdLogger.Panic(msg, fields...) }
func Fatal(msg string, fields ...zap.Field)      { stdLogger.Fatal(msg, fields...) }
