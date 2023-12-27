package end

import (
	"blog-api/db_server/tables"
	"blog-api/global"
	"fmt"
	"math"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type ExposeGiftParams struct {
	IDs []int `form:"ids"`
}

// 将ExtendsTypes转换为字符串拼接形式
func extends2Str(s []tables.ExtendsType) string {
	if len(s) == 0 {
		return ""
	}
	str := ""
	for i, v := range s {
		if i == 0 {
			str = v.ExtendsName
		} else {
			str += fmt.Sprintf("、%s", v.ExtendsName)
		}
	}
	return str
}

// 同理，将礼物标签也转换为字符串拼接形式
func tags2Str(s []tables.GiftTag) string {
	if len(s) == 0 {
		return ""
	}
	str := ""
	for i, v := range s {
		if i == 0 {
			str = v.GiftTagName
		} else {
			str += fmt.Sprintf("、%s", v.GiftTagName)
		}
	}
	return str
}

func DownLoadGiftExcel(c *gin.Context) {
	db := global.GlobalDB

	var params ExposeGiftParams
	err := c.ShouldBindQuery(&params)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeGetParamsFailed,
			"message": "获取参数失败",
			"data":    nil,
		})
		return
	}

	if len(params.IDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeParamError,
			"message": "参数为空",
			"data":    nil,
		})
		return
	}

	var gifts []tables.KKGifts
	result := db.Model(&tables.KKGifts{}).Where("giftId IN ?", params.IDs).Find(&gifts)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeQueryFailed,
			"message": "查询失败",
			"data":    nil,
		})
		return
	}

	// 创建excel文件
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeCreateFileError,
				"message": "文件有误",
				"data":    nil,
			})
			return
		}
	}()
	// 设置列名
	colTitles := []string{"礼物ID", "礼物名称", "礼物类型", "类型拓展", "礼物标签", "礼物价值(秀币)", "礼物角标", "创建日期"}
	err = f.SetSheetRow("Sheet1", "A1", &colTitles)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeCreateFileError,
			"message": "文件有误",
			"data":    nil,
		})
		return
	}
	// 设置列宽度
	err = f.SetColWidth("Sheet1", "A", "H", 24)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeCreateFileError,
			"message": "设置有误",
			"data":    nil,
		})
		return
	}
	// 设置列名的样式，加粗，加背景色等
	colTitleStyleID, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{ // 对齐方式
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: excelize.Fill{ // 背景色
			Type:    "pattern",
			Color:   []string{"#fff2cc"},
			Pattern: 1,
		},
		Font: &excelize.Font{ // 字体
			Bold: true,
			Size: 16,
		},
		Border: []excelize.Border{ // 边框
			{Type: "top", Color: "D3D3D3", Style: 1},
			{Type: "bottom", Color: "D3D3D3", Style: 1},
			{Type: "left", Color: "D3D3D3", Style: 1},
			{Type: "right", Color: "D3D3D3", Style: 1},
		},
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeCreateFileError,
			"message": "设置有误",
			"data":    nil,
		})
		return
	}
	err = f.SetRowStyle("Sheet1", 1, 1, colTitleStyleID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeCreateFileError,
			"message": "设置有误",
			"data":    nil,
		})
		return
	}
	// 添加数据
	for i, v := range gifts {
		err := f.SetCellInt("Sheet1", fmt.Sprintf("A%v", i+2), v.GiftID)
		err = f.SetCellValue("Sheet1", fmt.Sprintf("B%v", i+2), v.GiftName)
		err = f.SetCellValue("Sheet1", fmt.Sprintf("C%v", i+2), v.GiftType)
		err = f.SetCellValue("Sheet1", fmt.Sprintf("D%v", i+2), extends2Str(v.ExtendsTypes))
		err = f.SetCellValue("Sheet1", fmt.Sprintf("E%v", i+2), tags2Str(v.GiftTags))
		err = f.SetCellValue("Sheet1", fmt.Sprintf("F%v", i+2), v.GiftValue)
		err = f.SetCellValue("Sheet1", fmt.Sprintf("G%v", i+2), v.CornerMarkName)
		err = f.SetCellValue("Sheet1", fmt.Sprintf("H%v", i+2), time.UnixMilli(v.CreateDate).Format("2006-01-02 15:04:05"))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeCreateFileError,
				"message": "设置有误",
				"data":    nil,
			})
			return
		}
	}
	// 给其他行设置样式
	sheetStyleID, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{ // 对齐方式
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Font: &excelize.Font{ // 字体
			Bold: false,
			Size: 12,
		},
		Border: []excelize.Border{ // 边框 #D3D3D3是电脑上创建的Excel的边框默认色，这里为了保持一致也设置该色。
			{Type: "top", Color: "D3D3D3", Style: 1},
			{Type: "bottom", Color: "D3D3D3", Style: 1},
			{Type: "left", Color: "D3D3D3", Style: 1},
			{Type: "right", Color: "D3D3D3", Style: 1},
		},
	})
	err = f.SetRowStyle("Sheet1", 2, int(math.Max(float64(len(gifts)+1), 2)), sheetStyleID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeCreateFileError,
			"message": "设置有误",
			"data":    nil,
		})
		return
	}

	// 添加响应头
	c.Header("Content-Type", "application/octet-stream")
	// 设置"Content-Disposition"是为了让浏览器去下载该文件，而不是只响应不下载
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="gifts_%v.xlsx"`, time.Now().UnixMilli()))
	// 浏览器会将header里的敏感数据都隐藏起来，前端就无法通过header获取后端定义的文件名，所以这里要暴露给浏览器
	c.Header("Access-Control-Expose-Headers", "Content-Type, Content-Disposition")

	if err := f.Write(c.Writer); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
}
