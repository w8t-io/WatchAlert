package tools

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/zeromicro/go-zero/core/logc"
	"strconv"
	"time"
)

// TimeTransformToWeek 时间转换成周
func TimeTransformToWeek(ct time.Time) string {
	// 获取当前时间
	currentDate := ct.Format("2006-01-02")
	// 解析日期字符串为时间对象
	date, err := time.Parse("2006-01-02", currentDate)
	if err != nil {
		logc.Error(context.Background(), fmt.Sprintf("Time Transform To Week failed, err: %s", err.Error()))
		return ""
	}
	return date.Weekday().String()
}

// TimeTransformToSeconds // 时间转换成秒
func TimeTransformToSeconds(ct time.Time) int {
	cs := ct.Hour()*3600 + ct.Minute()*60
	return cs
}

// FormatTimeToUTC 格式化为 UTC 时间
func FormatTimeToUTC(t int64) string {
	utcTime := time.Unix(t, 0).UTC()
	utcTimeString := utcTime.Format("2006-01-02T15:04:05.999Z")
	return utcTimeString
}

// ParserDuration 获取时间区间的开始时间
func ParserDuration(curTime time.Time, logScope int, timeType string) time.Time {
	duration, err := time.ParseDuration(strconv.Itoa(logScope) + timeType)
	if err != nil {
		logrus.Error(err.Error())
		return time.Time{}
	}
	startsAt := curTime.Add(-duration)
	return startsAt
}
