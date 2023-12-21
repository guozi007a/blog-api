package end

import (
	"blog-api/db_server/tables"
	"blog-api/global"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SearchGiftsQuery struct {
	GiftId       int    `json:"giftId" form:"giftId"`
	GiftName     string `json:"giftName" form:"giftName"`
	GiftTypeID   int    `json:"giftTypeId" form:"giftTypeId"`
	ExtendsID    int    `json:"extendsId" form:"extendsId"`
	GiftTagID    int    `json:"giftTagId" form:"giftTagId"`
	MinGiftValue int64  `json:"minGiftValue" form:"minGiftValue"`
	MaxGiftValue int64  `json:"maxGiftValue" form:"maxGiftValue"`
	Page         int    `json:"page" form:"page"`
	PageSize     int    `json:"pageSize" form:"pageSize"`
}

type SearchGiftsRes struct {
	GiftList []tables.KKGifts `json:"giftList"`
	Total    int64            `json:"total"`
}

func SearchGifts(c *gin.Context) {
	db := global.GlobalDB

	var params SearchGiftsQuery
	err := c.ShouldBindQuery(&params)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeGetParamsFailed,
			"message": "参数解析失败",
			"data":    nil,
		})
		return
	}
	if params.Page == 0 || params.PageSize == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "参数解析失败",
			"data":    nil,
		})
		return
	}

	var gifts []tables.KKGifts
	var count int64
	if params.GiftId == 0 && params.GiftName == "" && params.GiftTypeID == 0 && params.ExtendsID == 0 && params.GiftTagID == 0 && params.MinGiftValue == 0 && params.MaxGiftValue == 0 { /* 搜索全部 */
		result := db.Model(&tables.KKGifts{}).Order("createDate desc").Limit(params.PageSize).Offset((params.Page - 1) * params.PageSize).Preload("GiftTags").Preload("ExtendsTypes").Find(&gifts)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
		result = db.Model(&tables.KKGifts{}).Count(&count)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data": SearchGiftsRes{
			GiftList: gifts,
			Total:    count,
		},
	})
	return
}
