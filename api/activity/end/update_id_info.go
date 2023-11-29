package end

import (
	"net/http"

	"blog-api/db_server/tables"
	"blog-api/global"

	"github.com/gin-gonic/gin"
)

type UpdateType struct {
	ParamType int         `json:"paramType"`
	UserId    int         `json:"userId"`
	Key       string      `json:"key"`
	Value     interface{} `json:"value"`
}

// 第一类 money, coupon, userLevel, actorLevel, familyId
// 第二类 nickName, avatar, familyName, birthday
// 第三类 gender-genderName, identity-identityName-isActor, talent-talentName
func UpdateIdInfo(c *gin.Context) {
	db := global.GlobalDB
	var updateInfo UpdateType
	if err := c.ShouldBind(&updateInfo); err != nil {
		panic(err)
	}
	upt := updateInfo.ParamType
	id := updateInfo.UserId
	if upt == 0 || id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "缺少必要参数",
			"data":    nil,
		})
		return
	}
	k := updateInfo.Key
	val := updateInfo.Value
	// val定义的是interface{}，即任意类型，但是实际使用时需要的是int和string类型，所以这里做一个断言
	var vi int
	var vs string
	// 这里需要注意的是，由于val没有规定类型，当值是数字时，默认类型是float64，而不是int
	if v, ok := val.(float64); ok {
		vi = int(v)
	} else if v, ok := val.(string); ok {
		vs = v
	}
	if upt == 1 {
		db.Model(&tables.IdInfo{}).Where("userId = ?", id).Update(k, vi)
	} else if upt == 2 {
		db.Model(&tables.IdInfo{}).Where("userId = ?", id).Update(k, vs)
	} else {
		switch k {
		case "gender":
			db.Model(&tables.IdInfo{}).Where("userId = ?", id).Updates(tables.IdInfo{Gender: vi, GenderName: genderList[vi-1]})
		case "identity":
			db.Model(&tables.IdInfo{}).Where("userId = ?", id).Select("Identity", "IdentityName", "IsActor").Updates(tables.IdInfo{Identity: vi, IdentityName: identityList[vi-1], IsActor: vi != 1})
		case "talent":
			db.Model(&tables.IdInfo{}).Where("userId = ?", id).Updates(tables.IdInfo{Talent: vi, TalentName: talentList[vi-1]})
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data":    true,
	})
}
