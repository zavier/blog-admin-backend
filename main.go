package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zavier/blog-admin-backend/handler"
	"log"
)

func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)
}

func main() {
	router := gin.Default()
	// 认证校验, CORS
	//router.Use(middleware.Auth(), middleware.Cors())

	// 登陆相关
	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)
	router.POST("/logout", handler.Logout)

	// 博客相关
	blogRouter := router.Group("/blog")
	blogRouter.POST("/save", handler.Save)
	blogRouter.POST("/update", handler.Update)
	blogRouter.GET("/list", handler.List)
	//blogRouter.POST("/upload", handler.Upload)

	router.Run(":8081")
}
