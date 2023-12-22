package end

import (
	"net/http"

	"blog-api/db_server/tables"
	"blog-api/global"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func UpdateGift(c *gin.Context) {
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

	result := db.Unscoped().Delete(&tables.KKGifts{}, params.GiftID)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeDeleteFailed,
			"message": "删除失败",
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

	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&gift)

	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data":    true,
	})
}
