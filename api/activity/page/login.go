package page

import (
	"net/http"

	"blog-api/db_server/tables"
	"blog-api/global"
	"blog-api/utils"

	"blog-api/plugins"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginInfo struct {
	UserId   int    `json:"userId"`
	Password string `json:"password"`
}

type SessionInfo struct {
	ACTIVITY_SESSION_ID string
}

func Login(c *gin.Context) {
	db := global.GlobalDB

	var info LoginInfo
	err := c.ShouldBind(&info)
	if err != nil {
		panic(err)
	}
	if info.UserId == 0 || info.Password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "缺少必要参数",
			"data":    nil,
		})
		return
	}
	var userInfo tables.IdInfo
	result := db.Table(userInfo.TableName()).Where("userId = ? AND password = ?", info.UserId, utils.Md5Str(info.Password)).Find(&userInfo)
	if result.Error != nil {
		panic(result.Error)
	}
	if userInfo.UserId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeQueryFailed,
			"message": "账号或密码错误",
			"data":    nil,
		})
		return
	}
	now := time.Now().UnixMilli()
	_token := plugins.CreateToken(userInfo.UserId, userInfo.NickName)
	result = db.Model(&userInfo).Updates(tables.IdInfo{
		IsLogin:       true,
		Token:         _token,
		LastLoginDate: now,
	})
	if result.Error != nil {
		panic(result.Error)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data": SessionInfo{
			ACTIVITY_SESSION_ID: utils.CreateSessionID(info.UserId),
		},
	})
}
