package end

import (
	"blog-api/db_server/tables"
	"blog-api/global"
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
)

func DownLoadGiftsJSON(c *gin.Context) {
	db := global.GlobalDB

	var gifts []tables.KKGifts
	result := db.Model(&tables.KKGifts{}).Find(&gifts)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeQueryFailed,
			"message": "查询失败",
			"data":    nil,
		})
		return
	}

	giftsJSON := GiftJsonParams{
		Gifts: gifts,
	}
	jsonData, err := json.MarshalIndent(giftsJSON, "", "    ")
	if err != nil {
		c.String(http.StatusOK, "文件解析错误")
		return
	}

	reader := bytes.NewReader(jsonData)
	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="gifts_%v.json"`, time.Now().UnixMilli()),
		// 浏览器会将header里的敏感数据都隐藏起来，前端就无法通过header获取后端定义的文件名，所以这里要暴露给浏览器
		"Access-Control-Expose-Headers": "Content-Type, Content-Disposition",
	}
	// 设置了"application/octet-stream"之后，就会在响应头中自动加入`content-type: "application/octet-stream"`
	c.DataFromReader(http.StatusOK, int64(len(jsonData)), "application/octet-stream", reader, extraHeaders)
}
