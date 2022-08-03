package times

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// 日期转周几
func GetWeekDay(date string) string {
	var WeekDayMap = map[string]string{
		"Monday":    "周一",
		"Tuesday":   "周二",
		"Wednesday": "周三",
		"Thursday":  "周四",
		"Friday":    "周五",
		"Saturday":  "周六",
		"Sunday":    "周日",
	}
	day, _ := time.Parse("2006-01-02", date)
	dayInt := day.Weekday().String()
	return WeekDayMap[dayInt]
}

// 根据生日算年龄
func BirthdayToAge(birthday string) (age int) {
	if len(strings.Split(birthday, `-`)) == 2 {
		birthday = birthday + `-01` // 补上日期
	}

	sBirthday := strings.Split(birthday, `-`)
	if len(sBirthday) == 3 && len(birthday) != 10 {
		month0, _ := strconv.Atoi(sBirthday[1])
		day0, _ := strconv.Atoi(sBirthday[2])
		birthday = fmt.Sprintf(`%s-%02d-%02d`, sBirthday[0], month0, day0)
	}

	tt, err := time.ParseInLocation("2006-01-02", birthday, time.Local)
	if err != nil {
		return
	}

	year1, month1, day1 := tt.Date()
	year2, month2, day2 := time.Now().Date()

	age = year2 - year1
	if fmt.Sprintf(`%02d%02d`, month2, day2) < fmt.Sprintf(`%02d%02d`, month1, day1) {
		age = age - 1
	}
	return
}

// 获取两个时间的秒数差
func GetTimeRemainSeconds(startTimeStr, endTimeStr string) int {
	timeLayOut := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	startTime, _ := time.ParseInLocation(timeLayOut, startTimeStr, loc)

	theTime, _ := time.ParseInLocation(timeLayOut, endTimeStr, loc) // 使用模板在对应时区转化为time.time类型
	remain := math.Floor(theTime.Sub(startTime).Seconds())          // 获取距离当前的秒数

	return int(remain)
}

// 获取某一天的0点时间
func GetZeroTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// 获取本周周一的日期
func GetMondayOfWeek(t time.Time, fmtStr string) (dayStr string) {
	dayObj := GetZeroTime(t)
	if t.Weekday() == time.Monday {
		//修改hour、min、sec = 0后格式化
		dayStr = dayObj.Format(fmtStr)
	} else {
		offset := int(time.Monday - t.Weekday())
		if offset > 0 {
			offset = -6
		}
		dayStr = dayObj.AddDate(0, 0, offset).Format(fmtStr)
	}
	return
}

// 获取上周周日日期
func GetNextWeekMonday(t time.Time, fmtStr string) (day string, err error) {
	monday := GetMondayOfWeek(t, fmtStr)
	dayObj, err := time.Parse(fmtStr, monday)
	if err != nil {
		return
	}
	day = dayObj.AddDate(0, 0, +7).Format(fmtStr)
	return
}

// 通过ISOWeek翻转得到周的日期时间
// @see https://blog.csdn.net/pingD/article/details/60964306
func FirstDayOfISOWeek(year int, week int, timezone *time.Location) time.Time {
	date := time.Date(year, 0, 0, 0, 0, 0, 0, timezone)
	isoYear, isoWeek := date.ISOWeek()
	for date.Weekday() != time.Monday { // iterate back to Monday
		date = date.AddDate(0, 0, -1)
		isoYear, isoWeek = date.ISOWeek()
	}
	for isoYear < year { // iterate forward to the first day of the first week
		date = date.AddDate(0, 0, 1)
		isoYear, isoWeek = date.ISOWeek()
	}
	for isoWeek < week { // iterate forward to the first day of the given week
		date = date.AddDate(0, 0, 1)
		isoYear, isoWeek = date.ISOWeek()
	}
	return date
}
