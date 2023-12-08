package end

import (
	"net/http"
	"time"

	"blog-api/db_server/tables"
	"blog-api/global"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func AddActivity(c *gin.Context) {
	db := global.GlobalDB

	branch := c.PostForm("branch")
	name := c.PostForm("name")
	tag := c.PostForm("tag")
	url := c.PostForm("url")
	dateStart := c.PostForm("dateStart")
	dateEnd := c.PostForm("dateEnd")
	moudleStart := c.PostForm("moudleStart")
	moudleEnd := c.PostForm("moudleEnd")
	if branch == "" || name == "" || tag == "" || url == "" || dateStart == "" || dateEnd == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "缺少必要参数",
			"data":    nil,
		})
		return
	}

	dst, err := time.Parse("2006-01-02 15:04:05", dateStart)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeFormatError,
			"message": "参数格式错误",
			"data":    nil,
		})
		return
	}
	det, err := time.Parse("2006-01-02 15:04:05", dateEnd)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeFormatError,
			"message": "参数格式错误",
			"data":    nil,
		})
		return
	}
	var mst time.Time
	if moudleStart != "" {
		_mst, err := time.Parse("2006-01-02 15:04:05", moudleStart)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeFormatError,
				"message": "参数格式错误",
				"data":    nil,
			})
			return
		}
		mst = _mst
	}
	var met time.Time
	if moudleEnd != "" {
		_met, err := time.Parse("2006-01-02 15:04:05", moudleEnd)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeFormatError,
				"message": "参数格式错误",
				"data":    nil,
			})
			return
		}
		met = _met
	}
	info := tables.ActivityListInfo{
		Branch:      branch,
		Name:        name,
		Tag:         tag,
		Url:         url,
		DateStart:   dst.UnixMilli(),
		DateEnd:     det.UnixMilli(),
		MoudleStart: mst.UnixMilli(),
		MoudleEnd:   met.UnixMilli(),
	}
	result := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&info)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeCreateDataFailed,
			"message": "添加失败",
			"data":    nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data":    true,
	})
}
