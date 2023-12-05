package end

import (
	"net/http"

	"blog-api/db_server/tables"
	"blog-api/global"

	"github.com/gin-gonic/gin"
)

type ParamConfig struct {
	ID int `json:"id"`
}

func ChargeDel(c *gin.Context) {
	db := global.GlobalDB

	var param ParamConfig
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeDeleteFailed,
			"message": "参数有误",
			"data":    nil,
		})
		return
	}
	result := db.Where("id = ?", param.ID).Unscoped().Delete(&tables.ChargeInfo{})
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeDeleteFailed,
			"message": "删除失败",
			"data":    nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data":    true,
	})
}
