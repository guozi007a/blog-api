package tables

type GiftTags struct {
	GiftID    int    `json:"giftId" gorm:"index;comment:礼物ID;column:giftId"`
	GiftTag   string `json:"giftTag" gorm:"comment:礼物标签，如活动礼物，年度礼物，战神礼物等;column:giftTag"`
	GiftTagID int    `json:"giftTagId" gorm:"comment:不同标签礼物对应的ID;column:giftTagId;primaryKey"`
}

type ExtendsTypes struct {
	GiftID    int `json:"giftId" gorm:"index;comment:礼物ID;column:giftId"`
	ExtendsID int `json:"extendsId" gorm:"comment:拓展ID;column:extendsId;primaryKey"`
}

type KKGifts struct {
	GiftID       int            `json:"giftId" gorm:"primaryKey;comment:礼物ID;column:giftId"`
	GiftName     string         `json:"giftName" gorm:"comment:礼物名称;column:giftName"`
	GiftType     string         `json:"giftType" gorm:"comment:礼物分类，如促销礼物，高级礼物，豪华礼物等;column:giftType"`
	GiftTypeID   int            `json:"giftTypeId" gorm:"comment:不同类别的礼物对应的不同ID;column:giftTypeId"`
	ExtendsTypes []ExtendsTypes `json:"extendsTypes" gorm:"comment:类型扩展;column:extendsTypes;foreignKey:GiftID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	GiftTags     []GiftTags     `json:"giftTags" gorm:"foreignKey:GiftID;column:giftTags;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	GiftValue    int            `json:"giftValue" gorm:"comment:礼物价格(秀币);column:giftValue"`
	CreateDate   int64          `json:"createDate" gorm:"comment:创建该记录的时间戳;autoCreateTime:milli;column:createDate"`
	UpdateDate   int64          `json:"updateDate" gorm:"comment:更新记录的时间戳;autoUpdateTime:milli;column:updateDate"`
	RoomID       int            `json:"roomId" gorm:"comment:有些礼物可能只在指定的房间展示，所以会带有房间ID;column:roomId;default:0"`
	GiftDescribe string         `json:"giftDescribe" gorm:"comment:礼物描述;column:giftDescribe"`
}

func (GiftTags) TableName() string {
	return "gift_tags"
}
func (ExtendsTypes) TableName() string {
	return "extends_types"
}
func (KKGifts) TableName() string {
	return "kk_gifts"
}
