package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Server() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Cookie"},
		AllowCredentials: true,
		AllowOrigins:     []string{"*"},
		MaxAge:           24 * time.Hour, // 表示在24小时内，同样的预检请求可以不再重复进行了
	}))

	groupRouter(r)

	r.Run(server)
}
