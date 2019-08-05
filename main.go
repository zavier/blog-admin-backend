package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zavier/blog-admin-backend/handler"
	"github.com/zavier/blog-admin-backend/middleware"
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
	userRouter.POST("/login", middleware.JwtMiddleware().LoginHandler)
	// todo 可能不需要退出了，让前端自己删除token即可，是否登陆也可以让前段自行判断
	userRouter.POST("/logout", handler.Logout)
	userRouter.GET("/isLogin", handler.IsLogIn)

	// 下面的功能开始需要进行认证了
	router.Use(middleware.JwtMiddleware().MiddlewareFunc())

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
