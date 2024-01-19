package goid

import (
	"strconv"

	"github.com/petermattis/goid"
)

func GetInt() int64 {
	return goid.Get()
}

func GetString() string {
	return strconv.FormatInt(goid.Get(), 10)
}
