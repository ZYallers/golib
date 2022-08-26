package logger

import (
	"go.uber.org/zap"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func Test_Zap_Use(t *testing.T) {
	SetLoggerDir("/Users/cloud/projects/ZYallers/golib")
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
