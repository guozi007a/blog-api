package utils

import (
	"blog-api/global"
	"fmt"
	"time"
)

/* 将time.Time格式的日期转换为global.COMMON_TIME_FORMAT形式的字符串 */
func FormatDateStr(t time.Time) string {
	return t.Format(global.COMMON_TIME_FORMAT)
}

/* 将当前时间转换为"2006-01-02"的字符串格式 */
func FormatNowShortStr() string {
	return time.Now().Format("2006-01-02")
}

/* 生成当前时间的毫秒级时间戳 */
func NowMilli() int64 {
	return time.Now().UnixMilli()
}

/* 计算出本地时间某一天的开始时间戳和结束时间戳(毫秒级) */
func DayMilli(t time.Time) (int64, int64) {
	loc, err := time.LoadLocation(global.SERVER_TIME_ZONE)
	if err != nil {
		return 0, 0
	}
	dayStr := t.Format(global.DAY_TIME_FORMAT)
	dayStartStr := fmt.Sprintf("%s 00:00:00", dayStr)

	ti, err := time.ParseInLocation(global.COMMON_TIME_FORMAT, dayStartStr, loc)
	if err != nil {
		return 0, 0
	}
	return ti.UnixMilli(), ti.Add(time.Hour * 24).UnixMilli()
}

/* 将时间字符串转换成本地时间戳 */
func LocMilli(timeStr string) int64 {
	loc, err := time.LoadLocation(global.SERVER_TIME_ZONE)
	if err != nil {
		return 0
	}
	ti, err := time.ParseInLocation(global.COMMON_TIME_FORMAT, timeStr, loc)
	if err != nil {
		return 0
	}
	return ti.UnixMilli()
}
