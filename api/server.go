package api

import (
	"time"

	"blog-api/global"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Server() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Cookie", "ACTIVITY_SESSION_ID"},
		AllowCredentials: true,
		AllowOrigins:     []string{"http://localhost:3001", "http://121.40.42.63"},
		MaxAge:           24 * time.Hour, // 表示在24小时内，同样的预检请求可以不再重复进行了
	}))

	// 设置静态文件的路由和对应的文件目录。当用户访问/static中的文件时，就会去./static目录下查找资源文件。
	// 如果不加，则用户无法访问go项目中的资源。和koa2中的static方法类似。
	r.Static("/"+global.StaticPath, "./"+global.StaticPath)

	groupRouter(r)

	r.Run(server)
}
