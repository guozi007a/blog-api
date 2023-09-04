package tables

type DateLogs struct {
	Date string `json:"date" gorm:"primaryKey;comment:日期如2023-09-04"`
}

func (DateLogs) TableName() string {
	return "dateLogs"
}

type DevLogs struct {
	ID       string   `json:"id" gorm:"primaryKey;comment:前端传入的日志唯一id"`
	Key      string   `json:"key" gorm:"comment:值同id"`
	Content  string   `json:"content" gorm:"comment:单条日志内容"`
	LogID    string   `json:"logId" gorm:"size:20;comment:外键;"`                                               // error: BLOB/TEXT column 'log_id' used in key specification without a key length，需要设置size解决报错。
	DateLogs DateLogs `gorm:"foreignKey:LogID;references:Date;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"` // 设置外键
}

func (DevLogs) TableName() string {
	return "devLogs"
}
