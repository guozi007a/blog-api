package play_2399

import (
	"fmt"
	"net/http"
	"time"

	"blog-api/db_server/tables"
	"blog-api/global"

	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

type RoundTotal struct {
	Total int64 `json:"total"`
}

type CardPrizeInfo struct {
	Position  int    `json:"position"`
	PrizeId   int    `json:"prizeId"`
	PrizeName string `json:"prizeName"`
}

func CardInfo(c *gin.Context) {
	db := global.GlobalDB

	userId := c.Request.Header.Get("userId")
	token := c.Request.Header.Get("token")

	if userId == "" || token == "" { // 未传参
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeGetParamsFailed,
			"message": "",
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
	if activityInfo.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeNoActivityInfo,
			"message": "未查询到活动信息",
			"data":    nil,
		})
		return
	}
	if activityInfo.MoudleStart == 0 || activityInfo.MoudleEnd == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeActivityInfoError,
			"message": "活动信息错误",
			"data":    nil,
		})
		return
	}
	if time.Now().UnixMilli() < activityInfo.MoudleStart {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeNotStart,
			"message": "活动未开始",
			"data":    nil,
		})
		return
	}
	if time.Now().UnixMilli() > activityInfo.MoudleEnd {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeFinished,
			"message": "活动已结束",
			"data":    nil,
		})
		return
	}
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
	fmt.Printf("turn cards info: %+v\n", turnCardsInfo)
	if turnCardsInfo.UserId == 0 {
		// 活动开始到现在的充值金额
		var roundTotal RoundTotal
		result := db.Model(&tables.ChargeInfo{}).Select("sum(count) as total").Where("userId = ? AND date BETWEEN ? AND ?", uid, activityInfo.MoudleStart, time.Now().UnixMilli()).Find(&roundTotal)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}

		var cards []tables.Play_2399_Card_Prize
		for i := 1; i <= 4; i++ {
			cardItem := tables.Play_2399_Card_Prize{
				Position: i,
			}
			cards = append(cards, cardItem)
		}

		newCardsInfo := tables.Play_2399_Turn_Cards{
			UserId: uid,
			Count:  roundTotal.Total,
			Cards:  cards,
		}
		db.Clauses(clause.OnConflict{DoNothing: true}).Create(&newCardsInfo)

		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeOK,
			"message": "success",
			"data": tables.Play_2399_Turn_Cards{
				Count: roundTotal.Total / 1000,
				Round: 1,
				Cards: cards,
			},
		})
	} else {
		// 查询当前轮次的金额
		var roundTotal RoundTotal
		var findStart int64 = activityInfo.MoudleStart
		if turnCardsInfo.Round > 1 {
			findStart = turnCardsInfo.UpdateRoundDate
		}
		result := db.Model(&tables.ChargeInfo{}).Select("sum(count) as total").Where("userId = ? AND date BETWEEN ? AND ?", uid, findStart, time.Now().UnixMilli()).Find(&roundTotal)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeOK,
			"message": "success",
			"data": tables.Play_2399_Turn_Cards{
				Count: roundTotal.Total / 1000,
				Round: turnCardsInfo.Round,
				Cards: turnCardsInfo.Cards,
			},
		})
	}
}
