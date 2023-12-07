package tables

type ChargeInfo struct {
	ID     int    `json:"id" gorm:"comment:活动编号;primaryKey;column:id"`
	UserId int    `json:"userId" gorm:"column:userId;comment:充值用户的id"`
	PayId  int    `json:"payId" gorm:"column:payId;comment:在支持给别人充值的情况下，花钱给UserId充值的那个人的id就是payId。一般情况下，payId是UserId，也就是自己给自己充值"`
	Type   string `json:"type" gorm:"comment:充值类型，有两种:money&coupon"`
	Count  int64  `json:"count" gorm:"comment:充值数量，单位对应类型的秀币或欢乐券，和rmb的比例是1:1000"`
	Date   int64  `json:"date" gorm:"comment:充值时间戳(毫秒级);autoCreateTime:milli"`
}

func (ChargeInfo) TableName() string {
	return "charge_info"
}
