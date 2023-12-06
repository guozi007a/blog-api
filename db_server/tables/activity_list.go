package tables

type ActivityListInfo struct {
	ID          int    `json:"id" gorm:"comment:活动编号;not null;unique;column:id"`
	Branch      string `json:"branch" gorm:"comment:活动分支,例play_2399;primaryKey"`
	Name        string `json:"name" gorm:"comment:活动名称,例感恩节活动"`
	Tag         string `json:"tag" gorm:"comment:活动类型,例如节日活动"`
	Url         string `json:"url" gorm:"comment:活动地址"`
	CreateDate  int64  `json:"createDate" gorm:"comment:记录新增该活动时的时间戳;column:createDate"`
	DateStart   string `json:"dateStart" gorm:"column:dateStart;comment:活动开始时间，如2023-12-06 10:46:32"`
	DateEnd     string `json:"dateEnd" gorm:"column:dateEnd;comment:活动结束时间"`
	MoudleStart string `json:"moudleStart" gorm:"column:moudleStart;comment:活动中某个模块开始时间，如2023-12-06 10:46:32"`
	MoudleEnd   string `json:"moudleEnd" gorm:"column:moudleEnd;comment:活动中某个模块结束时间"`
}

func (ActivityListInfo) TableName() string {
	return "activity_list_info"
}
