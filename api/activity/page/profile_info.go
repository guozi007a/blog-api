package page

import (
	"fmt"
	"net/http"

	"blog-api/db_server/tables"
	"blog-api/global"

	"blog-api/plugins"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetProfileInfo(c *gin.Context) {
	db := global.GlobalDB

	cookieUserId, err := c.Cookie("userId")
	if err != nil {
		panic(err)
	}
	cookieToken, err := c.Cookie("token")
	if err != nil {
		panic(err)
	}
	fmt.Printf("userId: %s\ntoken: %s\n", cookieUserId, cookieToken)
	if cookieUserId == "" || cookieToken == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeNotExist,
			"message": "not login",
			"data":    map[string]interface{}{},
		})
		return
	}
	var userInfo tables.IdInfo
	uid, err := strconv.Atoi(cookieUserId)
	if err != nil {
		panic(err)
	}
	claims := plugins.ParseToken(cookieToken)
	expTime := claims.RegisteredClaims.ExpiresAt.Unix()
	now := time.Now().Unix()
	// token过期了
	if now > expTime {
		db.Model(&userInfo).Where("userId = ? AND token = ?", uid, cookieToken).Select("isLogin", "token").Updates(tables.IdInfo{IsLogin: false, Token: ""})
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeNotExist,
			"message": "not login",
			"data":    map[string]interface{}{},
		})
		return
	} else {
		db.Model(&userInfo).Where("userId = ? AND token = ?", uid, cookieToken).Update("isLogin", true)
		result := db.Table(userInfo.TableName()).Where("userId = ? AND token = ?", uid, cookieToken).Find(&userInfo)
		if result.Error != nil {
			panic(result.Error)
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeOK,
			"message": "success",
			"data":    userInfo,
		})
	}
}
