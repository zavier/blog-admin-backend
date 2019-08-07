package common

import (
	"log"
	"path/filepath"
)

const (
	// jwt秘钥
	JwtSecretKey   = "^BrX!!7V%gzH4#kF9e8Z"
	JwtIdentityKey = "name"

	// ==== 文件路径相关 ====
	// 博客目录
	blogPath = "blog"
	// 博客管理文件
	blogManageFileName = "blogList.json"
	// 密码文件
	pwdFilePath string = "pwd.txt"
	// 博客索引文件名
	indexFile = "index"
)

var (
	BlogPath           string
	BlogManageFileName string
	PwdFilePath        string
	IndexFile          string

	// 文档是否有更新, 暂时供自动发布使用
	BlogUpdated = false
)

func init() {
	curPath, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		log.Fatal("get curPath error", err)
	}
	BlogPath = curPath + "/" + blogPath
	BlogManageFileName = curPath + "/" + blogManageFileName
	PwdFilePath = curPath + "/" + pwdFilePath
	IndexFile = curPath + "/" + indexFile
}
