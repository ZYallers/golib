package ants

import (
	"log"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewPoolAndConfig(t *testing.T) {
	NewPool(10,
		WithSubmitTimeout(8*time.Second),
		WithSubmitInterval(2*time.Second),
		WithExpiryDuration(5*time.Second),
		WithPanicHandler(func(r interface{}) {
			log.Println("panic handler:", r)
		}),
	)

	p := Pool()
	t.Logf("pool-> %#v\n", p)
	cfg := Config()
	t.Logf("config-> %#v\n", cfg)
	t.Log("pool cap:", p.Cap())
	p.Tune(200)
	t.Logf("config-> %#v\n", cfg)
	Go(func() {
		panic("test panic")
	})
	t.Log("pool cap:", p.Cap())
	time.Sleep(3 * time.Second)
}

var sum int32

func task(i int) func() {
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
	NewPool(10, WithSubmitTimeout(8*time.Second))
	for i := 0; i < 100; i++ {
		Go(task(i))
	}
	time.Sleep(30 * time.Second)
}
