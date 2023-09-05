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

// [{"id": "a", "key": "a", "content": "这是第1条日志哦~"}]

// 删除某个日期下的全部日志
func DeleteLogs(c *gin.Context) {

}

// 清空所有日志
func ClearAllLogs(c *gin.Context) {

}

type DateLog struct {
	ID      string `json:"id"`
	Key     string `json:"key"`
	Content string `json:"content"`
	// LogID   string `json:"log_id"` // 返回给前端时，不需要该字段了，就不写了
}

// 查询某个日期下的日志
func FindDateLogs(c *gin.Context) {
	db := global.GlobalDB

	date := c.Query("date")

	if date == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "缺少必要参数：date",
		})
		return
	}

	_, err := time.Parse("2016-01-02", date)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeFormatError,
			"message": "参数date格式错误",
		})
		return
	}

	var logs []DateLog

	result := db.Table("devlogs").Where("log_id = ?", date).Select("id", "key", "content").Find(&logs)

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeQueryFailed,
			"message": "查询失败",
		})
		panic(result.Error)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data":    logs,
	})
}

type Day struct {
	Date string `json:"date"`
}

type Dev struct {
	ID      string `json:"id"`
	Key     string `json:"key"`
	Content string `json:"content"`
}

// 查询所有日志
func FindAllLogs(c *gin.Context) {
	db := global.GlobalDB

	var days []Day

	// 这里使用降序查找，即最近的日志放在上面
	result := db.Table("datelogs").Order("date desc").Find(&days)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeQueryFailed,
			"message": "查询失败",
		})
		return
	}

	var devs []Dev
	var logs []map[string]interface{}

	for _, day := range days {
		result := db.Table("devlogs").Where("log_id = ?", day.Date).Find(&devs)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
			})
			return
		}

		// 这里需要先将day数据类型转换成map[string]interface{}类型，才能进行同类型拼接
		dayMap := make(map[string]interface{}, 1)
		dayMap["Date"] = day.Date
		logs = append(logs, dayMap)

		for _, dev := range devs {
			devMap := make(map[string]interface{}, 3)
			devMap["ID"] = dev.ID
			devMap["Key"] = dev.Key
			devMap["Content"] = dev.Content
			logs = append(logs, devMap)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "0",
		"message": "success",
		"data":    logs,
	})
}
