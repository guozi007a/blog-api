package backstage

import (
	"blog-api/db_server/tables"
	"blog-api/global"
	"net/http"
	"time"

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
	UID      string `json:"uid"`
	Size     int64  `json:"size"` // 这里的size要设置为int64类型，因为file.Size就是int64类型的，省去类型转换
	Name     string `json:"name"`
	Type     string `json:"type"`     // 文件的mime类型，例如"image/png"
	Sort     int    `json:"sort"`     // 文件已经完成的上传片段数量
	Describe string `json:"describe"` // 对预传文件的文字描述
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

type FileListData struct {
	Count    int64               `json:"count"`
	FileList []tables.SourceInfo `json:"fileList"`
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
	db := global.GlobalDB

	file, err := c.FormFile("file")
	uid := c.PostForm("uid")
	fileType := c.PostForm("type")
	describe := c.PostForm("describe")
	if err != nil || uid == "" || fileType == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "缺少必要参数",
		})
		return
	}

	var info tables.SourceInfo

	re := db.Table("source_info").Where("name = ?", file.Filename).Find(&info)
	if re.Error != nil {
		panic(re.Error)
	}

	// 根据文件名称查找，如果文件名已存在，判定为文件已存在，不支持上传，以免错误覆盖
	if info.ID != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeExist,
			"message": "上传内容已存在",
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

	var sourceInfo = tables.SourceInfo{
		UID:      uid,
		Name:     file.Filename,
		Date:     time.Now().UnixMilli(),
		Category: getCateByType(fileType),
		Size:     int(file.Size),
		Describe: describe,
		Temp:     false,
	}

	result := db.Create(&sourceInfo)
	if result.Error != nil {
		panic(result.Error)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data": FinishFileData{
			UID:      uid,
			Path:     fmt.Sprintf("%s/%s", global.GlobalOrigin(), outPath),
			Progress: 100,
		},
	})
}

// 图片预传
func PreUpload(c *gin.Context) {
	db := global.GlobalDB
	file, err := c.FormFile("file")
	uid := c.PostForm("uid")
	fileType := c.PostForm("type")
	describe := c.PostForm("describe")
	if err != nil || uid == "" || fileType == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "缺少必要参数",
		})
		return
	}

	var info tables.SourceInfo

	result := db.Table("source_info").Where("name = ?", file.Filename).Find(&info)
	if result.Error != nil {
		panic(result.Error)
	}

	if info.ID != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeExist,
			"message": "上传内容已存在",
		})
		return
	}

	preFile = PreFile{
		UID:      uid,
		Name:     file.Filename,
		Size:     file.Size,
		Type:     fileType,
		Sort:     0,
		Describe: describe,
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
	db := global.GlobalDB
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

	var sourceInfo = tables.SourceInfo{
		UID:      preFile.UID,
		Name:     preFile.Name,
		Date:     time.Now().UnixMilli(),
		Category: getCateByType(preFile.Type),
		Size:     int(preFile.Size),
		Describe: preFile.Describe,
		Temp:     false,
	}

	result := db.Create(&sourceInfo)
	if result.Error != nil {
		panic(result.Error)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data": FinishFileData{
			UID:      preFile.UID,
			Progress: 100,
			Path:     fmt.Sprintf("%s/%s", global.GlobalOrigin(), outFilePath),
		},
	})
}

// 根据筛选条件，获取单页文件列表
func SelectFileList(c *gin.Context) {
	db := global.GlobalDB

	option := c.Query("option")     // 是按上传时间排序还是文件大小排序 data size
	category := c.Query("category") // 文件类型 image av other
	sortType := c.Query("sortType") // 排序方式 true-正序 false倒叙
	pageSize := c.Query("pageSize") // 每次查询数量
	start := c.Query("start")       // 开始查询位置

	if option == "" || category == "" || sortType == "" || pageSize == "" || start == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "缺少必要参数",
		})
		return
	}

	sortTypeBool, err := strconv.ParseBool(sortType)
	if err != nil {
		panic(err)
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		panic(err)
	}

	startInt, err := strconv.Atoi(start)
	if err != nil {
		panic(err)
	}

	var fileList []tables.SourceInfo
	var order string

	if sortTypeBool {
		order = option
	} else {
		order = fmt.Sprintf("%s desc", option)
	}

	result := db.Table("source_info").Where("category = ? AND temp = ?", category, false).Order(order).Limit(pageSizeInt).Offset(startInt).Find(&fileList)
	if result.Error != nil {
		panic(result.Error)
	}

	var count int64

	searchCount := db.Table("source_info").Where("category = ? AND temp = ?", category, false).Count(&count)
	if searchCount.Error != nil {
		panic(searchCount.Error)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data": FileListData{
			Count:    count,
			FileList: fileList,
		},
	})
}

// 删除文件(置为临时文件)
func SetFileTemp(c *gin.Context) {
	db := global.GlobalDB

	uid := c.PostForm("uid")

	if uid == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "缺少必要参数",
		})
		return
	}

	db.Model(&tables.SourceInfo{}).Where("uid = ?", uid).Update("temp", true)

	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data":    true,
	})
}

// 获取临时文件列表
func QueryTempFileList(c *gin.Context) {
	db := global.GlobalDB

	var fileList []tables.SourceInfo

	result := db.Table("source_info").Where("temp = ?", true).Find(&fileList)
	if result.Error != nil {
		panic(result.Error)
	}

	count := int64(len(fileList))

	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data": FileListData{
			Count:    count,
			FileList: fileList,
		},
	})
}

type RestitutionParam struct {
	UIDs []string `json:"uids"`
}

// 还原和批量还原(将临时文件恢复为正常文件)
func RestitutionFiles(c *gin.Context) {
	db := global.GlobalDB

	var uids RestitutionParam

	if err := c.ShouldBind(&uids); err != nil {
		panic(err)
	}

	for _, uid := range uids.UIDs {

		db.Model(&tables.SourceInfo{}).Where("uid = ?", uid).Update("temp", false)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data":    true,
	})
}

// 删除和批量删除(彻底删除临时文件)
type DeleteParam struct {
	Temps []tables.SourceInfo `json:"temps"`
}

func DeleteThorough(c *gin.Context) {
	db := global.GlobalDB

	var temps DeleteParam

	if err := c.ShouldBind(&temps); err != nil {
		panic(err)
	}

	for _, temp := range temps.Temps {
		// 从数据库中删除
		db.Where("uid = ?", temp.UID).Unscoped().Delete(&tables.SourceInfo{})

		// 从资源目录中删除
		path := fmt.Sprintf("%s/%s/%s", global.StaticPath, temp.Category, temp.Name)
		err := os.RemoveAll(path)
		if err != nil {
			panic(err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data":    true,
	})
}

// 修改文件信息--比如重命名文件、修改文件描述等
func UpdateFileInfo(c *gin.Context) {
	db := global.GlobalDB

	var fileparam tables.SourceInfo

	if err := c.ShouldBind(&fileparam); err != nil {
		panic(err)
	}

	// 先检查文件是否存在
	var fileQuery tables.SourceInfo
	result := db.Table("source_info").Where("uid = ?", fileparam.UID).Find(&fileQuery)
	if result.Error != nil {
		panic(result.Error)
	}

	if fileQuery.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeNotExist,
			"message": "文件不存在，修改失败",
		})
		return
	}

	// 检查文件名是否被修改，如果被修改，再检查是否有文件名重复
	// 多层判断，可以减少对数据库的操作
	// 如果修改了文件名且不重名，还需要修改静态目录中的文件名
	if fileparam.Name == fileQuery.Name && fileparam.Describe == fileQuery.Describe {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeNotModify,
			"message": "内容未变更",
		})
		return
	}

	if fileparam.Name == fileQuery.Name && fileparam.Describe != fileQuery.Describe {
		db.Model(&fileQuery).Update("describe", fileparam.Describe)

		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeOK,
			"message": "success",
			"data":    true,
		})
		return
	}

	var fileCheck tables.SourceInfo

	result = db.Table("source_info").Where("name = ?", fileparam.Name).Find(&fileCheck)
	if result.Error != nil {
		panic(result.Error)
	}

	if fileCheck.ID != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeExist,
			"message": "文件名已存在，修改失败",
		})
		return
	}

	prePath := fmt.Sprintf("%s/%s/", global.StaticPath, fileparam.Category)

	err := os.Rename(prePath+fileQuery.Name, prePath+fileparam.Name)
	if err != nil {
		panic(err)
	}

	db.Model(&fileQuery).Updates(tables.SourceInfo{Name: fileparam.Name, Describe: fileparam.Describe})

	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data":    true,
	})
}
