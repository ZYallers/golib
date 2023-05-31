package logger

import (
	"github.com/ZYallers/golib/consts"
	"github.com/ZYallers/golib/funcs/safe"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"path/filepath"
	"time"
)

const (
	// The maximum size in megabytes of the log file before it gets rotated. It defaults to 100 megabytes.
	maxSize = 100

	// The maximum number of old log files to retain.
	// The default is to retain all old log files (though MaxAge may still cause them to get deleted.
	maxBackups = 20

	// The maximum number of running caches, if exceeded, will trigger a deletion mechanism to free memory.
	cacheMaxSize = 50

	// The log files suffix.
	suffix = ".log"
)

var (
	loggerDir    string
	loggerDict                        = safe.NewDict()
	levelEnabler zap.LevelEnablerFunc = func(lv zapcore.Level) bool { return lv >= zapcore.DebugLevel }
	jsonEncoder                       = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(consts.TimeFormatLogger))
		},
	})
)

func SetLoggerDir(dir string) {
	loggerDir = dir
}

func GetLoggerDir() string {
	return loggerDir
}

func Use(filename string) *zap.Logger {
	fn, dir := filename, loggerDir
	if fn == "" {
		fn = time.Now().Format("20060102")
	}
	if dir == "" {
		dir, _ = filepath.Abs(filepath.Dir("."))
	}
	fp, _ := filepath.Abs(dir + "/" + fn + suffix)
	return NewLogger(fp)
}

func NewLogger(filename string) *zap.Logger {
	v, ok := loggerDict.GetOrPutFunc(filename, func(fn string) (interface{}, error) {
		lk := &lumberjack.Logger{MaxSize: maxSize, MaxBackups: maxBackups, LocalTime: true, Compress: false, Filename: fn}
		logger := zap.New(zapcore.NewCore(jsonEncoder, zapcore.AddSync(lk), levelEnabler))
		return logger, nil
	})
	if !ok && loggerDict.Len() >= cacheMaxSize {
		// Exceeding maximum value, randomly delete half
		counter, clean := 0, cacheMaxSize/2
		for key, _ := range loggerDict.Data() {
			if key == filename {
				continue
			}
			if counter++; counter > clean {
				break
			}
			loggerDict.Delete(key)
		}
	}
	logger, _ := v.(*zap.Logger)
	return logger
}
