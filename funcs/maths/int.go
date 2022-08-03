package maths

import (
	"math/rand"
	"time"
)

func RandIntn(max int) int {
	rad := rand.New(rand.NewSource(time.Now().Unix()))
	return rad.Intn(max)
}
