package logger

import (
	"go.uber.org/zap"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func Test_Zap_Use(t *testing.T) {
	SetLoggerDir("./test_log")
	Use("ddd").Info("1234")
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			is := strconv.Itoa(i)
			Use("").Error(strings.Repeat(is, 5))
			Use("eeee").Info(is, zap.Any("len", loggerDict.Len()))
		}(i)
	}
	wg.Wait()
}

func Test_Zap_NewLogger(t *testing.T) {
	fileName := func(fn, dir string) string {
		if fn == "" {
			fn = time.Now().Format("20060102")
		}
		if dir == "" {
			dir, _ = filepath.Abs(filepath.Dir("."))
		}
		fp, _ := filepath.Abs(dir + "/" + fn + suffix)
		return fp
	}
	var wg sync.WaitGroup
	for i := 0; i < 52; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			is := strconv.Itoa(i)
			NewLogger(fileName("test@"+is, "./test_log")).Info(is)
			NewLogger(fileName("record", "")).Debug(is, zap.Int("len", loggerDict.Len()))
		}(i)
	}
	wg.Wait()
}
