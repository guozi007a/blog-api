package end

import (
	"fmt"
	"net/http"

	"blog-api/db_server/tables"
	"blog-api/global"

	"github.com/gin-gonic/gin"
)

func SearchActivityByBranch(c *gin.Context) {
	db := global.GlobalDB
	branch := c.Query("branch")
	if branch == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "缺少必要参数",
			"data":    nil,
		})
		return
	}
	var list []tables.ActivityListInfo
	// 通过%xx%进行模糊查询
	result := db.Where("branch LIKE ?", fmt.Sprintf("%%%s%%", branch)).Find(&list)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeQueryFailed,
			"message": "查询失败",
			"data":    nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data": ResultList{
			List:  list,
			Total: int64(len(list)),
		},
	})
}
