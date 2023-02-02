package maths

import (
	"math/rand"
	"time"
)

func RandIntn(max int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(max)
}
