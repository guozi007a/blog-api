package play_2399

import (
	"net/http"
	"time"

	"blog-api/db_server/tables"
	"blog-api/global"
	"blog-api/utils"

	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

type DayChargeTotal struct {
	Total int64 `json:"total"`
}

func SignInfo(c *gin.Context) {
	db := global.GlobalDB

	userId := c.Request.Header.Get("userId")
	token := c.Request.Header.Get("token")

	if userId == "" || token == "" { // 未传参
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeOK,
			"message": "",
			"data":    0,
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

	var activityInfo tables.ActivityListInfo
	result := db.Table(activityInfo.TableName()).Where("branch = ?", ACTIVITY_BRANCH).Find(&activityInfo)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeQueryFailed,
			"message": "查询失败",
			"data":    nil,
		})
		return
	}
	if time.Now().UnixMilli() < activityInfo.DateStart || time.Now().UnixMilli() > activityInfo.DateEnd { // 非活动时间
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeOK,
			"message": "",
			"data":    0,
		})
		return
	}

	s, e := utils.DayMilli(time.Now())
	var signInfo tables.Play_2399_Sign_List
	result = db.Table(signInfo.TableName()).Where("userId = ? AND createDate BETWEEN ? AND ?", uid, s, e).Find(&signInfo)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeQueryFailed,
			"message": "查询失败",
			"data":    nil,
		})
		return
	}

	var dayChargeTotal DayChargeTotal
	db.Model(&tables.ChargeInfo{}).Select("sum(count) as total").Where("userId = ? AND date BETWEEN ? AND ?", uid, s, e).Find(&dayChargeTotal)
	// fmt.Printf("充值总额：%v\n", dayChargeTotal.Total)
	status := 0
	if dayChargeTotal.Total >= int64(DAY_CHARGE_LIMIT*1000) {
		status = 1
	}

	if signInfo.ID == 0 { // 没查到记录
		newSignInfo := tables.Play_2399_Sign_List{ // 先初始化记录，再创建
			UserId: uid,
			Status: status,
		}
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&newSignInfo)
		if db.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeCreateDataFailed,
				"message": "创建失败",
				"data":    nil,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeOK,
			"message": "",
			"data":    status,
		})
		return
	} else { // 查到记录
		if signInfo.Status == 0 {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeOK,
				"message": "",
				"data":    status,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeOK,
				"message": "",
				"data":    signInfo.Status,
			})
		}
	}
}
