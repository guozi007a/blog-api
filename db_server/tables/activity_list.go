package tables

type ActivityListInfo struct {
	ID         int    `json:"id" gorm:"comment:活动编号;not null;unique;column:id"`
	Branch     string `json:"branch" gorm:"comment:活动分支,例play_2399;primaryKey"`
	Name       string `json:"name" gorm:"comment:活动名称,例感恩节活动"`
	Tag        string `json:"tag" gorm:"comment:活动类型,例如节日活动"`
	Date       string `json:"date" gorm:"comment:活动起始时间,例如11月23日 11:00-11月26日 24:00"`
	Url        string `json:"url" gorm:"comment:活动地址"`
	CreateDate int64  `json:"createDate" gorm:"comment:新增活动时间戳;column:createDate"`
}

func (ActivityListInfo) TableName() string {
	return "activity_list_info"
}
