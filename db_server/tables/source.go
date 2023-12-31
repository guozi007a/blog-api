package tables

type SourceInfo struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	UID      string `json:"uid"`      // 上传时携带的uid
	Name     string `json:"name"`     // 带后缀的文件名
	Date     int64  `json:"date"`     // 上传时的时间戳，精确到毫秒，和前端保持一致。如果值用时间戳表示，这里的类型就需要写为int64，而不能用time.Time。如果用类型time.Time，这里的Date值可以写为time.Now()
	Category string `json:"category"` // 分类：图片image 音频av 其他other 分别放在不同目录下
	Size     int    `json:"size"`     // 文件大小
	Describe string `json:"describe"` // 文件说明和描述
	Temp     bool   `json:"temp"`     // 是否已转为临时文件，true表示当前文件为临时文件。当一个上传的文件被删除时，会先转为待删除的临时文件。可以选择将文件恢复或彻底删除。
}

func (SourceInfo) TableName() string {
	return "source_info"
}
