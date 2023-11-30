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

	cookieUserId, err := c.Cookie("activityUserId")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	// if cookieUserId == "10323" {
	// 	c.JSON(http.StatusOK, gin.H{})
	// 	return
	// }
	cookieToken, err := c.Cookie("activityToken")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	if cookieUserId == "" || cookieToken == "" {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	var userInfo tables.IdInfo
	uid, err := strconv.Atoi(cookieUserId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{})
		panic(err)
	}
	claims := plugins.ParseToken(cookieToken)
	expTime := claims.RegisteredClaims.ExpiresAt.Unix()
	now := time.Now().Unix()
	// token过期了
	if now > expTime {
		fmt.Println("token过期了")
		db.Model(&userInfo).Where("userId = ? AND token = ?", uid, cookieToken).Select("isLogin", "token").Updates(tables.IdInfo{IsLogin: false, Token: ""})
		c.SetCookie("activityUserId", "", -1, "/", global.ActivityCookieAllowOrigin(), false, true)
		c.SetCookie("activityToken", "", -1, "/", global.ActivityCookieAllowOrigin(), false, true)
		c.JSON(http.StatusOK, gin.H{})
		return
	} else {
		fmt.Println("token没过期")
		db.Model(&userInfo).Where("userId = ? AND token = ?", uid, cookieToken).Update("isLogin", true)
		fmt.Printf("uid: %v, cookieToken: %s\n", uid, cookieToken)
		result := db.Table(userInfo.TableName()).Where("userId = ? AND token = ?", uid, cookieToken).Find(&userInfo)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{})
			panic(result.Error)
		}
		fmt.Printf("id: %v, token: %s\n", userInfo.UserId, userInfo.Token)
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeOK,
			"message": "success",
			"data":    userInfo,
		})
	}
}
