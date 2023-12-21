package end

import (
	"blog-api/db_server/tables"
	"blog-api/global"
	"fmt"
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
	} else if params.GiftId != 0 && params.GiftName == "" && params.GiftTypeID == 0 && params.ExtendsID == 0 && params.GiftTagID == 0 && params.MinGiftValue == 0 && params.MaxGiftValue == 0 { /* 按礼物id查找，此时不需要分页了，毕竟礼物id是唯一的，要么查不到，要么只能查到一条 */
		result := db.Model(&tables.KKGifts{}).Where("giftId = ?", params.GiftId).Order("createDate desc").Preload("GiftTags").Preload("ExtendsTypes").Find(&gifts)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
		count = int64(len(gifts))
	} else if params.GiftId == 0 && params.GiftName != "" && params.GiftTypeID == 0 && params.ExtendsID == 0 && params.GiftTagID == 0 && params.MinGiftValue == 0 && params.MaxGiftValue == 0 { /* 按礼物名称模糊匹配 */
		result := db.Model(&tables.KKGifts{}).Where("giftName Like ?", fmt.Sprintf("%%%s%%", params.GiftName)).Order("createDate desc").Limit(params.PageSize).Offset((params.Page - 1) * params.PageSize).Preload("GiftTags").Preload("ExtendsTypes").Find(&gifts)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
		result = db.Model(&tables.KKGifts{}).Where("giftName Like ?", fmt.Sprintf("%%%s%%", params.GiftName)).Count(&count)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
	} else if params.GiftId == 0 && params.GiftName == "" && params.GiftTypeID != 0 && params.ExtendsID == 0 && params.GiftTagID == 0 && params.MinGiftValue == 0 && params.MaxGiftValue == 0 { /* 按礼物类型查找 */
		result := db.Model(&tables.KKGifts{}).Where("giftTypeId = ?", params.GiftTypeID).Order("createDate desc").Limit(params.PageSize).Offset((params.Page - 1) * params.PageSize).Preload("GiftTags").Preload("ExtendsTypes").Find(&gifts)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
		result = db.Model(&tables.KKGifts{}).Where("giftTypeId = ?", params.GiftTypeID).Count(&count)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
	} else if params.GiftId == 0 && params.GiftName == "" && params.GiftTypeID == 0 && params.ExtendsID != 0 && params.GiftTagID == 0 && params.MinGiftValue == 0 && params.MaxGiftValue == 0 { /* 按礼物扩展类型查找 */
		/* 这里适用于外连接查询 */
		result := db.Model(&tables.KKGifts{}).Joins("left join extends_type on kk_gifts.giftId = extends_type.giftId and extends_type.extendsId = ?", params.ExtendsID).Order("createDate desc").Limit(params.PageSize).Offset((params.Page - 1) * params.PageSize).Preload("GiftTags").Preload("ExtendsTypes").Find(&gifts)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
		result = db.Model(&tables.ExtendsType{}).Where("extendsId = ?", params.ExtendsID).Count(&count)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
	} else if params.GiftId == 0 && params.GiftName == "" && params.GiftTypeID == 0 && params.ExtendsID == 0 && params.GiftTagID != 0 && params.MinGiftValue == 0 && params.MaxGiftValue == 0 { /* 按礼物标签查找 */
		result := db.Model(&tables.KKGifts{}).Joins("left join gift_tag on kk_gifts.giftId = gift_tag.giftId and gift_tag.giftTagId = ?", params.GiftTagID).Order("createDate desc").Limit(params.PageSize).Offset((params.Page - 1) * params.PageSize).Preload("GiftTags").Preload("ExtendsTypes").Find(&gifts)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
		result = db.Model(&tables.ExtendsType{}).Where("giftTagId = ?", params.GiftTagID).Count(&count)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
	} else if params.GiftId == 0 && params.GiftName == "" && params.GiftTypeID == 0 && params.ExtendsID == 0 && params.GiftTagID == 0 && params.MaxGiftValue != 0 { /* 按礼物价值区间查找 */
		result := db.Model(&tables.KKGifts{}).Where("giftValue BETWEEN ? AND ?", params.MinGiftValue, params.MaxGiftValue).Order("createDate desc").Limit(params.PageSize).Offset((params.Page - 1) * params.PageSize).Preload("GiftTags").Preload("ExtendsTypes").Find(&gifts)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
		result = db.Model(&tables.KKGifts{}).Where("giftValue BETWEEN ? AND ?", params.MinGiftValue, params.MaxGiftValue).Count(&count)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
	} else if params.GiftId != 0 && params.GiftName != "" && params.GiftTypeID == 0 && params.ExtendsID == 0 && params.GiftTagID != 0 && params.MinGiftValue == 0 && params.MaxGiftValue == 0 { /* 按礼物id和礼物名称查找 */
		result := db.Model(&tables.KKGifts{}).Where("giftId = ? AND giftName Like ?", params.GiftId, fmt.Sprintf("%%%s%%", params.GiftName)).Order("createDate desc").Preload("GiftTags").Preload("ExtendsTypes").Find(&gifts)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
		count = int64(len(gifts))
	} else if params.GiftId != 0 && params.GiftName == "" && params.GiftTypeID != 0 && params.ExtendsID == 0 && params.GiftTagID != 0 && params.MinGiftValue == 0 && params.MaxGiftValue == 0 { /* 按礼物id和礼物类型查找 */
		result := db.Model(&tables.KKGifts{}).Where("giftId = ? AND giftTypeId = ?", params.GiftId, params.GiftTypeID).Order("createDate desc").Preload("GiftTags").Preload("ExtendsTypes").Find(&gifts)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
		count = int64(len(gifts))
	} else if params.GiftId != 0 && params.GiftName == "" && params.GiftTypeID == 0 && params.ExtendsID != 0 && params.GiftTagID == 0 && params.MinGiftValue == 0 && params.MaxGiftValue == 0 { /* 按礼物id和礼物扩展类型查找 */
		result := db.Model(&tables.KKGifts{}).Where("giftId = ?", params.GiftId).Joins("left join extends_type on kk_gifts.giftId = extends_type.giftId and extends_type.extendsId = ?", params.ExtendsID).Order("createDate desc").Preload("GiftTags").Preload("ExtendsTypes").Find(&gifts)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
		count = int64(len(gifts))
	} else if params.GiftId != 0 && params.GiftName == "" && params.GiftTypeID == 0 && params.ExtendsID == 0 && params.GiftTagID != 0 && params.MinGiftValue == 0 && params.MaxGiftValue == 0 { /* 按礼物id和礼物标签查找 */
		result := db.Model(&tables.KKGifts{}).Where("giftId = ?", params.GiftId).Joins("left join gift_tag on kk_gifts.giftId = gift_tag.giftId and gift_tag.giftTagId = ?", params.GiftTagID).Order("createDate desc").Preload("GiftTags").Preload("ExtendsTypes").Find(&gifts)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
		count = int64(len(gifts))
	} else if params.GiftId != 0 && params.GiftName == "" && params.GiftTypeID == 0 && params.ExtendsID == 0 && params.GiftTagID == 0 && params.MaxGiftValue != 0 { /* 按礼物id和礼物价值区间查找 */
		result := db.Model(&tables.KKGifts{}).Where("giftId = ? AND giftValue BETWEEN ? AND ?", params.GiftId, params.MinGiftValue, params.MaxGiftValue).Order("createDate desc").Preload("GiftTags").Preload("ExtendsTypes").Find(&gifts)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
		count = int64(len(gifts))
	} else if params.GiftId == 0 && params.GiftName != "" && params.GiftTypeID != 0 && params.ExtendsID == 0 && params.GiftTagID == 0 && params.MinGiftValue == 0 && params.MaxGiftValue == 0 { /* 按礼物名称和礼物类型匹配 */
		result := db.Model(&tables.KKGifts{}).Where("giftName Like ? AND giftTypeId = ?", fmt.Sprintf("%%%s%%", params.GiftName), params.GiftTypeID).Order("createDate desc").Limit(params.PageSize).Offset((params.Page - 1) * params.PageSize).Preload("GiftTags").Preload("ExtendsTypes").Find(&gifts)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
		result = db.Model(&tables.KKGifts{}).Where("giftName Like ? AND giftTypeId = ?", fmt.Sprintf("%%%s%%", params.GiftName), params.GiftTypeID).Count(&count)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
	} else if params.GiftId == 0 && params.GiftName != "" && params.GiftTypeID == 0 && params.ExtendsID != 0 && params.GiftTagID == 0 && params.MinGiftValue == 0 && params.MaxGiftValue == 0 { /* 按礼物名称和礼物扩展类型查找 */
		result := db.Model(&tables.KKGifts{}).Where("giftName Like ?", fmt.Sprintf("%%%s%%", params.GiftName)).Joins("left join extends_type on kk_gifts.giftId = extends_type.giftId and extends_type.extendsId = ?", params.ExtendsID).Order("createDate desc").Limit(params.PageSize).Offset((params.Page - 1) * params.PageSize).Preload("GiftTags").Preload("ExtendsTypes").Find(&gifts)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
		result = db.Model(&tables.KKGifts{}).Where("giftName Like ?", fmt.Sprintf("%%%s%%", params.GiftName)).Joins("left join extends_type on kk_gifts.giftId = extends_type.giftId and extends_type.extendsId = ?", params.ExtendsID).Count(&count)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
	} else if params.GiftId == 0 && params.GiftName != "" && params.GiftTypeID == 0 && params.ExtendsID == 0 && params.GiftTagID != 0 && params.MinGiftValue == 0 && params.MaxGiftValue == 0 { /* 按礼物名称和礼物标签查找 */
		result := db.Model(&tables.KKGifts{}).Where("giftName Like ?", fmt.Sprintf("%%%s%%", params.GiftName)).Joins("left join gift_tag on kk_gifts.giftId = gift_tag.giftId and gift_tag.giftTagId = ?", params.GiftTagID).Order("createDate desc").Limit(params.PageSize).Offset((params.Page - 1) * params.PageSize).Preload("GiftTags").Preload("ExtendsTypes").Find(&gifts)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
		result = db.Model(&tables.KKGifts{}).Where("giftName Like ?", fmt.Sprintf("%%%s%%", params.GiftName)).Joins("left join gift_tag on kk_gifts.giftId = gift_tag.giftId and gift_tag.giftTagId = ?", params.GiftTagID).Count(&count)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
	} else if params.GiftId == 0 && params.GiftName != "" && params.GiftTypeID == 0 && params.ExtendsID == 0 && params.GiftTagID == 0 && params.MaxGiftValue != 0 { /* 按礼物名称和礼物价值区间查找 */
		result := db.Model(&tables.KKGifts{}).Where("giftName Like ? AND giftValue BETWEEN ? AND ?", fmt.Sprintf("%%%s%%", params.GiftName), params.MinGiftValue, params.MaxGiftValue).Order("createDate desc").Preload("GiftTags").Preload("ExtendsTypes").Find(&gifts)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
		result = db.Model(&tables.KKGifts{}).Where("giftName Like ? AND giftValue BETWEEN ? AND ?", fmt.Sprintf("%%%s%%", params.GiftName), params.MinGiftValue, params.MaxGiftValue).Count(&count)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
	} else if params.GiftId == 0 && params.GiftName == "" && params.GiftTypeID != 0 && params.ExtendsID != 0 && params.GiftTagID == 0 && params.MinGiftValue == 0 && params.MaxGiftValue == 0 { /* 按礼物类型和礼物扩展类型查找 */
		result := db.Model(&tables.KKGifts{}).Where("giftTypeId = ", params.GiftTypeID).Joins("left join extends_type on kk_gifts.giftId = extends_type.giftId and extends_type.extendsId = ?", params.ExtendsID).Order("createDate desc").Limit(params.PageSize).Offset((params.Page - 1) * params.PageSize).Preload("GiftTags").Preload("ExtendsTypes").Find(&gifts)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
		result = db.Model(&tables.KKGifts{}).Where("giftTypeId = ", params.GiftTypeID).Joins("left join extends_type on kk_gifts.giftId = extends_type.giftId and extends_type.extendsId = ?", params.ExtendsID).Count(&count)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
	} else if params.GiftId == 0 && params.GiftName == "" && params.GiftTypeID != 0 && params.ExtendsID == 0 && params.GiftTagID != 0 && params.MinGiftValue == 0 && params.MaxGiftValue == 0 { /* 按礼物类型和礼物标签查找 */
		result := db.Model(&tables.KKGifts{}).Where("giftTypeId = ", params.GiftTypeID).Joins("left join gift_tag on kk_gifts.giftId = gift_tag.giftId and gift_tag.giftTagId = ?", params.GiftTagID).Order("createDate desc").Limit(params.PageSize).Offset((params.Page - 1) * params.PageSize).Preload("GiftTags").Preload("ExtendsTypes").Find(&gifts)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
		result = db.Model(&tables.KKGifts{}).Where("giftTypeId = ", params.GiftTypeID).Joins("left join gift_tag on kk_gifts.giftId = gift_tag.giftId and gift_tag.giftTagId = ?", params.GiftTagID).Count(&count)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
	} else if params.GiftId == 0 && params.GiftName == "" && params.GiftTypeID != 0 && params.ExtendsID == 0 && params.GiftTagID == 0 && params.MaxGiftValue != 0 { /* 按礼物类型和礼物价值区间查找 */
		result := db.Model(&tables.KKGifts{}).Where("giftTypeId = ? AND giftValue BETWEEN ? AND ?", params.GiftTypeID, params.MinGiftValue, params.MaxGiftValue).Order("createDate desc").Preload("GiftTags").Preload("ExtendsTypes").Find(&gifts)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
		result = db.Model(&tables.KKGifts{}).Where("giftTypeId = ? AND giftValue BETWEEN ? AND ?", params.GiftTypeID, params.MinGiftValue, params.MaxGiftValue).Count(&count)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
	} else if params.GiftId == 0 && params.GiftName == "" && params.GiftTypeID == 0 && params.ExtendsID != 0 && params.GiftTagID != 0 && params.MinGiftValue == 0 && params.MaxGiftValue == 0 { /* 按礼物扩展类型和礼物标签查找 */
		result := db.Model(&tables.KKGifts{}).Joins("left join extends_type on kk_gifts.giftId = extends_type.giftId and extends_type.extendsId = ?", params.ExtendsID).Joins("left join gift_tag on kk_gifts.giftId = gift_tag.giftId and gift_tag.giftTagId = ?", params.GiftTagID).Order("createDate desc").Limit(params.PageSize).Offset((params.Page - 1) * params.PageSize).Preload("GiftTags").Preload("ExtendsTypes").Find(&gifts)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
		result = db.Model(&tables.KKGifts{}).Joins("left join extends_type on kk_gifts.giftId = extends_type.giftId and extends_type.extendsId = ?", params.ExtendsID).Joins("left join gift_tag on kk_gifts.giftId = gift_tag.giftId and gift_tag.giftTagId = ?", params.GiftTagID).Count(&count)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
	} else if params.GiftId == 0 && params.GiftName == "" && params.GiftTypeID == 0 && params.ExtendsID != 0 && params.GiftTagID == 0 && params.MaxGiftValue != 0 { /* 按礼物拓展和礼物价值区间查找 */
		result := db.Model(&tables.KKGifts{}).Where("giftValue BETWEEN ? AND ?", params.MinGiftValue, params.MaxGiftValue).Joins("left join extends_type on kk_gifts.giftId = extends_type.giftId and extends_type.extendsId = ?", params.ExtendsID).Order("createDate desc").Limit(params.PageSize).Offset((params.Page - 1) * params.PageSize).Preload("GiftTags").Preload("ExtendsTypes").Find(&gifts)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
		result = db.Model(&tables.KKGifts{}).Where("giftValue BETWEEN ? AND ?", params.MinGiftValue, params.MaxGiftValue).Joins("left join extends_type on kk_gifts.giftId = extends_type.giftId and extends_type.extendsId = ?", params.ExtendsID).Count(&count)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
	} else if params.GiftId == 0 && params.GiftName == "" && params.GiftTypeID == 0 && params.ExtendsID == 0 && params.GiftTagID != 0 && params.MaxGiftValue != 0 { /* 按礼物标签和礼物价值区间查找 */
		result := db.Model(&tables.KKGifts{}).Where("giftValue BETWEEN ? AND ?", params.MinGiftValue, params.MaxGiftValue).Joins("left join gift_tag on kk_gifts.giftId = gift_tag.giftId and gift_tag.giftTagId = ?", params.GiftTagID).Order("createDate desc").Limit(params.PageSize).Offset((params.Page - 1) * params.PageSize).Preload("GiftTags").Preload("ExtendsTypes").Find(&gifts)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
		result = db.Model(&tables.KKGifts{}).Where("giftValue BETWEEN ? AND ?", params.MinGiftValue, params.MaxGiftValue).Joins("left join gift_tag on kk_gifts.giftId = gift_tag.giftId and gift_tag.giftTagId = ?", params.GiftTagID).Count(&count)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			return
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeQueryFailed,
			"message": "暂不支持同时3种及以上的查询条件！",
			"data": SearchGiftsRes{
				GiftList: gifts,
				Total:    count,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data": SearchGiftsRes{
			GiftList: gifts,
			Total:    count,
		},
	})
}
