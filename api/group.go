package api

import (
	"blog-api/api/backstage"

	"github.com/gin-gonic/gin"
)

func groupRouter(r *gin.Engine) {

	v1 := r.Group("/v1")
	{
		v1.POST("/publishLogs", backstage.PublishLogs)
		v1.GET("/getDateLogs", backstage.FindDateLogs)
	}
}
