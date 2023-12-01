package end

import (
	"net/http"
	"time"

	"blog-api/db_server/tables"
	"blog-api/global"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

type ChargeUser struct {
	UserId int `json:"userId"`
	PayId  int `json:"payId"`
	Money  int `json:"money"`
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
		panic(err)
	}
	if chargeUser.UserId == 0 || chargeUser.Money == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "缺少必要参数",
			"data":    nil,
		})
		return
	}
	var chargeInfo tables.ChargeInfo
	var maxId *int64
	result := db.Table(chargeInfo.TableName()).Select("max(id)").Scan(&maxId)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeCreateDataFailed,
			"message": "充值失败",
			"data":    nil,
		})
		panic(result.Error)
	}
	var userInfo tables.IdInfo
	result = db.Table(userInfo.TableName()).Where("userId = ?", chargeUser.UserId).Find(&userInfo)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeCreateDataFailed,
			"message": "充值失败",
			"data":    nil,
		})
		panic(result.Error)
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
			panic(result.Error)
		}
		payNick = payUser.NickName
	}
	result = db.Clauses(clause.OnConflict{DoNothing: true}).Create(tables.ChargeInfo{
		ID:       int(*maxId) + 1,
		UserId:   chargeUser.UserId,
		PayId:    chargeUser.PayId,
		NickName: userInfo.NickName,
		PayNick:  payNick,
		Money:    chargeUser.Money,
		Date:     time.Now().UnixMilli(),
	})
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeCreateDataFailed,
			"message": "充值失败",
			"data":    nil,
		})
		panic(result.Error)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data":    true,
	})
}
