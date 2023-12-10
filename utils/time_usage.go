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

/* 计算出某一天的开始时间戳和结束时间戳(毫秒级) */
func DayMilli(t time.Time) (int64, int64) {
	dayStr := t.Format("2006-01-02")
	dayStartStr := fmt.Sprintf("%s 00:00:00", dayStr)

	ti, err := time.Parse(global.COMMON_TIME_FORMAT, dayStartStr)
	if err != nil {
		return 0, 0
	}
	return ti.UnixMilli(), ti.Add(time.Hour * 24).UnixMilli()
}

/*
根据给定的时间字符串，返回以北京时间为基准的UTC毫秒时间戳
该方法主要用于以北京时间进行固定的时间记录，比如活动开始时间，活动结束时间
*/
func CN_Milli(dateStr string) int64 {
	local, err := time.LoadLocation("Asia/Shanghai") // 加载中国东八区的时区
	if err != nil {
		return 0
	}

	ds, err := time.ParseInLocation(global.COMMON_TIME_FORMAT, dateStr, local) // 使用东八区的时区来解析时间字符串
	if err != nil {
		return 0
	}

	return ds.UTC().UnixMilli() // 输出对应于UTC的Unix毫秒时间戳
}
