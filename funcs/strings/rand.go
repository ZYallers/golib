package strings

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letter  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	number  = "0123456789"
)

// CreateCaptcha 生成指定位数的随机数
func CreateCaptcha(num int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%0"+strconv.Itoa(num)+"v", seededRand.Int31n(int32(math.Pow(10, float64(num)))))
}

// RandString 生成指定位数的随机字符串(大小写英文字母+0-9数字)
func RandString(length int) string {
	b := make([]byte, length)
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// RandLetter 生成指定位数的随机字母(大小写英文字母)
func RandLetter(length int) string {
	b := make([]byte, length)
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letter[seededRand.Intn(len(letter))]
	}
	return string(b)
}

// RandNumber 生成指定位数的随机数字(0-9数字)
func RandNumber(length int) string {
	b := make([]byte, length)
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = number[seededRand.Intn(len(number))]
	}
	return string(b)
}
