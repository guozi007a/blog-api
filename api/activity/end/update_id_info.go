package end

import (
	"net/http"

	"blog-api/db_server/tables"
	"blog-api/global"

	"github.com/gin-gonic/gin"

	"strconv"
)

// 第一类 userId, money, coupon, userLevel, actorLevel, familyId
// 第二类 username, nickName, avatar, familyName, birthday
// 第三类 gender-genderName, identity-identityName-isActor, talent-talentName
func UpdateIdInfo(c *gin.Context) {
	db := global.GlobalDB
	paramType := c.PostForm("paramType")
	userId := c.PostForm("userId")
	if paramType == "" || userId == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "缺少必要参数",
			"data":    nil,
		})
		return
	}
	t, err := strconv.Atoi(paramType)
	if err != nil {
		panic(err)
	}
	id, err := strconv.Atoi(userId)
	if err != nil {
		panic(err)
	}
	k := c.PostForm("key")
	val := c.PostForm("value")
	if t == 1 {
		v, err := strconv.Atoi(val)
		if err != nil {
			panic(err)
		}
		db.Model(&tables.IdInfo{}).Where("userId = ?", id).Update(k, v)
	} else if t == 2 {
		db.Model(&tables.IdInfo{}).Where("userId = ?", id).Update(k, val)
	} else {
		v, err := strconv.Atoi(val)
		if err != nil {
			panic(err)
		}
		switch k {
		case "gender":
			db.Model(&tables.IdInfo{}).Where("userId = ?", id).Updates(&tables.IdInfo{Gender: v, GenderName: genderList[v-1]})
		case "identity":
			db.Model(&tables.IdInfo{}).Where("userId = ?", id).Updates(&tables.IdInfo{Identity: v, IdentityName: identityList[v-1], IsActor: v != 1})
		case "talent":
			db.Model(&tables.IdInfo{}).Where("userId = ?", id).Updates(&tables.IdInfo{Talent: v, TalentName: talentList[v-1]})
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data":    true,
	})
}
