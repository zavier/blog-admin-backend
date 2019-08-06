package server

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/zavier/blog-admin-backend/constants"
	"github.com/zavier/blog-admin-backend/util"
	"log"
	"os"
	"strings"
)

type User struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 参数校验
func (user User) checkUserInfoParam() bool {
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

// 保存用户
func (user User) Save() (bool, error) {
	if !user.checkUserInfoParam() {
		return false, errors.New("参数错误，请检查用户名和密码长度")
	}
	if !util.Exists(constants.PwdFilePath) {
		_, err := os.Create(constants.PwdFilePath)
		if err != nil {
			return false, err
		}
	}
	hasRegistered, err := hasExistUserName(user.Name)
	if err != nil {
		return false, err
	}
	if hasRegistered {
		return false, errors.New("此用户名已存在")
	}

	file, err := os.OpenFile(constants.PwdFilePath, os.O_WRONLY|os.O_APPEND, 777)
	if err != nil {
		return false, err
	}
	defer func() {
		e := file.Close()
		if e != nil {
			log.Fatal("close file error")
		}
	}()

	hash := sha256.New()
	hash.Write([]byte(user.Password))
	shaPwd := hex.EncodeToString(hash.Sum(nil))
	_, err = file.WriteString(user.Name + ":" + shaPwd + "\n")
	if err != nil {
		return false, err
	}
	return true, nil
}

// 判断用户名称是否存在
func hasExistUserName(username string) (bool, error) {
	file, err := os.OpenFile(constants.PwdFilePath, os.O_RDONLY, 777)
	if err != nil {
		log.Printf("open pwd file error: %s\n", err.Error())
		return true, err
	}
	defer func() {
		e := file.Close()
		if e != nil {
			log.Fatal("close file error")
		}
	}()

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

// 登录(判断用户名和密码是否正确)
func (user User) CheckPassword() (correct bool, ex error) {
	file, err := os.OpenFile(constants.PwdFilePath, os.O_RDONLY, 777)
	if err != nil {
		return false, err
	}
	defer func() {
		e := file.Close()
		if e != nil {
			log.Fatal("close file error")
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		namePwdPair := scanner.Text()
		namePwd := strings.Split(namePwdPair, ":")
		if namePwd[0] == user.Name {
			hash := sha256.New()
			hash.Write([]byte(user.Password))
			shaPwd := hex.EncodeToString(hash.Sum(nil))
			if namePwd[1] == shaPwd {
				return true, nil
			} else {
				return false, errors.New("用户名或密码错误")
			}
		}
	}
	log.Printf("do not have this user:%s\n", user.Name)
	return false, errors.New("用户名或密码错误")
}
