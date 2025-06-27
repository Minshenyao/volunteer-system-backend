package utils

import (
	"time"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
	TimeZone   = "Asia/Shanghai"
)

// FormatStr2Time 将时间字符串转换为time.Time，并保持本地时区
func FormatStr2Time(timeStr string) time.Time {
	// 设置本地时区
	loc, err := time.LoadLocation(TimeZone)
	if err != nil {
		loc = time.Local
	}

	// 解析时间字符串，将其解释为本地时区的时间
	t, err := time.ParseInLocation(TimeFormat, timeStr, loc)
	if err != nil {
		return time.Now() // 或者根据需求返回零值 time.Time{}
	}

	return t
}

// FormatTime2Str 将time.Time转换为指定格式的字符串
func FormatTime2Str(t time.Time) string {
	loc, err := time.LoadLocation(TimeZone)
	if err != nil {
		loc = time.Local
	}
	return t.In(loc).Format(TimeFormat)
}
