package handler

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"github.com/zavier/blog-admin-backend/constants"
	"github.com/zavier/blog-admin-backend/util"
	"log"
	"net/http"
	"os"
	"strings"
)

const pwdFilePath string = "pwd.txt"

type User struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 注册
func Register(context *gin.Context) {
	var user User
	err := context.ShouldBind(&user)
	if err != nil {
		log.Fatal("bind user param error", err)
		context.JSON(http.StatusBadRequest, ErrorResult("参数错误"))
		return
	}

	//todo 用户信息相关校验、长度、特殊符号等
	name := user.Name
	if strings.Contains(name, ":") {
		context.JSON(http.StatusBadRequest, ErrorResult("名称中不能包含特殊符号"))
		return
	}

	if !util.Exists(pwdFilePath) {
		_, err := os.Create(pwdFilePath)
		if err != nil {
			log.Fatal("create pwd file error", err)
			context.JSON(http.StatusInternalServerError, ErrorResult("系统错误"))
			return
		}
	}

	file, err := os.OpenFile(pwdFilePath, os.O_APPEND, 777)
	if err != nil {
		log.Fatal("open pwd file error", err)
		context.JSON(http.StatusInternalServerError, ErrorResult("系统错误"))
		return
	}
	defer file.Close()

	//hash := sha256.New()
	//hash.Write([]byte(user.Password))
	//shaPwd := string(hash.Sum(nil))
	file.WriteString(user.Name + ":" + user.Password + "\n")

	context.JSON(http.StatusOK, SuccessResult())
}

// 登陆
func Login(context *gin.Context) {
	var user User
	err := context.ShouldBind(&user)
	if err != nil {
		log.Fatal("bind user param error", err)
		context.JSON(http.StatusBadRequest, ErrorResult("参数错误"))
		return
	}

	file, err := os.OpenFile(pwdFilePath, os.O_RDONLY, 777)
	if err != nil {
		log.Fatal("open pwd file error", err)
		context.JSON(http.StatusInternalServerError, ErrorResult("系统错误"))
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		namePwdPair := scanner.Text()
		log.Printf("read pwd context:%s\n", namePwdPair)
		namePwd := strings.Split(namePwdPair, ":")
		if namePwd[0] == user.Name {
			//hash := sha256.New()
			//hash.Write([]byte(namePwd[1]))
			//shaPwd := string(hash.Sum(nil))
			if namePwd[1] == user.Password {
				context.SetCookie(constants.TokenName, namePwd[1], 65535, "/", "localhost", false, true)
				context.JSON(http.StatusOK, SuccessResult())
				return
			} else {
				log.Fatal("password error, excected:%s actual:%s\n", namePwd[1], user.Password)
				context.JSON(http.StatusOK, ErrorResult("用户名或密码错误"))
				return
			}
		}
	}
	log.Fatal("do not have this user:%s\n", user.Name)
	context.JSON(http.StatusOK, ErrorResult("用户名或密码错误"))
}

// 登出
func Logout(context *gin.Context) {
	context.SetCookie(constants.TokenName, "", 65535, "/", "localhost", false, true)
	context.JSON(http.StatusOK, SuccessResult())
}
