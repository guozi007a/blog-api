package end

import (
	"fmt"
	"net/http"

	"blog-api/db_server/tables"
	"blog-api/global"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

type AddGiftsParamsConfig struct {
	GiftID         int                  `json:"giftId" form:"giftId"`
	GiftName       string               `json:"giftName" form:"giftName"`
	GiftType       string               `json:"giftType" form:"giftType"`
	GiftTypeID     int                  `json:"giftTypeId" form:"giftTypeId"`
	ExtendsTypes   []tables.ExtendsType `json:"extendsTypes" form:"extendsTypes"`
	GiftTags       []tables.GiftTag     `json:"giftTags" form:"giftTags"`
	GiftValue      int64                `json:"giftValue" form:"giftValue"`
	GiftDescribe   string               `json:"giftDescribe" form:"giftDescribe"`
	CornerMarkID   int                  `json:"cornerMarkId" form:"cornerMarkId"`
	CornerMarkName string               `json:"cornerMarkName" form:"cornerMarkName"`
}

func AddGift(c *gin.Context) {
	db := global.GlobalDB

	var params AddGiftsParamsConfig
	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeGetParamsFailed,
			"message": "获取参数失败",
			"data":    nil,
		})
		return
	}
	if params.GiftID == 0 || params.GiftName == "" || params.GiftType == "" || params.GiftTypeID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "缺少必要参数",
			"data":    nil,
		})
		return
	}
	gift := tables.KKGifts{
		GiftID:         params.GiftID,
		GiftName:       params.GiftName,
		GiftType:       params.GiftType,
		GiftTypeID:     params.GiftTypeID,
		ExtendsTypes:   params.ExtendsTypes,
		GiftTags:       params.GiftTags,
		GiftValue:      params.GiftValue,
		GiftDescribe:   params.GiftDescribe,
		CornerMarkID:   params.CornerMarkID,
		CornerMarkName: params.CornerMarkName,
	}

	fmt.Printf("gift: %+v\n", gift)

	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&gift)

	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data":    true,
	})
}
