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
		uid = 0
	}
	pid, err := strconv.Atoi(payId)
	if err != nil {
		pid = 0
	}
	start, err := strconv.Atoi(dateStart)
	if err != nil {
		start = 0
	}
	end, err := strconv.Atoi(dateEnd)
	if err != nil {
		end = 0
	}
	p, err := strconv.Atoi(page)
	if err != nil {
		p = 0
	}
	ps, err := strconv.Atoi(pageSize)
	if err != nil {
		ps = 0
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
	var chargeList []tables.ChargeInfo
	var total int64
	if uid == 0 && pid == 0 && start == 0 && end == 0 { /* search all */
		db.Table(tb.TableName()).Order("id desc").Limit(ps).Offset((p - 1) * ps).Find(&chargeList)
		db.Table(tb.TableName()).Count(&total)
	} else if uid != 0 && pid == 0 && start == 0 && end == 0 { /* search by uid */
		db.Table(tb.TableName()).Where("userId = ?", uid).Order("id desc").Limit(ps).Offset((p - 1) * ps).Find(&chargeList)
		db.Model(&tables.ChargeInfo{}).Where("userId = ?", uid).Count(&total)
	} else if uid == 0 && pid != 0 && start == 0 && end == 0 { /* search by pid */
		db.Table(tb.TableName()).Where("payId = ?", pid).Order("id desc").Limit(ps).Offset((p - 1) * ps).Find(&chargeList)
		db.Model(&tables.ChargeInfo{}).Where("payId = ?", pid).Count(&total)
	} else if uid == 0 && pid == 0 && ((start != 0 && end == 0) || (start == 0 && end != 0) || (start != 0 && end != 0)) { /* search by start and end */
		db.Table(tb.TableName()).Where("date >= ? AND date <= ?", start, end).Order("id desc").Limit(ps).Offset((p - 1) * ps).Find(&chargeList)
		db.Model(&tables.ChargeInfo{}).Where("date >= ? AND date <= ?", start, end).Count(&total)
	} else if uid != 0 && pid != 0 && start == 0 && end == 0 { /* search by uid and pid */
		db.Table(tb.TableName()).Where("userId = ? AND payId = ?", uid, pid).Order("id desc").Limit(ps).Offset((p - 1) * ps).Find(&chargeList)
		db.Model(&tables.ChargeInfo{}).Where("userId = ? AND payId = ?", uid, pid).Count(&total)
	} else { /* search by uid and pid and start and end */
		db.Table(tb.TableName()).Where("userId = ? AND payId = ? AND date >= ? AND date <= ?", uid, pid, start, end).Order("id desc").Limit(ps).Offset((p - 1) * ps).Find(&chargeList)
		db.Model(&tables.ChargeInfo{}).Where("userId = ? AND payId = ? AND date >= ? AND date <= ?", uid, pid, start, end).Count(&total)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    global.CodeOK,
		"message": "success",
		"data": map[string]interface{}{
			"chargeList": chargeList,
			"total":      total,
		},
	})
}
