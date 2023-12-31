package end

import (
	"net/http"

	"blog-api/db_server/tables"
	"blog-api/global"
	"blog-api/utils"

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

	dst := utils.LocMilli(dateStart)
	if dst == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeFormatError,
			"message": "参数格式错误",
			"data":    nil,
		})
		return
	}

	det := utils.LocMilli(dateEnd)
	if det == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeFormatError,
			"message": "参数格式错误",
			"data":    nil,
		})
		return
	}
	var mst int64 = 0
	if moudleStart != "" {
		_mst := utils.LocMilli(moudleStart)
		if _mst == 0 {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeFormatError,
				"message": "参数格式错误",
				"data":    nil,
			})
			return
		}
		mst = int64(_mst)
	}
	var met int64 = 0
	if moudleEnd != "" {
		_met := utils.LocMilli(moudleEnd)
		if _met == 0 {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeFormatError,
				"message": "参数格式错误",
				"data":    nil,
			})
			return
		}
		met = int64(_met)
	}
	info := tables.ActivityListInfo{
		Branch:      branch,
		Name:        name,
		Tag:         tag,
		Url:         url,
		DateStart:   dst,
		DateEnd:     det,
		MoudleStart: mst,
		MoudleEnd:   met,
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
