package backstage

import (
	"blog-api/global"
	"net/http"

	"github.com/gin-gonic/gin"
	// "gorm.io/gorm/clause"
	"fmt"
	"io"

	"os"
	"regexp"
	"sort"
	"strconv"
)

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

// 文件正则分类
var regImage = `image/.*`      // 图片
var regAV = `(video|audio)/.*` // 音频

// 根据文件类型进行分类
func getCateByType(fileType string) string {
	reg1 := regexp.MustCompile(regImage)
	reg2 := regexp.MustCompile(regAV)

	switch {
	case reg1.MatchString(fileType):
		return global.SourceCateImg
	case reg2.MatchString(fileType):
		return global.SourceCateAV
	default:
		return global.SourceCateOther
	}
}

// 文件直传
func UploadFileDirect(c *gin.Context) {
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

	outPath := fmt.Sprintf("%s/%s/%s", global.StaticPath, getCateByType(fileType), file.Filename)

	err = c.SaveUploadedFile(file, outPath)
	if err != nil {
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
	err = os.RemoveAll(global.TempPath)
	if err != nil {
		panic(err)
	}

	// 再创建该目录，为存储下一个文件的切片做准备
	err = os.MkdirAll(global.TempPath, 0755)
	if err != nil {
		panic(err)
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

	chunkname := fmt.Sprintf("%v_%v%s", preFile.UID, currentChunk, regexp.MustCompile(`\..*`).FindString(preFile.Name))

	chunkfile, err := os.Create(fmt.Sprintf("%s/%s", global.TempPath, chunkname))
	if err != nil {
		panic(err)
	}

	chunkReader, err := chunk.Open()
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(chunkfile, chunkReader)
	if err != nil {
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
	entries, err := os.ReadDir(global.TempPath)
	if err != nil {
		panic(err)
	}

	// 从切片文件名中解析出切片序号，组成新的文件列表
	var chunks []ChunkSort
	for _, entry := range entries {
		num, err := strconv.Atoi(regexp.MustCompile(`.*_(\d+)\..*`).FindStringSubmatch(entry.Name())[1])
		if err != nil {
			panic(err)
		}
		chunks = append(chunks, ChunkSort{entry.Name(), num})
	}

	// 对列表重新按升序排序
	sort.Slice(chunks, func(i, j int) bool {
		return chunks[i].Num < chunks[j].Num
	})

	outFilePath := fmt.Sprintf("%s/%s/%s", global.StaticPath, getCateByType(preFile.Type), preFile.Name)

	file, err := os.Create(outFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 遍历合并
	for _, chunkItem := range chunks {
		tempFilePath := fmt.Sprintf("%s/%s", global.TempPath, chunkItem.Name)
		fi, err := os.Open(tempFilePath)
		if err != nil {
			panic(err)
		}

		_, err = io.Copy(file, fi)
		if err != nil {
			panic(err)
		}
		fi.Close()

		err = os.Remove(tempFilePath)
		if err != nil {
			panic(err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data": FinishFileData{
			UID:      preFile.UID,
			Progress: 100,
			Path:     fmt.Sprintf("http://localhost:4001/%s", outFilePath),
		},
	})
}
