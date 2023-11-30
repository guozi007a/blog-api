package page

import (
	"net/http"

	"blog-api/db_server/tables"
	"blog-api/global"

	"github.com/gin-gonic/gin"
)

type LogoutInfo struct {
	UserId int `json:"userId"`
}

func Logout(c *gin.Context) {
	db := global.GlobalDB

	var info LogoutInfo
	err := c.ShouldBind(&info)
	if err != nil {
		panic(err)
	}
	var idInfo tables.IdInfo
	result := db.Model(&idInfo).Where("userId = ?", info.UserId).Select("isLogin", "token").Updates(tables.IdInfo{IsLogin: false, Token: ""})
	if result.Error != nil {
		panic(result.Error)
	}
	c.SetCookie("activityUserId", "", -1, "/", "", false, true)
	c.SetCookie("activityToken", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data":    true,
	})
}
