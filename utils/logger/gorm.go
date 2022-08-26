package logger

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"regexp"
	"runtime"
	"strings"
	"time"
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

type logger struct {
	infoStr, warnStr, errStr, traceStr, traceErrStr, traceWarnStr string
	Writer                                                        *zap.Logger
	gormLogger.Config
	GormLogSender
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
		l.Printf(gormLogger.Info, l.infoStr+msg, append([]interface{}{fileWithLineNum()}, data...)...)
	}
}

func (l logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Warn {
		l.Printf(gormLogger.Warn, l.warnStr+msg, append([]interface{}{fileWithLineNum()}, data...)...)
	}
}

func (l logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Error {
		l.Printf(gormLogger.Error, l.errStr+msg, append([]interface{}{fileWithLineNum()}, data...)...)
	}
}

func (l logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > 0 {
		elapsed := time.Since(begin)
		switch {
		case err != nil && l.LogLevel >= gormLogger.Error:
			sql, rows := fc()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				l.Printf(gormLogger.Info, l.traceStr, fileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
			} else {
				l.Printf(gormLogger.Error, l.traceErrStr, fileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormLogger.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
			l.Printf(gormLogger.Warn, l.traceWarnStr, slowLog, fileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		case l.LogLevel >= gormLogger.Info:
			sql, rows := fc()
			l.Printf(gormLogger.Info, l.traceStr, fileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
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
