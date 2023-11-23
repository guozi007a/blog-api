package end

import (
	"net/http"

	"blog-api/db_server/tables"
	"blog-api/global"

	"github.com/gin-gonic/gin"
)

type RemoveParam struct {
	Branch string `json:"branch"`
}

func RemoveActivity(c *gin.Context) {
	db := global.GlobalDB
	var branch RemoveParam
	err := c.ShouldBind(&branch)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeFormatError,
			"message": "参数有误",
			"data":    nil,
		})
		return
	}
	result := db.Where("branch = ?", branch.Branch).Unscoped().Delete(&tables.ActivityListInfo{})
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
