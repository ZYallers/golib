package ants

import (
	"log"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewPoolAndConfig(t *testing.T) {
	NewPool(
		WithPoolSize(10),
		WithSubmitTimeout(8*time.Second),
		WithSubmitInterval(2*time.Second),
		WithExpiryDuration(5*time.Second),
		WithPanicHandler(func(r interface{}) {
			log.Println("panic handler:", r)
		}),
	)
	t.Logf("pool-> %#v\n", Pool())
	t.Logf("config-> %#v\n", Config())
	Go(func() {
		panic("test panic")
	})
	time.Sleep(3 * time.Second)
}

var sum int32

func myTask(i int) func() {
	return func() {
		if i == 99 {
			panic("test panic")
		}
		atomic.AddInt32(&sum, 1)
		time.Sleep(1 * time.Second)
		log.Printf("No.%d-------%d\n", i, atomic.LoadInt32(&sum))
	}
}

func TestGo(t *testing.T) {
	NewPool(
		WithPoolSize(10),
		WithSubmitTimeout(8*time.Second),
	)
	for i := 0; i < 100; i++ {
		Go(myTask(i))
	}
	time.Sleep(30 * time.Second)
}
