package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zavier/blog-admin-backend/common"
	"github.com/zavier/blog-admin-backend/server"
	"log"
	"net/http"
)

// 注册用户
func Register(context *gin.Context) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("register error")
			context.JSON(http.StatusOK, common.ErrorResult(common.StatusInternalServerError, "系统错误"))
		}
	}()

	var user server.User
	err := context.ShouldBind(&user)
	if err != nil {
		log.Printf("bind user param error: %s\n", err.Error())
		context.JSON(http.StatusOK, common.ErrorResult(common.StatusBadRequest, "参数错误"))
		return
	}

	b, err := user.Save()
	if err != nil {
		context.JSON(http.StatusOK, common.ErrorResult(common.StatusInternalServerError, err.Error()))
		return
	}
	if !b {
		context.JSON(http.StatusOK, common.ErrorResult(common.StatusInternalServerError, "系统错误"))
		return
	}
	context.JSON(http.StatusOK, common.SuccessResult(true))
}
