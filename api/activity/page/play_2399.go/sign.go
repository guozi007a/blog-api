package play2399

import (
	"fmt"
	"net/http"
	"strconv"

	"blog-api/db_server/tables"
	"blog-api/global"

	"time"

	"github.com/gin-gonic/gin"
)

type ActivityInfo struct {
	DateStart int64 `json:"dateStart"`
	DateEnd   int64 `json:"dateEnd"`
}

type DayChargeSum struct {
	MoneyTotal int64 `json:"moneyTotal"`
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
	var tb tables.ActivityListInfo
	var info ActivityInfo
	result := db.Table(tb.TableName()).Where("branch = ?", ACTIVITY_BRANCH).Find(&info)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeQueryFailed,
			"message": "查询失败",
			"data":    nil,
		})
		return
	}
	if info.DateStart == 0 || info.DateEnd == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeDataError,
			"message": "数据有误",
			"data":    nil,
		})
		return
	}
	if now < info.DateStart {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeNotStart,
			"message": "活动未开始",
			"data":    nil,
		})
		return
	}
	if now > info.DateEnd {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeFinished,
			"message": "活动已结束",
			"data":    nil,
		})
		return
	}

	startOfDayStr := fmt.Sprintf("%d-%d-%d 00:00:00", time.Now().Year(), int(time.Now().Month()), time.Now().Day())
	startDayTime, err := time.Parse("2006-01-02 15:04:05", startOfDayStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeDataError,
			"message": "数据有误",
			"data":    nil,
		})
		return
	}
	startOfDay := startDayTime.UnixMilli()
	endOfDay := startDayTime.Add(24 * time.Hour).UnixMilli()

	// 获取信息
	var maxId *int64
	var lastId int
	result = db.Model(&tables.Play_2399_Sign_List{}).Select("max(id)").Scan(&maxId)
	if maxId == nil {
		lastId = 1
	} else {
		lastId = int(*maxId) + 1
	}

	// 验证今日是否已获取签到资格(今日累计充值是否达标)
	var chargeInfo tables.ChargeInfo
	var dayTotalCharge DayChargeSum
	result = db.Table(chargeInfo.TableName()).Where("userId = ? AND date BETWEEN ? AND ?", uid, startOfDay, endOfDay).Select("sum(money)").Scan(&dayTotalCharge)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeQueryFailed,
			"message": "查询失败",
			"data":    nil,
		})
		return
	}
	if dayTotalCharge.MoneyTotal < int64(DAY_CHARGE_LIMIT) {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeAuthyLimited,
			"message": "今日签到资格不足",
			"data":    nil,
		})
		return
	}

	// 验证今日是否已签到
	var signInfo tables.Play_2399_Sign_List
	result = db.Table(signInfo.TableName()).Where("userId = ? AND status = ? AND date BETWEEN ? AND ?", uid, 2, startOfDay, endOfDay).Find(&signInfo)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeQueryFailed,
			"message": "查询失败",
			"data":    nil,
		})
		return
	}
	if signInfo.ID != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeRunAgain,
			"message": "今日已签到",
			"data":    nil,
		})
		return
	}

	// 检查签到奖池是否超限
	var limitInfo tables.Play_2399_Sign_List
	result = db.Table(limitInfo.TableName()).Last(&limitInfo)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeQueryFailed,
			"message": "查询失败",
			"data":    nil,
		})
		return
	}
	if limitInfo.TotalAwardMoney >= int64(SIGN_AWARD_POOL_LIMIT*1000) {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeAwardPoolOutLimited,
			"message": "奖池超限",
			"data":    nil,
		})
		return
	}

	// 完成签到

}
