package api

import (
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
		MaxAge:           20,
	}))

	groupRouter(r)

	r.Run(server)
}
