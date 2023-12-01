package utils

import (
	"fmt"
	"math"
	"strings"
	"time"
)

var sp = []string{"M", "Z", "P", "E", "Q", "H", "X", "R", "K", "S"}

// 编码userId
func EncodeID(userId int) string {
	m := userId
	n := 0
	str := ""
	for m != 0 {
		n = m % 10
		str = sp[n] + str
		m = (m - n) / 10
	}

	return str
}

// 解码userId
func DecodeID(ss string) int {
	id := 0
	for k, v := range ss {
		for i, m := range sp {
			if string(v) == m {
				id += int(math.Pow10(len(ss)-k-1)) * i
				break
			}
		}
	}
	return id
}

// 创建sessionId
func CreateSessionID(userId int) string {
	return fmt.Sprintf("%s*%s", EncodeID(userId), Md5Str(fmt.Sprintf("%s%v", EncodeID(userId), time.Now().UnixMilli())))
}

// 解析sessionId 获取userId
func ParseSessionID(sessionId string) int {
	return DecodeID(strings.Split(sessionId, "*")[0])
}
