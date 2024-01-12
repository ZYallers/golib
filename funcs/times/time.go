package times

import (
	"time"
)

// TodayRemainSecond 获取当天剩余的秒数
func TodayRemainSecond() time.Duration {
	todayLast := time.Now().Format("2006-01-02") + " 23:59:59"
	lastTime, _ := time.ParseInLocation("2006-01-02 15:04:05", todayLast, time.Local)
	return time.Duration(lastTime.Unix()-time.Now().Local().Unix()) * time.Second
}

// NowTime 指定format的当前时间字符串
func NowTime(format ...string) string {
	f := "2006-01-02 15:04:05"
	if len(format) > 0 {
		f = format[0]
	}
	return time.Now().Local().Format(f)
}
