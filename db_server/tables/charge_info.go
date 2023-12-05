package tables

type ChargeInfo struct {
	ID       int    `json:"id" gorm:"comment:活动编号;primaryKey;column:id"`
	UserId   int    `json:"userId" gorm:"column:userId;comment:充值用户的id"`
	PayId    int    `json:"payId" gorm:"column:payId;comment:在支持给别人充值的情况下，花钱给UserId充值的那个人的id就是payId。一般情况下，payId是UserId，也就是自己给自己充值"`
	NickName string `json:"nickName" gorm:"column:nickName"`
	PayNick  string `json:"payNick" gorm:"column:payNick"`
	Money    int64  `json:"money" gorm:"comment:充值金额，单位是秀币，比如充值1元，Money就是1000"`
	Coupon   int64  `json:"coupon" gorm:"comment:充值的欢乐券"`
	Date     int64  `json:"date" gorm:"comment:充值时间戳(毫秒级)"`
}

func (ChargeInfo) TableName() string {
	return "charge_info"
}
