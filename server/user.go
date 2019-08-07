package server

import (
	"bufio"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"github.com/zavier/blog-admin-backend/common"
	"log"
	"os"
	"strings"
)

type User struct {
	Name     string `json:"name" binding:"required,min=2,max=8"`
	Password string `json:"password" binding:"required,min=5,max=12"`
}

func init() {
	exists, e := common.ExistsPath(common.PwdFilePath)
	if e != nil {
		log.Fatal("exist pwd file error", e)
	}
	if !exists {
		if _, err := os.Create(common.PwdFilePath); err != nil {
			log.Fatal("create pwd file error", err)
		}
	}
}

// 保存用户
func (user *User) Save() (bool, error) {
	hasRegistered, err := hasExistUserName(user.Name)
	if err != nil {
		return false, err
	}
	if hasRegistered {
		return false, errors.New("此用户名已存在")
	}

	file, err := os.OpenFile(common.PwdFilePath, os.O_WRONLY|os.O_APPEND, 777)
	if err != nil {
		return false, err
	}
	defer func() {
		e := file.Close()
		if e != nil {
			log.Fatal("close file error")
		}
	}()

	base64Name := base64.StdEncoding.EncodeToString([]byte(user.Name))

	hash := sha256.New()
	hash.Write([]byte(user.Password))
	shaPwd := hex.EncodeToString(hash.Sum(nil))

	if _, err = file.WriteString(base64Name + ":" + shaPwd + "\n"); err != nil {
		return false, err
	}
	return true, nil
}

// 判断用户名称是否存在
func hasExistUserName(username string) (bool, error) {
	file, err := os.OpenFile(common.PwdFilePath, os.O_RDONLY, 777)
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
		namePwd := strings.Split(namePwdPair, ":")
		if namePwd[0] == base64.StdEncoding.EncodeToString([]byte(username)) {
			return true, nil
		}
	}
	return false, nil
}

// 登录(判断用户名和密码是否正确)
func (user *User) CheckPassword() (correct bool, ex error) {
	file, err := os.OpenFile(common.PwdFilePath, os.O_RDONLY, 777)
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
		if namePwd[0] == base64.StdEncoding.EncodeToString([]byte(user.Name)) {
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
