package end

import (
	"fmt"
	"net/http"

	"blog-api/db_server/tables"
	"blog-api/global"

	"strconv"

	"github.com/gin-gonic/gin"
)

type ResultList struct {
	List  []tables.ActivityListInfo `json:"list"`
	Total int64                     `json:"total"`
}

// 分页查询
func SearchActivityList(c *gin.Context) {
	db := global.GlobalDB
	pageSize := c.Query("pageSize")
	page := c.Query("page")
	if pageSize == "" || page == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "缺少必要参数",
			"data":    nil,
		})
		return
	}
	ps, err := strconv.Atoi(pageSize)
	p, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeFormatError,
			"message": "参数有误",
			"data":    nil,
		})
		return
	}
	fmt.Printf("ps: %v, p: %v\n", ps, p)
	var list []tables.ActivityListInfo
	result := db.Order("id desc").Limit(ps).Offset((p - 1) * ps).Find(&list)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeQueryFailed,
			"message": "查询失败",
			"data":    nil,
		})
		return
	}
	var count int64
	tableName := tables.ActivityListInfo{}
	result = db.Table(tableName.TableName()).Count(&count)
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
			Total: count,
		},
	})
}
