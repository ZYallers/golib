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
	gormless "gorm.io/gorm/logger"
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
	gormless.Config
	infoStr      string
	warnStr      string
	errStr       string
	traceStr     string
	traceErrStr  string
	traceWarnStr string
	Writer       *zap.Logger
	Sender       GormLogSender
}

func NewGormLogger(name string, slowThreshold time.Duration, level gormless.LogLevel, sender GormLogSender) gormless.Interface {
	return &gormLogger{
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
		Sender:       sender,
		Writer:       Use(name),
		Config:       gormless.Config{SlowThreshold: slowThreshold, LogLevel: level, Colorful: false},
	}
}

func (l *gormLogger) LogMode(level gormless.LogLevel) gormless.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l gormLogger) Printf(level gormless.LogLevel, format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	switch level {
	case gormless.Error:
		if l.Writer != nil {
			l.Writer.Error(s)
		}
		if l.Sender != nil && l.Sender.Open() {
			l.Sender.Push(s)
		}
	case gormless.Warn:
		if l.Writer != nil {
			l.Writer.Warn(s)
		}
		if l.Sender != nil && l.Sender.Open() {
			l.Sender.Push(s)
		}
	default:
		if l.Writer != nil {
			l.Writer.Info(s)
		}
		if l.Sender != nil && l.Sender.Open() && l.Sender.Always() {
			l.Sender.Push(s)
		}
	}
}

func (l gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormless.Info {
		l.Printf(gormless.Info, l.infoStr+msg, append([]interface{}{fileWithLineNum()}, data...)...)
	}
}

func (l gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormless.Warn {
		l.Printf(gormless.Warn, l.warnStr+msg, append([]interface{}{fileWithLineNum()}, data...)...)
	}
}

func (l gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormless.Error {
		l.Printf(gormless.Error, l.errStr+msg, append([]interface{}{fileWithLineNum()}, data...)...)
	}
}

func (l gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > 0 {
		elapsed := time.Since(begin)
		switch {
		case err != nil && l.LogLevel >= gormless.Error:
			sql, rows := fc()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				l.Printf(gormless.Info, l.traceStr, fileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
			} else {
				l.Printf(gormless.Error, l.traceErrStr, fileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormless.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
			l.Printf(gormless.Warn, l.traceWarnStr, slowLog, fileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		case l.LogLevel >= gormless.Info:
			sql, rows := fc()
			l.Printf(gormless.Info, l.traceStr, fileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
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
