package server

import (
	"github.com/zavier/blog-admin-backend/constants"
	"github.com/zavier/blog-admin-backend/util"
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
---`

// 全部发布
func HexoDeployAll() {
	sourcePath := os.Getenv("SOURCE_PATH")
	log.Printf("hexo source path: %s\n", sourcePath)
	err := os.RemoveAll(sourcePath)
	if err != nil {
		log.Fatal("remove old path error", err)
	}
	err = os.MkdirAll(sourcePath, 777)
	if err != nil {
		log.Fatal("revert hexo source path error", err)
	}

	list := BlogList()
	blogMap := make(map[string]BlogBase)
	for _, b := range list {
		title := b.Title
		blogMap[title+".md"] = b
	}

	infos, err := ioutil.ReadDir(constants.BlogPath)
	if err != nil {
		log.Fatal("ReadDir error", err)
	}
	for _, f := range infos {
		if !util.Exists(f.Name()) {
			_, err := os.Create(f.Name())
			if err != nil {
				log.Fatal("create dest file error", err)
			}
		}

		bytes, err := ioutil.ReadFile(f.Name())
		if err != nil {
			log.Fatal("read file error", err)
		}
		blog := blogMap[f.Name()]
		newHeader := strings.ReplaceAll(hexoHeader, "${title}", blog.Title)
		newHeader = strings.ReplaceAll(newHeader, "${date}", time.Now().Format("2006-01-02 15:04:05"))
		s := blog.Categories
		if s != "" {
			newHeader = strings.ReplaceAll(newHeader, "${tags}", "["+blog.Categories+"]")
		} else {
			newHeader = strings.ReplaceAll(newHeader, "${tags}", "")
		}
		var context = newHeader + string(bytes)

		newFile, err := os.Create(sourcePath + "/" + f.Name())
		if err != nil {
			log.Fatal("create new file error", err)
		}
		_, err = newFile.WriteString(context)
		if err != nil {
			log.Fatal("write new file error", err)
		}
	}

	err = os.Chdir(sourcePath + "/../..")
	if err != nil {
		log.Fatal("chdir fail", err)
	}

	cmd := exec.Command("hexo", "g")
	// 获取输出对象，可以从该对象中读取输出结果
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("execute deploy command error", err)
	}
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	res, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal("read command result error", err)
	}
	log.Println(string(res))
}
