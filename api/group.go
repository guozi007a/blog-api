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
		v1.GET("/getAllLogs", backstage.FindAllLogs)
		v1.PUT("/updateDateLogs", backstage.UpdateDateLogs)
		v1.DELETE("/deleteDateLogs", backstage.DeleteDateLogs)
		v1.DELETE("/clearAllLogs", backstage.ClearAllLogs)

		v1.POST("/uploadDirect", backstage.UploadFileDirect)
		v1.POST("/preUpload", backstage.PreUpload)
		v1.POST("/chunkUpload", backstage.UploadChunk)
		v1.GET("/mergeChunks", backstage.MergeChunks)
		v1.GET("/selectFileList", backstage.SelectFileList)
		v1.PUT("/setFileTemp", backstage.SetFileTemp)
		v1.GET("/queryTempFileList", backstage.QueryTempFileList)
		v1.PUT("/restitutionFiles", backstage.RestitutionFiles)
		v1.DELETE("/deleteThorough", backstage.DeleteThorough)
	}
}
