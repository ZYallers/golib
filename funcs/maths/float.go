package maths

import (
	"math/rand"
	"strconv"
	"time"
)

// 区间范围内获取随机数
func RandFloat64(min, max float64, decimalNum int) float64 {
	rand.Seed(time.Now().UnixNano())
	limitFloat64 := rand.Float64()*float64(max-min)*100 + float64(min)*100
	limitStr := strconv.FormatFloat(limitFloat64/100, 'f', decimalNum, 64)
	rank, _ := strconv.ParseFloat(limitStr, 64)
	return rank
}
