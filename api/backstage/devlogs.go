package backstage

import (
	"net/http"

	"blog-api/db_server/tables"
	"blog-api/global"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

type Day struct {
	Date string `json:"date"`
}

type Log struct {
	ID      string `json:"id"`
	Key     string `json:"key"`
	Content string `json:"content"`
}

type PubLogParams struct {
	Logs []Log `json:"logs"`
}

type UpdateLogParams struct {
	Date string `json:"date"`
	Logs []Log  `json:"logs"`
}

// 发布日志
func PublishLogs(c *gin.Context) {
	db := global.GlobalDB // 这一句初始化，一定要放置接口函数内，不然就panic

	var params PubLogParams

	err1 := c.ShouldBind(&params)

	if err1 != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "缺少必要参数：日志列表",
		})
		return
	}

	logList := params.Logs

	publishDate := time.Now().Local().Format("2006-01-02")

	for _, log := range logList {
		// Cloauses用于说明在创建时对冲突的处理方式，这里是对冲突不做任何反应
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

// 更新某个日期的日志
func UpdateDateLogs(c *gin.Context) {
	db := global.GlobalDB

	var params UpdateLogParams

	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeGetParamsFailed,
			"message": "获取参数失败",
		})
		return
	}

	date := params.Date
	logs := params.Logs

	if date == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "缺少必要参数：date",
		})
		return
	}

	if _, err := time.Parse("2006-01-02", date); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeFormatError,
			"message": "参数date格式错误",
		})
		return
	}

	// 先删除
	db.Where("date = ?", date).Unscoped().Delete(&tables.DateLogs{})
	// 再重新创建
	for _, log := range logs {
		// Cloauses用于说明在创建时对冲突的处理方式，这里是对冲突不做任何反应
		result := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&tables.DevLogs{
			ID:       log.ID,
			Key:      log.Key,
			Content:  log.Content,
			LogID:    date,
			DateLogs: tables.DateLogs{Date: date},
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

// 删除某个日期下的全部日志
func DeleteDateLogs(c *gin.Context) {
	db := global.GlobalDB

	date := c.Query("date")

	if date == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "缺少必要参数：date",
		})
		return
	}

	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeFormatError,
			"message": "参数date格式错误",
		})
		return
	}

	db.Where("date = ?", date).Unscoped().Delete(&tables.DateLogs{})

	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data":    true,
	})
}

// 清空所有日志
func ClearAllLogs(c *gin.Context) {
	db := global.GlobalDB

	db.Where("1 = 1").Unscoped().Delete(&tables.DateLogs{})

	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data":    true,
	})
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

	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeFormatError,
			"message": "参数date格式错误",
		})
		return
	}

	var logs []Log

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

	var devs []Log
	var logs []interface{} // 定义一个元素是任意类型的空切片。因为需要把不同类型的数据拼接在一起

	for _, day := range days {
		result := db.Table("devlogs").Where("log_id = ?", day.Date).Find(&devs)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
			})
			return
		}

		logs = append(logs, day)

		for _, dev := range devs {
			logs = append(logs, dev)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "0",
		"message": "success",
		"data":    logs,
	})
}

// 删除单条日志
func DeleteOneLog(c *gin.Context) {
	db := global.GlobalDB

	id := c.Query("id")
	date := c.Query("date")

	if id == "" || date == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "缺少必要参数",
		})
		return
	}

	var logs []tables.DevLogs

	result := db.Table("devlogs").Where("log_id = ?", date).Find(&logs)
	if result.Error != nil {
		panic(result.Error)
	}

	if len(logs) > 1 {
		db.Where("id = ?", id).Unscoped().Delete(&tables.DevLogs{})
	} else {
		db.Where("date = ?", date).Unscoped().Delete(&tables.DateLogs{})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data":    true,
	})
}
