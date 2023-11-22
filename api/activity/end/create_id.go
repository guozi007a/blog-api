package end

import (
	"net/http"
	"time"

	"blog-api/db_server/tables"
	"blog-api/global"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"

	"crypto/md5"
	"fmt"
	"io"
)

type InfoType struct {
	UserId     int    `json:"userId" form:"userId"`
	UserName   string `json:"username" form:"username"`
	NickName   string `json:"nickName" form:"nickName"`
	Avatar     string `json:"avatar" form:"avatar"`
	Password   string `json:"password" form:"password"`
	Money      int32  `json:"money" form:"money"`
	Coupon     int32  `json:"coupon" form:"coupon"`
	Gender     int    `json:"gender" form:"gender"`
	Identity   int    `json:"identity" form:"identity"`
	UserLevel  int    `json:"userLevel" form:"userLevel"`
	ActorLevel int    `json:"actorLevel" form:"actorLevel"`
	Talent     int    `json:"talent" form:"talent"`
	FamilyId   int    `json:"familyId" form:"familyId"`
	FamilyName string `json:"familyName" form:"familyName"`
	Birthday   string `json:"birthday" form:"birthday"`
}

type InfoPatch struct {
	InfoType
	GenderName   string `json:"genderName"`
	IdentityName string `json:"identityName"`
	TalentName   string `json:"talentName"`
	CreateDate   int64  `json:"createDate"`
	IsActor      bool   `json:"isActor"`
}

var genderList [3]string = [3]string{"男", "女", "保密"}
var identityList [4]string = [4]string{"用户", "普通主播", "情感厅房主", "情感厅普通主播"}
var talentList [5]string = [5]string{"唱歌", "跳舞", "二次元", "搞笑", "无"}

// 数据加密
func secrete(str string) string {
	h := md5.New()
	io.WriteString(h, str)

	psd := fmt.Sprintf("%x", h.Sum(nil))

	return psd
}

func CreateId(c *gin.Context) {
	db := global.GlobalDB
	var info InfoType
	err := c.ShouldBind(&info)
	if err != nil {
		panic(err)
	}

	if info.UserId == 0 || info.UserName == "" || info.NickName == "" || info.Password == "" || info.Identity == 0 || info.Talent == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "缺少必要参数",
		})
		return
	}
	// 查询userId username nickName是否有重复
	var doubleInfo []InfoType
	tablename := tables.IdInfo{}
	result := db.Table(tablename.TableName()).Where("userId = ?", info.UserId).Or("username = ?", info.UserName).Or("nickName = ?", info.NickName).Find(&doubleInfo)
	if result.Error != nil {
		panic(result.Error)
	}
	if len(doubleInfo) != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeExist,
			"message": "用户id或用户名或昵称已存在",
			"data":    nil,
		})
		return
	}

	idInfo := &tables.IdInfo{
		UserId:       info.UserId,
		UserName:     info.UserName,
		NickName:     info.NickName,
		Avatar:       info.Avatar,
		Password:     secrete(info.Password),
		Money:        info.Money,
		Coupon:       info.Coupon,
		Gender:       info.Gender,
		Identity:     info.Identity,
		UserLevel:    info.UserLevel,
		ActorLevel:   info.ActorLevel,
		Talent:       info.Talent,
		FamilyId:     info.FamilyId,
		FamilyName:   info.FamilyName,
		Birthday:     info.Birthday,
		GenderName:   genderList[info.Gender-1],
		IdentityName: identityList[info.Identity-1],
		TalentName:   talentList[info.Talent-1],
		CreateDate:   time.Now().UnixMilli(),
		IsActor:      info.Identity != 1,
	}

	result = db.Clauses(clause.OnConflict{DoNothing: true}).Create(idInfo)
	if result.Error != nil {
		panic(result.Error)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data":    true,
	})
}
