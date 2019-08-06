package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zavier/blog-admin-backend/handler"
	"github.com/zavier/blog-admin-backend/middleware"
	"log"
	"net/http"
)

func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)
}

func main() {
	router := gin.Default()

	// 登陆相关
	userRouter := router.Group("/user")
	userRouter.POST("/register", handler.Register)
	userRouter.POST("/login", middleware.JwtMiddleware().LoginHandler)

	// 下面的功能开始需要进行认证了
	router.Use(middleware.JwtMiddleware().MiddlewareFunc())
	// 校验是否登录
	router.GET("/check/isLogin", func(context *gin.Context) {
		context.JSON(http.StatusOK, handler.SuccessResult(true))
	})

	// 博客相关
	blogRouter := router.Group("/blog")
	blogRouter.POST("/save", handler.Save)
	blogRouter.POST("/update", handler.Update)
	blogRouter.GET("/list", handler.List)
	blogRouter.GET("/getBlog", handler.GetBlog)
	//blogRouter.POST("/upload", handler.Upload)
	blogRouter.POST("/deployAll", handler.DeployAll)

	err := router.Run(":8081")
	if err != nil {
		log.Fatal("server start error", err)
	}
}
