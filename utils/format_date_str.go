/* 将time.Time格式的日期转换为"2016-01-02 15:04:05"形式的字符串 */
package utils

import (
	"time"
)

func FormatDateStr(t time.Time) string {
	return t.Format("2016-01-02 15:04:05")
}
