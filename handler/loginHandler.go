package handler

import (
	"bufio"
	"crypto/sha256"
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

// 校验，todo 后续可考虑通过binding实现
func (user User) checkUserInfo() bool {
	var nameLen = len(user.Name)
	if nameLen < 2 || nameLen > 8 {
		log.Printf("username len:%d is invalid\n", nameLen)
		return false
	}
	if strings.Contains(user.Name, ":") {
		log.Printf("username len:%d is invalid\n", nameLen)
		return false
	}

	pwdLen := len(user.Password)
	if pwdLen < 5 || pwdLen > 12 {
		log.Printf("userpwd len:%d is invalid", pwdLen)
		return false
	}
	return true
}

// 注册
func Register(context *gin.Context) {
	var user User
	err := context.ShouldBind(&user)
	if err != nil {
		log.Printf("bind user param error: %s\n", err.Error())
		context.JSON(http.StatusBadRequest, ErrorResult("参数错误"))
		return
	}

	if !user.checkUserInfo() {
		context.JSON(http.StatusBadRequest, ErrorResult("参数错误，请检查用户名和密码长度"))
		return
	}

	if !util.Exists(pwdFilePath) {
		_, err := os.Create(pwdFilePath)
		if err != nil {
			log.Printf("create pwd file error: %s\n", err.Error())
			context.JSON(http.StatusInternalServerError, ErrorResult("系统错误"))
			return
		}
	}

	file, err := os.OpenFile(pwdFilePath, os.O_WRONLY|os.O_APPEND, 777)
	if err != nil {
		log.Printf("open pwd file error: %s\n", err.Error())
		context.JSON(http.StatusInternalServerError, ErrorResult("系统错误"))
		return
	}
	defer file.Close()

	hash := sha256.New()
	hash.Write([]byte(user.Password))
	shaPwd := string(hash.Sum(nil))
	file.WriteString(user.Name + ":" + shaPwd + "\n")

	context.JSON(http.StatusOK, SuccessResult())
}

func hasRegister(username string) (bool, error) {
	file, err := os.OpenFile(pwdFilePath, os.O_RDONLY, 777)
	if err != nil {
		log.Printf("open pwd file error: %s\n", err.Error())
		return true, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		namePwdPair := scanner.Text()
		log.Printf("read pwd context:%s\n", namePwdPair)
		namePwd := strings.Split(namePwdPair, ":")
		if namePwd[0] == username {
			return true, nil
		}
	}
	return false, nil
}

// 登陆
func Login(context *gin.Context) {
	var user User
	err := context.ShouldBind(&user)
	if err != nil {
		log.Printf("bind user param error: %s\n", err.Error())
		context.JSON(http.StatusBadRequest, ErrorResult("参数错误"))
		return
	}

	file, err := os.OpenFile(pwdFilePath, os.O_RDONLY, 777)
	if err != nil {
		log.Printf("open pwd file error: %s\n", err.Error())
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
			hash := sha256.New()
			hash.Write([]byte(namePwd[1]))
			shaPwd := string(hash.Sum(nil))
			if shaPwd == user.Password {
				context.SetCookie(constants.TokenName, namePwd[1], 65535, "/", "localhost", false, true)
				context.JSON(http.StatusOK, SuccessResult())
				return
			} else {
				log.Printf("password error, excected:%s actual:%s\n", namePwd[1], user.Password)
				context.JSON(http.StatusOK, ErrorResult("用户名或密码错误"))
				return
			}
		}
	}
	log.Printf("do not have this user:%s\n", user.Name)
	context.JSON(http.StatusOK, ErrorResult("用户名或密码错误"))
}

// 登出
func Logout(context *gin.Context) {
	context.SetCookie(constants.TokenName, "", 65535, "/", "localhost", false, true)
	context.JSON(http.StatusOK, SuccessResult())
}
