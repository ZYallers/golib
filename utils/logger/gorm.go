package logger

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	gormLogger "gorm.io/gorm/logger"
	"regexp"
	"runtime"
	"strings"
	"time"
)

const (
	regular      = "/(service|model|logic)/"
	infoStr      = "%s\n[info] "
	warnStr      = "%s\n[warn] "
	errStr       = "%s\n[error] "
	traceStr     = "%s\n[runtime:%.3fms] [rows:%v]\n%s"
	traceWarnStr = "%s\n%s\n[runtime:%.3fms] [rows:%v]\n%s"
	traceErrStr  = "%s\n%s\n[runtime:%.3fms] [rows:%v]\n%s"
)

type GormLogSender interface {
	Push(string)
	Open() bool
	Always() bool
}

type logger struct {
	GormLogSender
	gormLogger.Config
	Writer                                                        *zap.Logger
	infoStr, warnStr, errStr, traceStr, traceErrStr, traceWarnStr string
}

func NewGormLogger(name string, slowThreshold time.Duration, level gormLogger.LogLevel, sender GormLogSender) gormLogger.Interface {
	return &logger{
		GormLogSender: sender,
		Writer:        Use(name),
		Config:        gormLogger.Config{SlowThreshold: slowThreshold, LogLevel: level, Colorful: false},
		infoStr:       infoStr,
		warnStr:       warnStr,
		errStr:        errStr,
		traceStr:      traceStr,
		traceWarnStr:  traceWarnStr,
		traceErrStr:   traceErrStr,
	}
}

func (l *logger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l logger) Printf(level gormLogger.LogLevel, format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	switch level {
	case gormLogger.Error:
		if l.Writer != nil {
			l.Writer.Error(s)
		}
		if l.Open() {
			l.Push(s)
		}
	case gormLogger.Warn:
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

func (l logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Info {
		l.Printf(gormLogger.Info, l.infoStr+msg, append([]interface{}{fileWithLine()}, data...)...)
	}
}

func (l logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Warn {
		l.Printf(gormLogger.Warn, l.warnStr+msg, append([]interface{}{fileWithLine()}, data...)...)
	}
}

func (l logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Error {
		l.Printf(gormLogger.Error, l.errStr+msg, append([]interface{}{fileWithLine()}, data...)...)
	}
}

func (l logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > 0 {
		elapsed := time.Since(begin)
		switch {
		case err != nil && l.LogLevel >= gormLogger.Error:
			sql, rows := fc()
			if rows == -1 {
				l.Printf(gormLogger.Error, l.traceErrStr, fileWithLine(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				l.Printf(gormLogger.Error, l.traceErrStr, fileWithLine(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormLogger.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
			if rows == -1 {
				l.Printf(gormLogger.Warn, l.traceWarnStr, slowLog, fileWithLine(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				l.Printf(gormLogger.Warn, l.traceWarnStr, slowLog, fileWithLine(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case l.LogLevel >= gormLogger.Info:
			sql, rows := fc()
			if rows == -1 {
				l.Printf(gormLogger.Info, l.traceStr, fileWithLine(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				l.Printf(gormLogger.Info, l.traceStr, fileWithLine(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		}
	}
}

func fileWithLine() string {
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && strings.HasSuffix(file, ".go") && regexp.MustCompile(regular).MatchString(file) {
			return file + ":" + cast.ToString(line)
		}
	}
	return ""
}
