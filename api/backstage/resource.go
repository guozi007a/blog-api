package backstage

import (
	"blog-api/global"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	// "gorm.io/gorm/clause"
	"fmt"
	"io"
	"mime"
	"os"
	"sort"
	"strconv"
	"strings"
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

type ChunkRes struct {
	UID  string `json:"uid"`
	Sort int    `json:"sort"`
}

// 切片排序
type ChunkSort struct {
	Name string
	Num  int
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

// 图片预传
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

	// 移除临时文件目录 防止掺杂其他切片文件
	err1 := os.RemoveAll("temp")
	if err1 != nil {
		panic(err1)
	}

	// 再创建该目录，为存储下一个文件的切片做准备
	err2 := os.MkdirAll("temp", 0755)
	if err2 != nil {
		panic(err2)
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
	chunk, err := c.FormFile("chunk")
	currentChunk := c.PostForm("currentChunk")
	if err != nil || currentChunk == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "缺少必要参数",
		})
		return
	}

	// 拼接切片临时名称 uid + _currentChunk + ext
	exts, err := mime.ExtensionsByType(preFile.Type)
	if err != nil {
		panic(err)
	}

	ext := exts[0]
	chunkname := fmt.Sprintf("%v_%v%s", preFile.UID, currentChunk, ext)

	chunkfile, err := os.Create("temp/" + chunkname)
	if err != nil {
		panic(err)
	}

	chunkReader, err := chunk.Open()
	if err != nil {
		panic(err)
	}

	_, err1 := io.Copy(chunkfile, chunkReader)
	if err1 != nil {
		panic(err)
	}
	// 这里一定要记得都用Close关闭，不然一直打开状态，在合并切片后，
	// 也会一直占用进程，无法删除切片文件
	chunkReader.Close()
	chunkfile.Close()

	intCurrentChunk, err := strconv.Atoi(currentChunk)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "0",
		"message": "success",
		"data": ChunkRes{
			UID:  preFile.UID,
			Sort: intCurrentChunk,
		},
	})
}

// 合并切片
func MergeChunks(c *gin.Context) {
	// 获取切片列表
	entries, err := os.ReadDir("temp")
	if err != nil {
		panic(err)
	}

	// 从切片文件名中解析出切片序号，组成新的文件列表
	var chunks []ChunkSort
	for _, entry := range entries {
		chunkPiece := strings.Split(entry.Name(), "_")[1]
		num, err := strconv.Atoi(strings.Split(chunkPiece, ".")[0])
		if err != nil {
			panic(err)
		}
		chunks = append(chunks, ChunkSort{entry.Name(), num})
	}

	// 对列表重新按升序排序
	sort.Slice(chunks, func(i, j int) bool {
		return chunks[i].Num < chunks[j].Num
	})

	// 创建static目录，有就创建，没有就忽略
	err1 := os.MkdirAll("static", 0755)
	if err1 != nil {
		panic(err1)
	}

	file, err := os.Create(fmt.Sprintf("static/%s", preFile.Name))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 遍历合并
	for _, chunkItem := range chunks {
		fi, err := os.Open(fmt.Sprintf("temp/%s", chunkItem.Name))
		if err != nil {
			panic(err)
		}

		_, err2 := io.Copy(file, fi)
		if err2 != nil {
			panic(err2)
		}
		fi.Close()

		err1 := os.Remove(fmt.Sprintf("temp/%s", chunkItem.Name))
		if err1 != nil {
			panic(err1)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data": FinishFileData{
			UID:      preFile.UID,
			Progress: 100,
			Path:     fmt.Sprintf("http://localhost:4001/static/%s", preFile.Name),
		},
	})
}
