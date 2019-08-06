package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zavier/blog-admin-backend/constants"
	"github.com/zavier/blog-admin-backend/server"
	"log"
	"net/http"
	"strconv"
)

// 博客保存
func Save(context *gin.Context) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("save error")
			context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, "系统错误"))
		}
	}()

	var blog server.Blog
	err := context.ShouldBind(&blog)
	if err != nil {
		log.Printf("blog bind error: %s\n", err.Error())
		context.JSON(http.StatusOK, ErrorResult(StatusBadRequest, "参数错误"))
		return
	}

	err = blog.SaveBlog()
	if err != nil {
		log.Printf("save blog error: %s\n", err.Error())
		context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, err.Error()))
		return
	}
	context.JSON(http.StatusOK, SuccessResult(true))
}

// 更新博客
func Update(context *gin.Context) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("save error")
			context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, "系统错误"))
		}
	}()

	var blog server.Blog
	err := context.ShouldBind(&blog)
	if err != nil {
		log.Printf("blog bind error: %s\n", err.Error())
		context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, "参数错误"))
		return
	}
	err = blog.UpdateBlog()
	if err != nil {
		log.Printf("update blog error: %s\n", err.Error())
		context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, err.Error()))
		return
	}
	context.JSON(http.StatusOK, SuccessResult(true))
}

// 查询博客内容
func GetBlog(context *gin.Context) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("getBlog error")
			context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, "系统错误"))
		}
	}()

	id := context.Query("id")
	if id == "" {
		context.JSON(http.StatusOK, ErrorResult(StatusBadRequest, "ID不能为空"))
	} else {
		i, e := strconv.Atoi(id)
		if e != nil {
			context.JSON(http.StatusOK, ErrorResult(StatusBadRequest, "参数错误"))
		} else {
			blog, err := server.GetBlog(i)
			if err != nil {
				context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, err.Error()))
			} else {
				context.JSON(http.StatusOK, SuccessResult(blog))
			}
		}
	}
}

// 查询博客列表
func List(context *gin.Context) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("listBlog error")
			context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, "系统错误"))
		}
	}()

	list, err := server.BlogList()
	if err != nil {
		context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, err.Error()))
	} else {
		context.JSON(http.StatusOK, SuccessResult(list))
	}
}

// 发布全部博客
func DeployAll(context *gin.Context) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("deployAll Blog error")
			context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, "系统错误"))
		}
	}()

	err := server.HexoDeployAll()
	if err != nil {
		context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, err.Error()))
	}
	context.JSON(http.StatusOK, SuccessResult(true))
}

// todo 上传博客文件
func Upload(context *gin.Context) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("upload error")
			context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, "系统错误"))
		}
	}()

	file, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, "上传文件失败"))
		return
	}
	log.Printf("upload file name %s\n", file.Filename)

	err = context.SaveUploadedFile(file, constants.BlogPath+"/"+file.Filename)
	if err != nil {
		context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, "上传文件失败"))
		return
	}

	context.JSON(http.StatusOK, SuccessResult(true))
}
