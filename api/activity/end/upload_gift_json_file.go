package end

import (
	"blog-api/db_server/tables"
	"blog-api/global"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

type GiftJsonParams struct {
	Gifts []tables.KKGifts `json:"gifts"`
}

func UploadGiftJsonFile(c *gin.Context) {
	db := global.GlobalDB

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeParamError,
			"message": "文件出错",
			"data":    nil,
		})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeParamError,
			"message": "文件出错",
			"data":    nil,
		})
		return
	}

	data, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeParamError,
			"message": "文件出错",
			"data":    nil,
		})
		return
	}

	var gifts GiftJsonParams
	err = json.Unmarshal(data, &gifts)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeParamError,
			"message": "文件出错",
			"data":    nil,
		})
		return
	}

	// fmt.Printf("上传的文件: %+v\n", gifts)
	uploadGifts := gifts.Gifts
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&uploadGifts)

	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "message",
		"data":    true,
	})
}
