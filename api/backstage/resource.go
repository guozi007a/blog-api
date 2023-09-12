package backstage

import (
	"blog-api/global"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	// "gorm.io/gorm/clause"
	"fmt"
	// "io"
	"os"
)

// 存储到db中的文件
type DBFile struct {
	ID       int       `json:"id" gorm:"primaryKey"`
	UID      string    `json:"uid"`
	Path     string    `json:"path"`
	Date     time.Time `json:"date"`
	Category string    `json:"category"` // 分类：图片 音频 其他
	Size     int       `json:"size"`
}

type FinishFileData struct {
	UID      string  `json:"uid"`
	Path     string  `json:"path"`
	Progress float32 `json:"progress"`
}

type PreFile struct {
	UID  string `json:"uid"`
	Size int64  `json:"size"` // 这里的size要设置为int64类型，因为file.Size就是int64类型的，省去类型转换
	Name string `json:"name"`
	Type string `json:"type"` // 文件的mime类型，例如"image/png"
	Sort int    `json:"sort"` // 文件已经完成的上传片段数量
}

var preFile PreFile

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

// 图片预传 暂时略
func PreUpload(c *gin.Context) {
	// db := global.GlobalDB
	file, err := c.FormFile("file")
	uid := c.PostForm("uid")
	fileType := c.PostForm("type")
	if err != nil || uid == "" || fileType == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "缺少必要参数",
		})
		return
	}

	preFile = PreFile{
		UID:  uid,
		Name: file.Filename,
		Size: file.Size,
		Type: fileType,
		Sort: 0,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data": map[string]interface{}{
			"uid":  uid,
			"sort": 0,
		},
	})
}

// 上传切片
func UploadChunk(c *gin.Context) {
	// chunk, err := c.FormFile("chunk")
	// currentChunk := c.PostForm("currentChunk")
	// if err != nil || currentChunk == "" {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"code":    global.CodeLackRequired,
	// 		"message": "缺少必要参数",
	// 	})
	// 	return
	// }
}
