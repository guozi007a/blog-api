package test_api

import (
	"net/http"
	"time"

	"blog-api/global"

	"github.com/gin-gonic/gin"
)

func CheckTime(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "测试时间",
		"data":    time.Now().UnixMilli(),
	})
}
