package backstage

import (
	"encoding/json"
	"fmt"
	"net/http"

	"blog-api/db_server/tables"
	"blog-api/global"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

var db = global.GlobalDB

// 发布日志
func PublishLogs(c *gin.Context) {
	// date := c.PostForm("date")
	logs := c.PostForm("logs")

	fmt.Println(logs)

	if logs == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "缺少必要参数：日志列表",
		})
		return
	}

	var logList []tables.DevLogs

	err := json.Unmarshal([]byte(logs), &logList)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeUnmarshalFailed,
			"message": fmt.Sprintf("数据解码失败：%s", err),
		})
		return
	}

	fmt.Printf("解码后：%+v", logList)

	// if date == "" { // 没有date，说明是发布日志，需要在后端生成一个日期，如2023-09-04
	publishDate := time.Now().Local().Format("2006-01-02")

	for _, log := range logList {
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&tables.DevLogs{
			ID:       log.ID,
			Key:      log.Key,
			Content:  log.Content,
			DateLogs: tables.DateLogs{Date: publishDate},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data":    true,
	})
	// }

}

// 删除某个日期下的全部日志
func DeleteLogs(c *gin.Context) {

}

// 清空所有日志
func ClearAllLogs(c *gin.Context) {

}

// 查询某个日期下的日志
func FindDayLogs(c *gin.Context) {

}

// 查询所有日志
func FindAllLogs(c *gin.Context) {

}

// 获取日志列表
func LogsList(c *gin.Context) {
	data := map[string]interface{}{
		"name":   "dilireba",
		"age":    18,
		"gender": "女",
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    data,
	})
}
