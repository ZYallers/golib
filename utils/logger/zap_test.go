package logger

import (
	"github.com/spf13/cast"
	"go.uber.org/zap"
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
			Use("").Error(strings.Repeat(cast.ToString(i), 5))
			Use("eeee").Info(cast.ToString(i), zap.Any("len", loggerDict.Len()))
		}(i)
	}
	wg.Wait()
}
