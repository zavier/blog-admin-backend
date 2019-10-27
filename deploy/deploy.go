package deploy

import (
	"github.com/zavier/blog-admin-backend/common"
	"github.com/zavier/blog-admin-backend/server"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var hexoHeader = `---
title: ${title}
date: ${date}
tags: ${tags}
---
`

// 暂停自动发布功能（已将方法名从init改为initStop）
func initStop() {
	updatedUser := "zheng"
	// 启动定时任务，定时扫描状态，进行发布
	common.StartTimer(func() {
		if !common.BlogUpdated {
			return
		}
		log.Println("blog has update, start deploying...")
		common.BlogUpdated = false
		err := HexoDeployAll(updatedUser)
		if err != nil {
			log.Printf("deploy error: %s\n", err.Error())
		}
	}, 5*time.Minute)
}

// 全部发布(会清空之前内容)
func HexoDeployAll(username string) error {
	sourcePath := os.Getenv("SOURCE_PATH")
	log.Printf("hexo source path: %s\n", sourcePath)
	err := os.RemoveAll(sourcePath)
	if err != nil {
		return err
	}
	err = os.MkdirAll(sourcePath, 777)
	if err != nil {
		return err
	}

	// 遍历所有博客，拷贝文件
	baseBlogList, err := server.BlogList()
	if err != nil {
		return err
	}
	for _, blog := range baseBlogList {
		if blog.Author != username {
			continue
		}
		bytes, err := ioutil.ReadFile(blog.Location)
		if err != nil {
			return err
		}
		newHeader := strings.ReplaceAll(hexoHeader, "${title}", blog.Title)
		newHeader = strings.ReplaceAll(newHeader, "${date}", time.Now().Format("2006-01-02 15:04:05"))
		s := blog.Categories
		if s != "" {
			newHeader = strings.ReplaceAll(newHeader, "${tags}", "["+blog.Categories+"]")
		} else {
			newHeader = strings.ReplaceAll(newHeader, "${tags}", "")
		}
		var context = newHeader + string(bytes)
		newFile, err := os.OpenFile(sourcePath+"/"+blog.Title+".md", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
		if err != nil {
			return err
		}
		_, err = newFile.WriteString(context)
		if err != nil {
			return err
		}

	}

	// 执行发布
	go doDeploy()

	return nil
}

func doDeploy() {
	sourcePath := os.Getenv("SOURCE_PATH")
	// 执行hexo命令发布
	err := os.Chdir(sourcePath + "/../..")
	if err != nil {
		log.Fatal("chdir error", err)
	}

	cmd := exec.Command("hexo", "g")
	// 获取输出对象，可以从该对象中读取输出结果
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("cmd stdoutPipe error", err)
	}
	defer func() {
		e := stdout.Close()
		if e != nil {
			log.Fatal("close stdout error", e)
		}
	}()
	// 运行命令
	if err := cmd.Start(); err != nil {
		log.Fatal("cmd start error", err)
	}
	res, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal("read stdout error", err)
	}
	log.Println(string(res))
	log.Println("deploy success!")
}
