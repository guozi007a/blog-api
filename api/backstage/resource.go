package backstage

import (
	"blog-api/global"
	"net/http"

	"github.com/gin-gonic/gin"

	// "gorm.io/gorm/clause"
	"fmt"
	// "io"
	"os"
)

type UploadFile struct {
	UID          string `json:"uid"`
	LastModified int    `json:"lastModified"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Type         string `json:"type"`
}

type FinishFileData struct {
	UID      string  `json:"uid"`
	Path     string  `json:"path"`
	Progress float32 `json:"progress"`
}

// 文件直传
func UploadFileDirect(c *gin.Context) {
	// db := global.GlobalDB

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "缺少必要参数：file",
		})
		return
	}

	uid := c.PostForm("uid")

	_, err1 := os.Stat(global.StaticPath)
	if os.IsNotExist(err1) {
		os.Mkdir(global.StaticPath, 0777)
	}

	outPath := fmt.Sprintf("%s/%s", global.StaticPath, file.Filename)

	err2 := c.SaveUploadedFile(file, outPath)
	if err2 != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeSaveFileError,
			"message": "文件保存失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data": FinishFileData{
			UID:      uid,
			Path:     fmt.Sprintf("http://localhost:4001/%s", outPath),
			Progress: 100,
		},
	})
}

// 图片预传
func PreUpload(c *gin.Context) {
	// db := global.GlobalDB

}
