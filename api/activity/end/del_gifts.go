package end

import (
	"net/http"

	"blog-api/db_server/tables"
	"blog-api/global"

	"github.com/gin-gonic/gin"
)

type DelGiftsParams struct {
	IDs []int `json:"ids"` /* 接受字段为giftId组成的切片 */
}

func DelGifts(c *gin.Context) {
	db := global.GlobalDB

	var params DelGiftsParams
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeGetParamsFailed,
			"message": "获取参数失败",
			"data":    nil,
		})
		return
	}

	result := db.Unscoped().Delete(&tables.KKGifts{}, "giftId IN ?", params.IDs)
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
