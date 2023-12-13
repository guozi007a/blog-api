package play_2399

import (
	"net/http"
	"slices"
	"strconv"

	"blog-api/db_server/tables"
	"blog-api/global"

	"time"

	"blog-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ParamsConfig struct {
	Position int `json:"position"`
}

type RoundChargeTotal struct {
	Total int64 `json:"total"`
}

type TurnResultInfo struct {
	Position  int    `json:"position"`
	PrizeId   int    `json:"prizeId"`
	PrizeName string `json:"prizeName"`
}

func TurnCard(c *gin.Context) {
	db := global.GlobalDB

	var pos ParamsConfig
	err := c.BindJSON(&pos)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeParamError,
			"message": "参数有误",
			"data":    nil,
		})
		return
	}
	if slices.Contains[[]int, int]([]int{1, 2, 3, 4}, pos.Position) == false {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeParamError,
			"message": "参数有误",
			"data":    nil,
		})
		return
	}

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
	// 获取活动时间
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
	// 未开始
	if time.Now().UnixMilli() < activityInfo.MoudleStart {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeNotStart,
			"message": "活动未开始",
			"data":    nil,
		})
		return
	}
	// 已结束
	if time.Now().UnixMilli() > activityInfo.MoudleEnd {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeFinished,
			"message": "活动已结束",
			"data":    nil,
		})
		return
	}
	// 选择翻牌位置对应的金额
	cardMoney := CARDS_LIMIT[pos.Position-1]
	// 查询当前轮次信息
	var turnCardsInfo tables.Play_2399_Turn_Cards
	result = db.Model(&tables.Play_2399_Turn_Cards{}).Where("userId = ?", uid).Preload("Cards").Find(&turnCardsInfo)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeQueryFailed,
			"message": "查询失败",
			"data":    nil,
		})
		return
	}
	// 点击的该卡片已经翻过
	if turnCardsInfo.Cards[pos.Position-1].PrizeId != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeRunAgain,
			"message": "卡片已翻开",
			"data":    nil,
		})
		return
	}
	// 当前轮次已经翻开的牌子的数量
	turnedCount := 0
	for _, v := range turnCardsInfo.Cards {
		if v.PrizeId != 0 {
			turnedCount++
		}
	}
	// 如果在第一轮，就需要查询从活动开始到当前时间的充值总额
	if turnCardsInfo.Round == 1 {
		var roundChargeTotal RoundChargeTotal
		result = db.Model(&tables.ChargeInfo{}).Select("sum(count) as total").Where("userId = ? AND date BETWEEN ? AND ?", uid, activityInfo.DateStart, time.Now().UnixMilli()).Find(&roundChargeTotal)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
		// 如果充值金额不足
		if roundChargeTotal.Total < int64(cardMoney*1000) {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeChargeNotEnough,
				"message": "CodeChargeNotEnough",
				"data":    nil,
			})
			return
		}
		// 如果当前要翻开卡片的是本轮次最后一张卡片
		if turnedCount == 3 {
			// 轮次增加1 更新轮次的时间 新轮次的充值为0
			db.Model(&tables.Play_2399_Turn_Cards{}).Where("userId = ?", uid).Updates(map[string]interface{}{"round": gorm.Expr("round + ?", 1), "updateRoundDate": time.Now().UnixMilli(), "count": 0})
			// 将对应的卡片列表初始化
			db.Model(&tables.Play_2399_Card_Prize{}).Where("userId = ?", uid).Updates(map[string]interface{}{"prizeId": 0, "prizeName": ""})

			// 中奖
			award := CARDS_AWARDS[pos.Position-1][utils.RandInt(3)]
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeOK,
				"message": "success",
				"data": TurnResultInfo{
					Position:  pos.Position,
					PrizeId:   award.PrizeId,
					PrizeName: award.PrizeName,
				},
			})
			return
		} else { // 不是最后一张卡片
			// 更新对应位置卡片的奖励
			award := CARDS_AWARDS[pos.Position-1][utils.RandInt(3)]
			db.Model(&tables.Play_2399_Card_Prize{}).Where("userId = ? AND position = ?", uid, pos.Position).Updates(tables.Play_2399_Card_Prize{PrizeId: award.PrizeId, PrizeName: award.PrizeName})
			// 中奖
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeOK,
				"message": "success",
				"data": TurnResultInfo{
					Position:  pos.Position,
					PrizeId:   award.PrizeId,
					PrizeName: award.PrizeName,
				},
			})
			return
		}
	} else { // 不是第一轮
		// 当前轮次的充值金额就是本轮次开始到当前时间的充值金额
		var roundChargeTotal RoundChargeTotal
		result = db.Model(&tables.ChargeInfo{}).Select("sum(count) as total").Where("userId = ? AND date BETWEEN ? AND ?", uid, turnCardsInfo.UpdateRoundDate, time.Now().UnixMilli()).Find(&roundChargeTotal)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
		// 如果充值金额不足
		if roundChargeTotal.Total < int64(cardMoney*1000) {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeChargeNotEnough,
				"message": "CodeChargeNotEnough",
				"data":    nil,
			})
			return
		}
		// 如果当前要翻开卡片的是本轮次最后一张卡片
		if turnedCount == 3 {
			// 轮次增加1 更新轮次的时间 新轮次的充值为0
			db.Model(&tables.Play_2399_Turn_Cards{}).Where("userId = ?", uid).Updates(map[string]interface{}{"round": gorm.Expr("round + ?", 1), "updateRoundDate": time.Now().UnixMilli(), "count": 0})
			// 将对应的卡片列表初始化
			db.Model(&tables.Play_2399_Card_Prize{}).Where("userId = ?", uid).Updates(map[string]interface{}{"prizeId": 0, "prizeName": ""})

			// 中奖
			award := CARDS_AWARDS[pos.Position-1][utils.RandInt(3)]
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeOK,
				"message": "success",
				"data": TurnResultInfo{
					Position:  pos.Position,
					PrizeId:   award.PrizeId,
					PrizeName: award.PrizeName,
				},
			})
			return
		} else { // 不是最后一张卡片
			// 更新对应位置卡片的奖励
			award := CARDS_AWARDS[pos.Position-1][utils.RandInt(3)]
			db.Model(&tables.Play_2399_Card_Prize{}).Where("userId = ? AND position = ?", uid, pos.Position).Updates(tables.Play_2399_Card_Prize{PrizeId: award.PrizeId, PrizeName: award.PrizeName})
			// 中奖
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeOK,
				"message": "success",
				"data": TurnResultInfo{
					Position:  pos.Position,
					PrizeId:   award.PrizeId,
					PrizeName: award.PrizeName,
				},
			})
			return
		}
	}
}
