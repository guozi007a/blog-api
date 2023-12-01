package end

import (
	"net/http"

	"blog-api/db_server/tables"
	"blog-api/global"

	"strconv"

	"github.com/gin-gonic/gin"
)

type ChargeListItem struct {
	UserId    int   `json:"userId"`
	PayId     int   `json:"payId"`
	DateStart int64 `json:"dateStart"`
	DateEnd   int64 `json:"dateEnd"`
	PageSize  int   `json:"pageSize"`
	Page      int   `json:"page"`
}

func GetChargeList(c *gin.Context) {
	db := global.GlobalDB
	userId := c.Query("userId")
	payId := c.Query("payId")
	dateStart := c.Query("dateStart")
	dateEnd := c.Query("dateEnd")
	pageSize := c.Query("pageSize")
	page := c.Query("page")
	uid, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeInvalid,
			"message": "参数无效",
			"data":    nil,
		})
		panic(err)
	}
	pid, err := strconv.Atoi(payId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeInvalid,
			"message": "参数无效",
			"data":    nil,
		})
		panic(err)
	}
	start, err := strconv.Atoi(dateStart)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeInvalid,
			"message": "参数无效",
			"data":    nil,
		})
		panic(err)
	}
	end, err := strconv.Atoi(dateEnd)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeInvalid,
			"message": "参数无效",
			"data":    nil,
		})
		panic(err)
	}
	p, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeInvalid,
			"message": "参数无效",
			"data":    nil,
		})
		panic(err)
	}
	ps, err := strconv.Atoi(pageSize)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeInvalid,
			"message": "参数无效",
			"data":    nil,
		})
		panic(err)
	}
	if p == 0 || ps == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    global.CodeLackRequired,
			"message": "缺少必要参数",
			"data":    nil,
		})
		return
	}
	var tb tables.ChargeInfo
	var chargeList []ChargeListItem
	if uid == 0 && pid == 0 && start == 0 && end == 0 { /* search all */
		result := db.Table(tb.TableName()).Order("id desc").Limit(ps).Offset((p - 1) * ps).Find(&chargeList)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			panic(result.Error)
		}
	} else if uid != 0 && pid == 0 && start == 0 && end == 0 { /* search by uid */
		result := db.Table(tb.TableName()).Where("userId = ?", uid).Order("id desc").Limit(ps).Offset((p - 1) * ps).Find(&chargeList)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			panic(result.Error)
		}
	} else if uid == 0 && pid != 0 && start == 0 && end == 0 { /* search by pid */
		result := db.Table(tb.TableName()).Where("payId = ?", pid).Order("id desc").Limit(ps).Offset((p - 1) * ps).Find(&chargeList)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			panic(result.Error)
		}
	} else if uid == 0 && pid == 0 && ((start != 0 && end == 0) || (start == 0 && end != 0) || (start != 0 && end != 0)) { /* search by start and end */
		result := db.Table(tb.TableName()).Where("date >= ? AND date <= ?", start, end).Order("id desc").Limit(ps).Offset((p - 1) * ps).Find(&chargeList)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			panic(result.Error)
		}
	} else if uid != 0 && pid != 0 && start == 0 && end == 0 { /* search by uid and pid */
		result := db.Table(tb.TableName()).Where("userId = ? AND payId = ?", uid, pid).Order("id desc").Limit(ps).Offset((p - 1) * ps).Find(&chargeList)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			panic(result.Error)
		}
	} else { /* search by uid and pid and start and end */
		result := db.Table(tb.TableName()).Where("userId = ? AND payId = ? AND date >= ? AND date <= ?", uid, pid, start, end).Order("id desc").Limit(ps).Offset((p - 1) * ps).Find(&chargeList)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    global.CodeQueryFailed,
				"message": "查询失败",
				"data":    nil,
			})
			panic(result.Error)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data":    chargeList,
	})
}
