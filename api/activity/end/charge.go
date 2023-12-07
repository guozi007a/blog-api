package end

import (
	"net/http"
	"time"

	"blog-api/db_server/tables"
	"blog-api/global"

	"blog-api/api/activity/page/play_2399"
	"blog-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ChargeUser struct {
	UserId int    `json:"userId"`
	PayId  int    `json:"payId"`
	Type   string `json:"type"`
	Count  int64  `json:"count"`
}

type DayChargeTotal struct {
	Total int64 `json:"total"`
}

func Charge(c *gin.Context) {
	db := global.GlobalDB
	var chargeUser ChargeUser
	err := c.BindJSON(&chargeUser)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeGetParamsFailed,
			"message": "参数解析失败",
			"data":    nil,
		})
		return
	}
	if chargeUser.UserId == 0 || (chargeUser.Type != "money" && chargeUser.Type != "coupon") {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "缺少必要参数",
			"data":    nil,
		})
		return
	}

	// 创建充值记录
	if chargeUser.PayId == 0 {
		info := tables.ChargeInfo{
			UserId: chargeUser.UserId,
			Type:   chargeUser.Type,
			Count:  chargeUser.Count * 1000,
		}
		result := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&info)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeCreateDataFailed,
				"message": "创建记录失败",
				"data":    nil,
			})
			return
		}
	} else {
		info := tables.ChargeInfo{
			UserId: chargeUser.UserId,
			PayId:  chargeUser.PayId,
			Type:   chargeUser.Type,
			Count:  chargeUser.Count * 1000,
		}
		result := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&info)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeCreateDataFailed,
				"message": "创建记录失败",
				"data":    nil,
			})
			return
		}
	}

	// 根据充值类型，更新用户信息中的秀币/欢乐券数量
	switch chargeUser.Type {
	case "money":
		result := db.Model(&tables.IdInfo{}).Where("userId = ?", chargeUser.UserId).Update("money", gorm.Expr("money + ?", chargeUser.Count))
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeUpdateFailed,
				"message": "更新失败",
				"data":    nil,
			})
			return
		}
	case "coupon":
		result := db.Model(&tables.IdInfo{}).Where("userId = ?", chargeUser.UserId).Update("coupon", gorm.Expr("coupon + ?", chargeUser.Count))
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeUpdateFailed,
				"message": "更新失败",
				"data":    nil,
			})
			return
		}
	}

	// 在活动时间内，更新play_2399的签到信息
	var activityInfo tables.ActivityListInfo
	result := db.Where("branch = ?", play_2399.ACTIVITY_BRANCH).Find(&activityInfo)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeQueryFailed,
			"message": "查询信息失败",
			"data":    nil,
		})
		return
	}
	if utils.NowMilli() > activityInfo.DateStart && utils.NowMilli() < activityInfo.DateEnd {
		var signInfo tables.Play_2399_Sign_List
		s, e := utils.DayMilli(time.Now())
		result = db.Where("userId = ? AND createDate BETWEEN ? AND ?", chargeUser.UserId, s, e).Find(&signInfo)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询信息失败",
				"data":    nil,
			})
			return
		}

		var totalCharge DayChargeTotal
		result = db.Model(&tables.ChargeInfo{}).Where("userId = ? AND date BETWEEN ? AND ? AND type = ?", chargeUser.UserId, s, e, "money").Select("sum(count)").Scan(&totalCharge)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询信息失败",
				"data":    nil,
			})
			return
		}

		status := 0
		if totalCharge.Total >= int64(play_2399.DAY_CHARGE_LIMIT*1000) {
			status = 1
		}

		if signInfo.ID == 0 {
			sign_info := tables.Play_2399_Sign_List{
				UserId: chargeUser.UserId,
				Status: status,
			}
			db.Clauses(clause.OnConflict{DoNothing: true}).Create(&sign_info)
		} else {
			if signInfo.Status == 0 && status == 1 {
				db.Model(&tables.Play_2399_Sign_List{}).Where("userId = ? AND createDate BETWEEN ? AND ?", chargeUser.UserId, s, e).Update("status", 1)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data":    true,
	})
}
