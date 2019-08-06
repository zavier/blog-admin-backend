package server

import (
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

// 全部发布
func HexoDeployAll(username string) error {
	// 清空文件夹
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
	baseBlogList, err := BlogList()
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
		newFile, err := os.OpenFile(sourcePath+"/"+blog.Title+".md", os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			return err
		}
		_, err = newFile.WriteString(context)
		if err != nil {
			return err
		}

	}

	// 执行hexo命令发布
	err = os.Chdir(sourcePath + "/../..")
	if err != nil {
		return err
	}

	cmd := exec.Command("hexo", "g")
	// 获取输出对象，可以从该对象中读取输出结果
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer func() {
		e := stdout.Close()
		if e != nil {
			log.Fatal("close stdout error", e)
		}
	}()
	// 运行命令
	if err := cmd.Start(); err != nil {
		return err
	}
	res, err := ioutil.ReadAll(stdout)
	if err != nil {
		return err
	}
	log.Println(string(res))
	return nil
}
