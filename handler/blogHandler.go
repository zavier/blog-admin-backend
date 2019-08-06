package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zavier/blog-admin-backend/common"
	"github.com/zavier/blog-admin-backend/server"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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

	if blog.Author == "" {
		blog.Author = context.GetString(common.JwtIdentityKey)
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

	var blogId = blog.Id
	if blogId <= 0 {
		context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, "ID不能小于等于0"))
		return
	}

	getBlog, err := server.GetBlog(blogId)
	if err != nil {
		context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, err.Error()))
		return
	}
	author := getBlog.Author
	if author != context.GetString(common.JwtIdentityKey) {
		context.JSON(http.StatusOK, ErrorResult(StatusUnauthorized, "您无权修改此博客"))
		return
	}

	if blog.Author == "" {
		blog.Author = context.GetString(common.JwtIdentityKey)
	}

	err = blog.UpdateBlog()
	if err != nil {
		log.Printf("update blog error: %s\n", err.Error())
		context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, err.Error()))
		return
	}
	context.JSON(http.StatusOK, SuccessResult(true))
}

// 删除博客
func DelBlog(context *gin.Context) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("DelBlog error")
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
			e := server.DelBlog(i)
			if e != nil {
				context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, e.Error()))
			} else {
				context.JSON(http.StatusOK, SuccessResult(true))
			}
		}
	}
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
				author := blog.Author
				if author != context.GetString(common.JwtIdentityKey) {
					context.JSON(http.StatusOK, ErrorResult(StatusUnauthorized, "无权访问此博客"))
				} else {
					blog.Location = ""
					context.JSON(http.StatusOK, SuccessResult(blog))
				}
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
		destBlogList := make([]server.BlogBase, 0)
		for _, b := range list {
			if b.Author == context.GetString(common.JwtIdentityKey) {
				b.Location = ""
				destBlogList = append(destBlogList, b)
			}
		}
		context.JSON(http.StatusOK, SuccessResult(destBlogList))
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

	err := server.HexoDeployAll(context.GetString(common.JwtIdentityKey))
	if err != nil {
		context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, err.Error()))
		return
	}
	context.JSON(http.StatusOK, SuccessResult(true))
}

// 上传博客文件
func Upload(context *gin.Context) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("upload error")
			context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, "系统错误"))
		}
	}()

	file, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, err.Error()))
		return
	}
	log.Printf("upload file name %s\n", file.Filename)

	if !strings.HasSuffix(file.Filename, ".md") {
		context.JSON(http.StatusOK, ErrorResult(StatusBadRequest, "文件格式错误"))
		return
	}

	src, err := file.Open()
	if err != nil {
		context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, err.Error()))
		return
	}
	defer func() {
		e := src.Close()
		if e != nil {
			log.Fatal("close file error")
		}
	}()

	bytes, err := ioutil.ReadAll(src)
	if err != nil {
		context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, err.Error()))
		return
	}

	blog := server.Blog{
		BlogBase: server.BlogBase{
			Title:  file.Filename[:len(file.Filename)-3],
			Author: context.GetString(common.JwtIdentityKey),
		},
		Content: string(bytes),
	}

	err = blog.SaveBlog()
	if err != nil {
		context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, err.Error()))
	} else {
		context.JSON(http.StatusOK, SuccessResult(true))
	}
}
