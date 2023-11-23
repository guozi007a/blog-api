package end

import (
	"net/http"
	"time"

	"blog-api/db_server/tables"
	"blog-api/global"

	"github.com/gin-gonic/gin"
)

type ParamsInfo struct {
	Branch string `json:"branch" form:"branch"`
	Name   string `json:"name" form:"name"`
	Tag    string `json:"tag" form:"tag"`
	Date   string `json:"date" form:"date"`
	Url    string `json:"url" form:"url"`
}

func AddActivity(c *gin.Context) {
	db := global.GlobalDB
	var paramsInfo ParamsInfo
	err := c.ShouldBind(&paramsInfo)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": err,
			"data":    nil,
		})
		return
	}
	// 查询已有的最大id
	tableName := &tables.ActivityListInfo{}
	// 之所以定义类型为*int64，而不用int，是因为值可能为null，*int64可以接收null，但是int不可以
	var maxId *int64
	result := db.Model(tableName).Select("max(id)").Scan(&maxId)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeQueryFailed,
			"message": "查询失败",
			"data":    nil,
		})
		panic(result.Error)
	}
	// fmt.Printf("maxId: %v, maxId == nil: %t\n", maxId, maxId == nil)
	var newId int
	if maxId == nil {
		newId = 1
	} else {
		newId = int(*maxId) + 1
	}
	info := tables.ActivityListInfo{
		ID:         newId,
		Branch:     paramsInfo.Branch,
		Name:       paramsInfo.Name,
		Tag:        paramsInfo.Tag,
		Date:       paramsInfo.Date,
		Url:        paramsInfo.Url,
		CreateDate: time.Now().UnixMilli(),
	}
	result = db.Create(&info)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeCreateDataFailed,
			"message": "新增失败",
			"data":    nil,
		})
		panic(result.Error)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data":    true,
	})
}
