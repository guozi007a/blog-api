package end

import (
	"database/sql"
	"net/http"
	"time"

	"blog-api/db_server/tables"
	"blog-api/global"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

type ChargeUser struct {
	UserId int   `json:"userId"`
	PayId  int   `json:"payId"`
	Money  int64 `json:"money"`
	Coupon int64 `json:"coupon"`
}

func Charge(c *gin.Context) {
	db := global.GlobalDB
	var chargeUser ChargeUser
	err := c.ShouldBind(&chargeUser)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeGetParamsFailed,
			"message": "参数解析失败",
			"data":    nil,
		})
		return
	}
	// 单次充值只能充值秀币或者欢乐券，不能同时充值两种，即二选一
	if chargeUser.UserId == 0 || (chargeUser.Money == 0 && chargeUser.Coupon == 0) {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "缺少必要参数",
			"data":    nil,
		})
		return
	}
	var chargeInfo tables.ChargeInfo
	// 这里定义sql.NullInt64为了处理表格中为空的情况
	var _maxId sql.NullInt64
	var maxId int
	result := db.Table(chargeInfo.TableName()).Select("max(id)").Scan(&_maxId)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeCreateDataFailed,
			"message": "充值失败",
			"data":    nil,
		})
		return
	}
	if _maxId.Valid { // 如果_maxId不为空
		maxId = int(_maxId.Int64)
	} else { // 为空
		maxId = 0
	}
	var userInfo tables.IdInfo
	result = db.Table(userInfo.TableName()).Where("userId = ?", chargeUser.UserId).Find(&userInfo)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeCreateDataFailed,
			"message": "充值失败",
			"data":    nil,
		})
		return
	}
	if userInfo.UserId == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeNotExist,
			"message": "用户不存在",
			"data":    nil,
		})
		return
	}
	payNick := ""
	if chargeUser.PayId != 0 {
		var payUser tables.IdInfo
		result := db.Table(payUser.TableName()).Where("userId = ?", chargeUser.PayId).Find(&payUser)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeCreateDataFailed,
				"message": "充值失败",
				"data":    nil,
			})
			return
		}
		if payUser.UserId == 0 {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeNotExist,
				"message": "用户不存在",
				"data":    nil,
			})
			return
		}
		payNick = payUser.NickName
	}
	result = db.Clauses(clause.OnConflict{DoNothing: true}).Create(tables.ChargeInfo{
		ID:       maxId + 1,
		UserId:   chargeUser.UserId,
		PayId:    chargeUser.PayId,
		NickName: userInfo.NickName,
		PayNick:  payNick,
		Money:    chargeUser.Money * 1000,
		Coupon:   chargeUser.Coupon * 1000,
		Date:     time.Now().UnixMilli(),
	})
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeCreateDataFailed,
			"message": "充值失败",
			"data":    nil,
		})
		return
	}
	// 更新用户列表中的金额数量
	chargeType := ""
	var count int64 = 0
	var total int64 = 0
	if chargeUser.Money != 0 {
		chargeType = "money"
		count = userInfo.Money
		total = count + chargeUser.Money*1000
	} else if chargeUser.Coupon != 0 {
		chargeType = "coupon"
		count = userInfo.Coupon
		total = count + chargeUser.Coupon*1000
	}
	result = db.Model(&tables.IdInfo{}).Where("userId = ?", chargeUser.UserId).Update(chargeType, total)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeCreateDataFailed,
			"message": "充值失败",
			"data":    nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data":    true,
	})
}
