package api

import (
	"blog-api/api/activity/end"
	"blog-api/api/activity/page"
	"blog-api/api/activity/page/play_2399"
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
		v1.DELETE("/deleteOneLog", backstage.DeleteOneLog)

		v1.POST("/uploadDirect", backstage.UploadFileDirect)
		v1.POST("/preUpload", backstage.PreUpload)
		v1.POST("/chunkUpload", backstage.UploadChunk)
		v1.GET("/mergeChunks", backstage.MergeChunks)
		v1.GET("/selectFileList", backstage.SelectFileList)
		v1.PUT("/setFileTemp", backstage.SetFileTemp)
		v1.GET("/queryTempFileList", backstage.QueryTempFileList)
		v1.PUT("/restitutionFiles", backstage.RestitutionFiles)
		v1.PUT("/deleteThorough", backstage.DeleteThorough)
		v1.PUT("/updateFileInfo", backstage.UpdateFileInfo)
	}

	v2 := r.Group("/v2")
	{
		v2.POST("/createId", end.CreateId)
		v2.GET("/searchId", end.SearchId)
		v2.POST("/updateIdInfo", end.UpdateIdInfo)
		v2.POST("/addActivity", end.AddActivity)
		v2.GET("/searchActivityList", end.SearchActivityList)
		v2.GET("/searchActivityByBranch", end.SearchActivityByBranch)
		v2.POST("/removeActivity", end.RemoveActivity)
		v2.POST("/charge", end.Charge)
		v2.GET("/chargeList", end.GetChargeList)
		v2.POST("/chargeDel", end.ChargeDel)
		v2.GET("/searchGifts", end.SearchGifts)
		v2.POST("/addGift", end.AddGift)
		v2.POST("/deleteGifts", end.DelGifts)
		v2.POST("/updateGift", end.UpdateGift)
		v2.POST("/uploadGiftJsonFile", end.UploadGiftJsonFile)
	}

	v3 := r.Group("/v3")
	{
		v3.POST("/login", page.Login)
		v3.POST("/logout", page.Logout)
		v3.GET("/profileInfo", page.GetProfileInfo)
	}

	activity := r.Group("/activity")
	{
		activity.POST("/2399/sign", play_2399.Sign)
		activity.GET("/2399/signInfo", play_2399.SignInfo)
		activity.GET("/2399/cardInfo", play_2399.CardInfo)
		activity.POST("/2399/turnCard", play_2399.TurnCard)
	}
}
