package logger

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	ormLogger "gorm.io/gorm/logger"
)

const (
	regular      = "/(service|model|logic)/"
	infoStr      = "INFO: \n%s\n[info] "
	warnStr      = "WARN: \n%s\n[warn] "
	errStr       = "ERR: \n%s\n[error] "
	traceStr     = "INFO: \n%s\n[time:%.3fms] [rows:%v]\n%s"
	traceWarnStr = "WARN: \n%s\n%s\n[time:%.3fms] [rows:%v]\n%s"
	traceErrStr  = "ERR: \n%s\n%s\n[time:%.3fms] [rows:%v]\n%s"
)

type GormLogSender interface {
	Push(string)
	Open() bool
	Always() bool
}

type gormLogger struct {
	infoStr, warnStr, errStr, traceStr, traceErrStr, traceWarnStr string
	Writer                                                        *zap.Logger
	ormLogger.Config
	GormLogSender
}

func NewGormLogger(name string, slowThreshold time.Duration, level ormLogger.LogLevel, sender GormLogSender) ormLogger.Interface {
	return &gormLogger{
		GormLogSender: sender,
		Writer:        Use(name),
		Config:        ormLogger.Config{SlowThreshold: slowThreshold, LogLevel: level, Colorful: false},
		infoStr:       infoStr,
		warnStr:       warnStr,
		errStr:        errStr,
		traceStr:      traceStr,
		traceWarnStr:  traceWarnStr,
		traceErrStr:   traceErrStr,
	}
}

func (l *gormLogger) LogMode(level ormLogger.LogLevel) ormLogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l gormLogger) Printf(level ormLogger.LogLevel, format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	switch level {
	case ormLogger.Error:
		if l.Writer != nil {
			l.Writer.Error(s)
		}
		if l.Open() {
			l.Push(s)
		}
	case ormLogger.Warn:
		if l.Writer != nil {
			l.Writer.Warn(s)
		}
		if l.Open() {
			l.Push(s)
		}
	default:
		if l.Writer != nil {
			l.Writer.Info(s)
		}
		if l.Open() && l.Always() {
			l.Push(s)
		}
	}
}

func (l gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= ormLogger.Info {
		l.Printf(ormLogger.Info, l.infoStr+msg, append([]interface{}{fileWithLineNum()}, data...)...)
	}
}

func (l gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= ormLogger.Warn {
		l.Printf(ormLogger.Warn, l.warnStr+msg, append([]interface{}{fileWithLineNum()}, data...)...)
	}
}

func (l gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= ormLogger.Error {
		l.Printf(ormLogger.Error, l.errStr+msg, append([]interface{}{fileWithLineNum()}, data...)...)
	}
}

func (l gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > 0 {
		elapsed := time.Since(begin)
		switch {
		case err != nil && l.LogLevel >= ormLogger.Error:
			sql, rows := fc()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				l.Printf(ormLogger.Info, l.traceStr, fileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
			} else {
				l.Printf(ormLogger.Error, l.traceErrStr, fileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= ormLogger.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
			l.Printf(ormLogger.Warn, l.traceWarnStr, slowLog, fileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		case l.LogLevel >= ormLogger.Info:
			sql, rows := fc()
			l.Printf(ormLogger.Info, l.traceStr, fileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

func fileWithLineNum() string {
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && strings.HasSuffix(file, ".go") && regexp.MustCompile(regular).MatchString(file) {
			return fmt.Sprintf("%s:%d", file, line)
		}
	}
	return ""
}
