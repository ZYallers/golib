package strings

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"math"
	"math/rand"
	"strconv"
	"time"
)

// 字符串首字母转成大写
func StrFirstToUpper(str string) string {
	if len(str) < 1 {
		return ""
	}
	strArray := []rune(str)
	if strArray[0] >= 97 && strArray[0] <= 122 {
		strArray[0] -= 32
	}
	return string(strArray)
}

// 字符串首字母转成小写
func StrFirstToLower(str string) string {
	if len(str) < 1 {
		return ""
	}
	strArr := []rune(str)
	if strArr[0] < 97 || strArr[0] > 122 {
		strArr[0] += 32
	}
	return string(strArr)
}

func MD5(str string) string {
	w := md5.New()
	_, _ = io.WriteString(w, str)
	return fmt.Sprintf("%x", w.Sum(nil))
}

// Nl2br nl2br()
// \n\r, \r\n, \r, \n
func Nl2br(str string, isXhtml bool) string {
	r, n, runes := '\r', '\n', []rune(str)
	var br []byte
	if isXhtml {
		br = []byte("<br />")
	} else {
		br = []byte("<br>")
	}
	skip := false
	length := len(runes)
	var buf bytes.Buffer
	for i, v := range runes {
		if skip {
			skip = false
			continue
		}
		switch v {
		case n, r:
			if (i+1 < length) && (v == r && runes[i+1] == n) || (v == n && runes[i+1] == r) {
				buf.Write(br)
				skip = true
				continue
			}
			buf.Write(br)
		default:
			buf.WriteRune(v)
		}
	}
	return buf.String()
}

// 生成指定位数的随机数
func CreateCaptcha(num int) string {
	return fmt.Sprintf("%0"+strconv.Itoa(num)+"v",
		rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(int32(math.Pow(10, float64(num)))))
}

// 截取字符串（支持中文）
func SubString(str string, begin, length int) string {
	rs := []rune(str)
	lth := len(rs)
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length

	if end > lth {
		end = lth
	}
	return string(rs[begin:end])
}
