package trace

import (
	"github.com/ZYallers/golib/goid"
	"testing"
	"time"
)

func TestSetCache(t *testing.T) {
	cache := NewFreeCache(10 * 1024 * 1024)
	SetCache(cache)
	SetTraceId("1", NewTraceId())
	time.Sleep(5 * time.Second)
	value := GetTraceId("1")
	t.Log(value)
}

func TestGetCache(t *testing.T) {
	t.Logf("%#v", GetCache())
}

func TestNewTraceId(t *testing.T) {
	t.Log(NewTraceId())
}

func TestSetTraceId(t *testing.T) {
	SetTraceId("1", NewTraceId())
	time.Sleep(5 * time.Second)
	t.Logf("%#v\n", GetTraceId("1"))
	t.Logf("%#v", GetTraceId("2"))
}

func TestGetTraceId(t *testing.T) {
	t.Logf("%#v\n", GetTraceId("1"))
	SetTraceId("1", NewTraceId())
	t.Logf("%#v\n", GetTraceId("1"))
}

func TestGetGoIdTraceId(t *testing.T) {
	t.Logf("%#v\n", GetGoIdTraceId())
	time.Sleep(5 * time.Second)
	SetTraceId(goid.GetString(), NewTraceId())
	t.Logf("%#v\n", GetGoIdTraceId())
}

func TestHasTraceId(t *testing.T) {
	SetCache(NewFreeCache(10 * 1024 * 1024))
	t.Logf("%#v\n", HasTraceId("1"))
	SetTraceId("1", NewTraceId())
	time.Sleep(5 * time.Second)
	t.Logf("%#v\n", HasTraceId("1"))
}

func TestDelTraceId(t *testing.T) {
	SetTraceId("1", NewTraceId())
	t.Logf("%#v, %#v\n", HasTraceId("1"), GetTraceId("1"))
	DelTraceId("1")
	time.Sleep(5 * time.Second)
	t.Logf("%#v, %#v\n", HasTraceId("1"), GetTraceId("1"))
}
