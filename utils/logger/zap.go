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
	suffix              = "log"
	cacheMaxSize        = 100
	perTimeDeleteNumber = 10
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

func Use(f string, opts ...Option) *zap.Logger {
	filename, dir := f, loggerDir
	if filename == "" {
		filename = time.Now().Format("20060102")
	}
	if dir == "" {
		dir, _ = filepath.Abs(filepath.Dir("."))
	}
	fp, _ := filepath.Abs(dir + "/" + filename + "." + suffix)
	return NewLogger(fp, opts...)
}

func NewLogger(f string, opts ...Option) *zap.Logger {
	v, exist := loggerDict.GetOrPutFunc(f, func(filename string) (interface{}, error) {
		cfg := &defaultConfig
		for _, opt := range opts {
			opt(cfg)
		}
		//fmt.Printf("%+v\n", cfg)
		lk := &lumberjack.Logger{Filename: filename, MaxSize: cfg.maxSize, MaxAge: cfg.maxAge,
			MaxBackups: cfg.maxBackups, LocalTime: cfg.localTime, Compress: cfg.compress}
		logger := zap.New(zapcore.NewCore(jsonEncoder, zapcore.AddSync(lk), levelEnabler))
		return logger, nil
	})
	// Randomly delete ${perTimeDeleteNumber} cached logger after exceeding ${cacheMaxSize}
	if !exist && loggerDict.Len() >= cacheMaxSize {
		counter, delNum := 0, perTimeDeleteNumber
		for key, _ := range loggerDict.Data() {
			if key == f {
				continue
			}
			if counter++; counter > delNum {
				break
			}
			loggerDict.Delete(key)
		}
	}
	logger, _ := v.(*zap.Logger)
	return logger
}
