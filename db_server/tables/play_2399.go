package tables

type Play_2399_Sign_List struct {
	ID         int      `json:"id" gorm:"primaryKey;column:id;comment:唯一编号"`
	UserId     int      `json:"userId" gorm:"column:userId;comment:用户id"`
	Status     int      `json:"status" gorm:"comment:签到状态，0-未达到签到条件 1-有签到资格但还未签到 2-已签到"`
	Date       int64    `json:"date" gorm:"comment:签到时间，未签到时值为0，签到后值为签到时间的时间戳"`
	Awards     []string `json:"awards" gorm:"comment:签到获得的奖励列表，未签到时为空列表"`
	CreateDate int64    `json:"createDate" gorm:"column:createDate;autoCreateTime:milli;comment:创建该记录的毫秒级时间戳"`
}

func (Play_2399_Sign_List) TableName() string {
	return "play_2399_sign_list"
}
