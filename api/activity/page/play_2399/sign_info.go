package play_2399

import (
	"net/http"

	"blog-api/global"

	"github.com/gin-gonic/gin"
)

func SignInfo(c *gin.Context) {
	// db := global.GlobalDB

	userId := c.Request.Header.Get("userId")
	token := c.Request.Header.Get("token")

	if userId == "" || token == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeAuthyLimited,
			"message": "无法签到",
			"data":    0,
		})
		return
	}
}
