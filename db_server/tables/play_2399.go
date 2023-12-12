package tables

type Play_2399_Sign_List struct {
	ID         int         `json:"id" gorm:"primaryKey;column:id;comment:唯一编号"`
	UserId     int         `json:"userId" gorm:"column:userId;comment:用户id"`
	Status     int         `json:"status" gorm:"comment:签到状态，0-未达到签到条件 1-有签到资格但还未签到 2-已签到"`
	Date       int64       `json:"date" gorm:"comment:签到时间，未签到时值为0，签到后值为签到时间的时间戳"`
	Awards     StringSlice `json:"awards" gorm:"comment:签到获得的奖励列表，未签到时为空列表"`
	CreateDate int64       `json:"createDate" gorm:"column:createDate;autoCreateTime:milli;comment:创建该记录的毫秒级时间戳"`
}

func (Play_2399_Sign_List) TableName() string {
	return "play_2399_sign_list"
}

type Play_2399_Card_Prize struct {
	Position  int    `json:"position" gorm:"comment:对应的卡片位置,1,2,3,4"`
	PrizeId   int    `json:"prizeId" gorm:"comment:礼物id,prizeId=0时表示未翻开;column:prizeId;default:0"`
	PrizeName string `json:"prizeName" gorm:"comment:礼物名称;column:prizeName"`
	UserId    int    `json:"userId" gorm:"index;comment:外键;column:userId"`
}

type Play_2399_Turn_Cards struct {
	UserId          int                    `json:"userId" gorm:"column:userId;primaryKey"`
	Count           int64                  `json:"count" gorm:"default:0;comment:当前轮次充值金额总数-秀币"`
	Round           int                    `json:"round" gorm:"default:1;comment:当前轮次"`
	UpdateRoundDate int64                  `json:"updateRoundDate" gorm:"comment:更新轮次的时间点;column:updateRoundDate"`
	Cards           []Play_2399_Card_Prize `json:"cards" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // 设置外键并做外键约束
}

func (Play_2399_Card_Prize) TableName() string {
	return "play_2399_card_prize"
}

func (Play_2399_Turn_Cards) TableName() string {
	return "play_2399_turn_cards"
}
