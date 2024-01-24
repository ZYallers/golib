package goid

import (
	"strconv"

	"github.com/petermattis/goid"
)

func Get() string {
	return GetString()
}

func GetInt() int {
	return int(goid.Get())
}

func GetInt64() int64 {
	return goid.Get()
}

func GetString() string {
	return strconv.FormatInt(goid.Get(), 10)
}
