package tables

type GiftTag struct {
	ID          int    `json:"id" gorm:"column:id;primaryKey"`
	GiftID      int    `json:"giftId" gorm:"index;comment:礼物ID;column:giftId"`
	GiftTagName string `json:"giftTagName" gorm:"comment:礼物标签，如活动礼物/年度礼物/战神礼物等;column:giftTagName"`
	GiftTagID   int    `json:"giftTagId" gorm:"comment:不同标签礼物对应的ID;column:giftTagId;index"`
}

func (GiftTag) TableName() string {
	return "gift_tag"
}

type ExtendsType struct {
	ID          int    `json:"id" gorm:"column:id;primaryKey"`
	GiftID      int    `json:"giftId" gorm:"index;comment:礼物ID;column:giftId"`
	ExtendsID   int    `json:"extendsId" gorm:"comment:拓展ID;column:extendsId;index"`
	ExtendsName string `json:"extendsName" gorm:"comment:拓展分类名称;column:extendsName"`
}

func (ExtendsType) TableName() string {
	return "extends_type"
}

type KKGifts struct {
	GiftID         int           `json:"giftId" gorm:"primaryKey;comment:礼物ID;column:giftId"`
	GiftName       string        `json:"giftName" gorm:"comment:礼物名称;column:giftName"`
	GiftType       string        `json:"giftType" gorm:"comment:礼物分类，如促销礼物，高级礼物，豪华礼物等;column:giftType"`
	GiftTypeID     int           `json:"giftTypeId" gorm:"comment:不同类别的礼物对应的不同ID;column:giftTypeId"`
	ExtendsTypes   []ExtendsType `json:"extendsTypes" gorm:"foreignKey:GiftID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	GiftTags       []GiftTag     `json:"giftTags" gorm:"foreignKey:GiftID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	GiftValue      int64         `json:"giftValue" gorm:"comment:礼物价格(秀币);column:giftValue"`
	CreateDate     int64         `json:"createDate" gorm:"comment:创建该记录的时间戳;autoCreateTime:milli;column:createDate"`
	UpdateDate     int64         `json:"updateDate" gorm:"comment:更新记录的时间戳;autoUpdateTime:milli;column:updateDate"`
	GiftDescribe   string        `json:"giftDescribe" gorm:"comment:礼物描述;column:giftDescribe"`
	CornerMarkID   int           `json:"cornerMarkId" gorm:"comment:角标类型;column:cornerMarkId"`
	CornerMarkName string        `json:"cornerMarkName" gorm:"comment:角标名称;column:cornerMarkName"`
}

func (KKGifts) TableName() string {
	return "kk_gifts"
}
