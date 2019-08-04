package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zavier/blog-admin-backend/util"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const blogFilePath string = "blog"

type Blog struct {
	Title   string `form:"title" json:"title" binding:"required"`
	Context string `form:"context" json:"context" binding:"required"`
}

func init() {
	if !util.Exists(blogFilePath) {
		os.Mkdir(blogFilePath, 0777)
	}
}

func Save(context *gin.Context) {
	var blog Blog
	err := context.ShouldBind(&blog)
	if err != nil {
		log.Fatal("blog bind error", err)
		context.JSON(http.StatusOK, ErrorResult("参数错误"))
		return
	}

	//todo 参数长度相关校验

	filePath := blogFilePath + "/" + blog.Title + ".md"
	if !util.Exists(filePath) {
		log.Printf("filePath:%s does not exist, create.", filePath)
		_, err := os.Create(filePath)
		if err != nil {
			log.Fatalf("create file error", err)
			context.JSON(http.StatusInternalServerError, ErrorResult("系统错误"))
			return
		}
	}
	file, err := os.OpenFile(filePath, os.O_RDWR, 0777)
	if err != nil {
		log.Fatal("open file error", err)
		context.JSON(http.StatusInternalServerError, ErrorResult("系统错误"))
		return
	}
	file.WriteString(blog.Context)

	context.JSON(http.StatusOK, SuccessResult())
	return
}

func List(context *gin.Context) {
	infos, err := ioutil.ReadDir(blogFilePath)
	if err != nil {
		log.Fatal("read blog path error", err)
		context.JSON(http.StatusInternalServerError, ErrorResult("系统错误"))
		return
	}
	var blogs = make([]Blog, 0)
	for _, file := range infos {
		name := file.Name()
		blogs = append(blogs, Blog{
			Title: name,
		})
	}
	context.JSON(http.StatusOK, blogs)
}
