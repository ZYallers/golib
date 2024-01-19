package logger

import (
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestSetLoggerDir(t *testing.T) {
	SetLoggerDir(".")
}

func TestGetLoggerDir(t *testing.T) {
	GetLoggerDir()
}

func TestUse(t *testing.T) {
	SetLoggerDir(".")
	Use("test").Debug("message")
}

func TestUse2(t *testing.T) {
	SetLoggerDir("./test")
	var wg sync.WaitGroup
	for i := 0; i < 105; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			is := strconv.Itoa(i)
			Use(is).Debug(is, zap.Any("len", loggerDict.Len()))
		}(i)
	}
	wg.Wait()
}

func TestNewLogger(t *testing.T) {
	NewLogger("./test2.log").Debug("message")
}

func TestNewLogger2(t *testing.T) {
	fileName := func(fn, dir string) string {
		if fn == "" {
			fn = time.Now().Format("20060102")
		}
		if dir == "" {
			dir, _ = filepath.Abs(filepath.Dir("."))
		}
		fp, _ := filepath.Abs(dir + "/" + fn + ".log")
		return fp
	}
	var wg sync.WaitGroup
	for i := 0; i < 52; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			is := strconv.Itoa(i)
			NewLogger(fileName("test@"+is, "./test")).Info(is)
			NewLogger(fileName("record", "")).Debug(is, zap.Int("len", loggerDict.Len()))
		}(i)
	}
	wg.Wait()
}
