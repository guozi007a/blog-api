package test

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Testapi() {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.GET("/about", About)
	}

	v2 := r.Group("/v2")
	{
		v2.POST("/list", List)
	}

	r.Run(":4001")
}

func About(c *gin.Context) {
	data := map[string]interface{}{
		"name":   "dilireba",
		"age":    18,
		"gender": "女",
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    data,
	})
}

func List(c *gin.Context) {
	data := map[string]interface{}{
		"name":   "dilireba",
		"age":    18,
		"gender": "女",
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    data,
	})
}
