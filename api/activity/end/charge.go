package end

import (
	"net/http"

	"blog-api/db_server/tables"
	"blog-api/global"

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
		result := db.Model(&tables.IdInfo{}).Where("userId = ?", chargeUser.UserId).Update("money", gorm.Expr("money + ?", chargeUser.Count*1000))
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeUpdateFailed,
				"message": "更新失败",
				"data":    nil,
			})
			return
		}
	case "coupon":
		result := db.Model(&tables.IdInfo{}).Where("userId = ?", chargeUser.UserId).Update("coupon", gorm.Expr("coupon + ?", chargeUser.Count*1000))
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeUpdateFailed,
				"message": "更新失败",
				"data":    nil,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data":    true,
	})
}
