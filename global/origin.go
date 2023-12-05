package global

import (
	"github.com/gin-gonic/gin"
)

func GlobalOrigin() string {
	var origin string

	if gin.Mode() == "release" {
		origin = "https://multi-app-blog.fun:9000" // 线上源
	} else {
		origin = "http://127.0.0.1:4001" // 开发本地源
	}

	return origin
}
