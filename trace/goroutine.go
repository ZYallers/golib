package trace

import "github.com/ZYallers/golib/goid"

func Go(f func()) {
	go func(mainTraceId string) {
		defer func() { recover() }()
		id := goid.Get()
		defer DelTraceId(id)
		SetTraceId(id, mainTraceId)
		f()
	}(GetTraceId(goid.Get()))
}
