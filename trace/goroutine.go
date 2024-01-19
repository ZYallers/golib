package trace

import "github.com/ZYallers/golib/goid"

func Go(f func()) {
	go func(parentTraceId string) {
		defer func() { recover() }()
		goId := goid.GetString()
		defer DelTraceId(goId)
		SetTraceId(goId, parentTraceId)
		f()
	}(GetTraceId(goid.GetString()))
}
