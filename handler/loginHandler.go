package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zavier/blog-admin-backend/constants"
	"github.com/zavier/blog-admin-backend/server"
	"log"
	"net/http"
)

// 注册
func Register(context *gin.Context) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("register error")
			context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, "系统错误"))
		}
	}()

	var user server.User
	err := context.ShouldBind(&user)
	if err != nil {
		log.Printf("bind user param error: %s\n", err.Error())
		context.JSON(http.StatusOK, ErrorResult(StatusBadRequest, "参数错误"))
		return
	}

	b, err := server.Register(user)
	if err != nil || !b {
		context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, "系统错误"))
		return
	}
	context.JSON(http.StatusOK, SuccessResult(true))
}

// 登陆
func Login(context *gin.Context) {
	var user server.User
	err := context.ShouldBind(&user)
	if err != nil {
		log.Printf("bind user param error: %s\n", err.Error())
		context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, "参数错误"))
		return
	}

	token, err := server.Login(user)
	if err != nil {
		context.JSON(http.StatusOK, ErrorResult(StatusInternalServerError, err.Error()))
		return
	} else {
		context.SetCookie(constants.TokenName, token, constants.CookieMaxAge, "/", constants.Domain, false, true)
		context.JSON(http.StatusOK, SuccessResult(true))
	}
}

// 是否已登陆
func IsLogIn(context *gin.Context) {
	token, err := context.Request.Cookie(constants.TokenName)
	if err != nil || token == nil || token.Value == "" {
		context.JSON(http.StatusOK, SuccessResult(false))
	} else {
		log.Printf("token:%v token name:%s, token value:%s", token, token.Name, token.Value)
		context.JSON(http.StatusOK, SuccessResult(true))
	}
}

// 登出
func Logout(context *gin.Context) {
	context.SetCookie(constants.TokenName, "", 1, "/", constants.Domain, false, true)
	context.JSON(http.StatusOK, SuccessResult(true))
}
