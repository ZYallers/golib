package times

import (
	"github.com/ZYallers/golib/consts"
	"time"
)

//  TodayRemainSecond 获取当天剩余的秒数
//  @author Cloud|2021-12-12 12:49:41
//  @return time.Duration ...
func TodayRemainSecond() time.Duration {
	todayLast := time.Now().Format("2006-01-02") + " 23:59:59"
	lastTime, _ := time.ParseInLocation(consts.TimeFormat, todayLast, time.Local)
	return time.Duration(lastTime.Unix()-time.Now().Local().Unix()) * time.Second
}

//  NowTime 当前时间
//  @author Cloud|2021-12-13 09:40:41
//  @param format ...string ...
//  @return string ...
func NowTime(format ...string) string {
	f := consts.TimeFormat
	if len(format) > 0 {
		f = format[0]
	}
	return time.Now().Local().Format(f)
}
