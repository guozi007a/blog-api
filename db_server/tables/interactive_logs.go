package tables

type InteractiveLogs struct {
	ID       int    `json:"id" gorm:"unique;not null;column:id;comment:唯一id"`
	Date     int64  `json:"date" gorm:"comment:时间戳"`
	UserId   int    `json:"userId" gorm:"primaryKey;comment:用户id;column:userId"`
	NickName string `json:"nickName" gorm:"comment:昵称;column:nickName"`
	Activity string `json:"activity" gorm:"comment:参与的活动名称"`
	Events   string `json:"events" gorm:"comment:交互事件，如点击"`
	Target   string `json:"target" gorm:"comment:交互的目标，如签到领奖按钮"`
	Action   string `json:"action" gorm:"comment:交互后的行为，比如获得、取消等"`
	Result   string `json:"result" gorm:"comment:交互的结果"`
	Content  string `json:"content" gorm:"comment:整个交互过程的描述，如这是一个用户啊 (10323) 在2023-11-26 21:35:44:324 参与感恩节回馈活动时， 点击了 签到领奖， 并获得了 狂欢盲盒。"`
}

func (InteractiveLogs) TableName() string {
	return "interactive_logs"
}
