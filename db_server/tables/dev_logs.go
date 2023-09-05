package tables

type DateLogs struct {
	// 这里直接用date做主键，利用主键的唯一性，避免日期重复
	Date string `json:"date" gorm:"primaryKey;size:20;comment:日期如2023-09-04"`
}

func (DateLogs) TableName() string {
	return "datelogs"
}

type DevLogs struct {
	ID       string   `json:"id" gorm:"primaryKey;comment:前端传入的日志唯一id"`
	Key      string   `json:"key" gorm:"comment:值同id"`
	Content  string   `json:"content" gorm:"comment:单条日志内容"`
	LogID    string   `json:"log_id" gorm:"size:20;comment:外键;"`                                             // error: BLOB/TEXT column 'log_id' used in key specification without a key length，需要设置size解决报错。这里logid的数据类型和size等，要和Date保持一致。
	DateLogs DateLogs `gorm:"foreignKey:LogID;references:Date;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // 设置外键 OnDelete的值设置为set null时，表示此时如果删除datelogs中的date，在devlogs中相关联的数据并不会被删除；如果设置为CASCADE，就会被删除，建立关联性。
}

func (DevLogs) TableName() string {
	return "devlogs"
}
