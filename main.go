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
	userRouter := router.Group("/user")
	userRouter.POST("/register", handler.Register)
	userRouter.POST("/login", handler.Login)
	userRouter.POST("/logout", handler.Logout)
	userRouter.GET("/isLogin", handler.IsLogIn)

	// 博客相关
	blogRouter := router.Group("/blog")
	blogRouter.POST("/save", handler.Save)
	blogRouter.POST("/update", handler.Update)
	blogRouter.GET("/list", handler.List)
	blogRouter.GET("/getBlog", handler.GetBlog)
	//blogRouter.POST("/upload", handler.Upload)
	blogRouter.POST("/deployAll", handler.DeployAll)

	router.Run(":8081")
}
