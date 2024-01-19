package maths

import (
	"math/rand"
	"time"
)

// RandIntn 随机获取不大于max的随机整数
func RandIntn(max int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(max)
}
