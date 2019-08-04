package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zavier/blog-admin-backend/constants"
	"github.com/zavier/blog-admin-backend/server"
	"log"
	"net/http"
)

// 博客保存
func Save(context *gin.Context) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("save error")
			context.JSON(http.StatusOK, ErrorResult("系统错误"))
		}
	}()

	var blog server.Blog
	err := context.ShouldBind(&blog)
	if err != nil {
		log.Printf("blog bind error: %s\n", err.Error())
		context.JSON(http.StatusOK, ErrorResult("参数错误"))
		return
	}

	err = blog.SaveBlog()
	if err != nil {
		log.Printf("save blog error: %s\n", err.Error())
		context.JSON(http.StatusInternalServerError, ErrorResult("系统错误"))
		return
	}
	context.JSON(http.StatusOK, SuccessResult())
}

// 更新博客
func Update(context *gin.Context) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("save error")
			context.JSON(http.StatusOK, ErrorResult("系统错误"))
		}
	}()

	var blog server.Blog
	err := context.ShouldBind(&blog)
	if err != nil {
		log.Printf("blog bind error: %s\n", err.Error())
		context.JSON(http.StatusOK, ErrorResult("参数错误"))
		return
	}
	err = blog.UpdateBlog()
	if err != nil {
		log.Printf("update blog error: %s\n", err.Error())
		context.JSON(http.StatusInternalServerError, ErrorResult("系统错误"))
		return
	}
	context.JSON(http.StatusOK, SuccessResult())
}

// 查询博客列表
func List(context *gin.Context) {
	list := server.BlogList()
	context.JSON(http.StatusOK, list)
}

// todo 上传博客文件
func Upload(context *gin.Context) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("save error")
			context.JSON(http.StatusOK, ErrorResult("系统错误"))
		}
	}()

	file, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusInternalServerError, ErrorResult("上传文件失败"))
		return
	}
	log.Printf("upload file name %s\n", file.Filename)

	err = context.SaveUploadedFile(file, constants.BlogPath+"/"+file.Filename)
	if err != nil {
		context.JSON(http.StatusInternalServerError, ErrorResult("上传文件失败"))
		return
	}

	context.JSON(http.StatusOK, SuccessResult())
}
