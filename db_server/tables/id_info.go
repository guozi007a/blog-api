package tables

type IdInfo struct {
	UserId        int    `json:"userId" gorm:"primaryKey;column:userId;comment:用户id"`
	UserName      string `json:"username" gorm:"column:username;comment:用户名"`
	NickName      string `json:"nickName" gorm:"column:nickName;comment:用户昵称"`
	Avatar        string `json:"avatar" gorm:"column:avatar;comment:头像图片地址"`
	Password      string `json:"password" gorm:"column:password;comment:MD5加密后的密码"`
	Money         int64  `json:"money" gorm:"column:money;comment:秀币余额"`
	Coupon        int64  `json:"coupon" gorm:"column:coupon;comment:欢乐券余额"`
	Gender        int    `json:"gender" gorm:"column:gender;comment:性别编号"`
	Identity      int    `json:"identity" gorm:"column:identity;comment:身份类型编号"`
	UserLevel     int    `json:"userLevel" gorm:"column:userLevel;comment:用户等级编号"`
	ActorLevel    int    `json:"actorLevel" gorm:"column:actorLevel;comment:主播等级编号"`
	Talent        int    `json:"talent" gorm:"column:talent;comment:主播分区编号"`
	FamilyId      int    `json:"familyId" gorm:"column:familyId;comment:公会id，默认为10222"`
	FamilyName    string `json:"familyName" gorm:"column:familyName;comment:公会名称，默认为星互娱"`
	Birthday      string `json:"birthday" gorm:"column:birthday;comment:生日日期，如2023-11-22"`
	GenderName    string `json:"genderName" gorm:"column:genderName;comment:性别名称"`
	IdentityName  string `json:"identityName" gorm:"column:identityName;comment:身份类型名称"`
	TalentName    string `json:"talentName" gorm:"column:talentName;comment:分区名称"`
	CreateDate    int64  `json:"createDate" gorm:"column:createDate;comment:创建该id的时间戳;autoCreateTime:milli"`
	IsActor       bool   `json:"isActor" gorm:"column:isActor;comment:是否是主播身份"`
	IsLogin       bool   `json:"isLogin" gorm:"column:isLogin;comment:是否处于登录状态"`
	LastLoginDate int64  `json:"lastLoginDate" gorm:"column:lastLoginDate;comment:上一次主动登录的时间戳"`
	Token         string `json:"token" gorm:"comment:用于登录验证的token"`
}

func (IdInfo) TableName() string {
	return "id_info"
}
