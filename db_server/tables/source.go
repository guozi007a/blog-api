package tables

import (
	"time"
)

type SourceInfo struct {
	ID       int       `json:"id" gorm:"primaryKey"`
	UID      string    `json:"uid"`      // 上传时携带的uid
	Name     string    `json:"name"`     // 带后缀的文件名
	Date     time.Time `json:"date"`     // 上传的日期，精确到秒
	Category string    `json:"category"` // 分类：图片image 音频av 其他other 分别放在不同目录下
	Size     int       `json:"size"`     // 文件大小
	Describe string    `json:"describe"` // 文件说明和描述
}

func (SourceInfo) TableName() string {
	return "source_info"
}
