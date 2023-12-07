package play_2399

import (
	"net/http"
	"strconv"

	"blog-api/db_server/tables"
	"blog-api/global"

	"time"

	"blog-api/utils"

	"github.com/gin-gonic/gin"
)

type ActivityInfo struct {
	DateStart int64 `json:"dateStart"`
	DateEnd   int64 `json:"dateEnd"`
}

type DayChargeSum struct {
	MoneyTotal int64 `json:"moneyTotal"`
}

type UpdateInfo struct {
	Status int      `json:"status"`
	Date   int64    `json:"date"`
	Awards []string `json:"awards"`
}

func Sign(c *gin.Context) {
	db := global.GlobalDB

	userId := c.Request.Header.Get("userId")
	token := c.Request.Header.Get("token")
	if userId == "" || token == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeNotLogin,
			"message": "请先登录",
			"data":    nil,
		})
		return
	}
	uid, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeParamError,
			"message": "参数有误",
			"data":    nil,
		})
		return
	}

	now := time.Now().UnixMilli()

	// 验证是否在活动时间内
	var activityInfo ActivityInfo
	result := db.Where("branch = ?", ACTIVITY_BRANCH).Find(&activityInfo)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeQueryFailed,
			"message": "查询失败",
			"data":    nil,
		})
		return
	}
	if activityInfo.DateStart == 0 || activityInfo.DateEnd == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeDataError,
			"message": "数据有误",
			"data":    nil,
		})
		return
	}
	if now < activityInfo.DateStart {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeNotStart,
			"message": "活动未开始",
			"data":    nil,
		})
		return
	}
	if now > activityInfo.DateEnd {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeFinished,
			"message": "活动已结束",
			"data":    nil,
		})
		return
	}

	s, e := utils.DayMilli(time.Now())

	// 是否已签到 是否可以签到
	var signInfo tables.Play_2399_Sign_List
	result = db.Where("userId = ? AND createDate BETWEEN ? AND ?", uid, s, e).Find(&signInfo)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeQueryFailed,
			"message": "查询失败",
			"data":    nil,
		})
		return
	}
	if signInfo.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeAuthyLimited,
			"message": "未达到签到条件",
			"data":    nil,
		})
		return
	} else {
		switch signInfo.Status {
		case 0:
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeAuthyLimited,
				"message": "未达到签到条件",
				"data":    nil,
			})
		case 1:
			result = db.Model(&tables.Play_2399_Sign_List{}).Where("userId = ? AND createDate BETWEEN ? AND ?", uid, s, e).Updates(UpdateInfo{
				Status: 2,
				Date:   time.Now().UnixMilli(),
				Awards: []string{"a"},
			})
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeOK,
				"message": "success",
				"data":    true,
			})
		case 2:
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeRunAgain,
				"message": "今日已签到",
				"data":    nil,
			})
		}
	}
}
