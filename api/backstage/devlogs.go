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

type Log struct {
	ID      string `json:"id"`
	Key     string `json:"key"`
	Content string `json:"content"`
}

// 发布日志
func PublishLogs(c *gin.Context) {
	db := global.GlobalDB // 这一句初始化，一定要放置接口函数内，不然就panic

	logs := c.PostForm("logs")

	if logs == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "缺少必要参数：日志列表",
		})
		return
	}

	var logList []Log

	err := json.Unmarshal([]byte(logs), &logList)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeUnmarshalFailed,
			"message": fmt.Sprintf("数据解码失败：%s", err),
		})
		return
	}

	publishDate := time.Now().Local().Format("2006-01-02")

	for _, log := range logList {
		result := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&tables.DevLogs{
			ID:       log.ID,
			Key:      log.Key,
			Content:  log.Content,
			LogID:    publishDate,
			DateLogs: tables.DateLogs{Date: publishDate},
		})
		if result.Error != nil {
			panic(result.Error)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data":    true,
	})
}

// [{"id": "c", "key": "c", "content": "这是第3条日志哦~"}]

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
