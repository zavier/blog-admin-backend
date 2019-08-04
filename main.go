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

type Blog struct {
	Title   string `form:"title" json:"title" binding:"required"`
	Context string `form:"context" json:"context" binding:"required"`
}

type Result struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func main() {
	router := gin.Default()
	// 认证校验, CORS
	router.Use(middleware.Auth(), middleware.Cors())

	// 登陆相关
	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)
	router.POST("/logout", handler.Logout)

	// 博客相关
	blogRouter := router.Group("/blog")
	blogRouter.POST("/save", handler.Save)
	blogRouter.GET("/list", handler.List)

	router.POST("/upload", func(c *gin.Context) {
		file, err1 := c.FormFile("file")
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, Result{
				Success: false,
				Data:    nil,
				Message: "上传文件失败",
			})
			return
		}
		log.Println(file.Filename)

		err := c.SaveUploadedFile(file, "/Users/zhengwei/go/src/github.com/zavier/blog-admin-backend/"+file.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Result{
				Success: false,
				Data:    nil,
				Message: "上传文件失败",
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	})
	router.Run()
}
