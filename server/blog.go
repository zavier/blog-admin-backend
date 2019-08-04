package server

import (
	"bufio"
	"encoding/json"
	"github.com/zavier/blog-admin-backend/common"
	"github.com/zavier/blog-admin-backend/constants"
	"github.com/zavier/blog-admin-backend/util"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)
	SyncBlogJsonFile()
}

type BlogBase struct {
	Id         int    `json:"id"`
	Title      string `form:"title" json:"title" binding:"required"`
	Location   string `json:"location"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
	Author     string `json:"author"`
	Categories string `json:"categories"`
}

type Blog struct {
	BlogBase
	Context string `form:"context" json:"context" binding:"required"`
}

func (blog Blog) ToString() string {
	bytes, _ := json.Marshal(blog)
	return string(bytes)
}

// 同步博客json文件及博客索引
func SyncBlogJsonFile() {
	// 初始化创建路径及必须文件
	if util.Exists(constants.BlogPath) {
		err := os.RemoveAll(constants.BlogPath)
		if err != nil {
			log.Fatal("create blog path error ", err)
		}
	}
	if !util.Exists(constants.BlogPath) {
		err := os.Mkdir(constants.BlogPath, 0777)
		if err != nil {
			log.Printf("create blog path error:%s\n", err.Error())
			return
		}
	}

	if util.Exists(constants.BlogJsonFileName) {
		err := os.Remove(constants.BlogJsonFileName)
		if err != nil {
			log.Fatal("create blog json file error ", err)
		}
	}
	if !util.Exists(constants.BlogJsonFileName) {
		_, err := os.Create(constants.BlogJsonFileName)
		if err != nil {
			log.Fatal("create file error", err)
		}
	}
	jsonFile, err := os.OpenFile(constants.BlogJsonFileName, os.O_RDWR, 777)
	if err != nil {
		log.Fatal("open file error", err)
	}
	defer func() {
		e := jsonFile.Close()
		if e != nil {
			log.Fatal("close file error", e)
		}
	}()

	// 初始读取文件夹下所有文件
	infos, err := ioutil.ReadDir(constants.BlogPath)
	if err != nil {
		log.Printf("read blog path error: %s\n", err)
		return
	}
	var index = 0
	for _, file := range infos {
		name := file.Name()
		path, err := filepath.Abs(filepath.Dir(file.Name()))
		if err != nil {
			log.Fatal("get abs path error", err)
		}
		index++
		blogBase := BlogBase{
			Id:       index,
			Title:    name,
			Location: path,
		}
		bytes, err := json.Marshal(blogBase)
		if err != nil {
			log.Fatal("json Marshal error", err.Error())
		}
		context := string(bytes)
		_, err = jsonFile.WriteString(context + "\n")
		if err != nil {
			log.Fatal("write file error", err)
		}
	}

	// 初始化索引
	common.InitIndex(strconv.Itoa(index))
}

// 保存博客
func (blog Blog) SaveBlog() error {
	// 创建要保存的文件
	fileName := blog.Title + ".md"
	filePath := constants.BlogPath + "/" + fileName
	if util.Exists(filePath) {
		log.Printf("filePath:%s has existed", filePath)
		return common.CheckError{
			Message: "文件路径已存在",
		}
	}
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("create file error: %s\n", err.Error())
		return err
	}
	defer file.Close()

	// 文件写入
	_, err = file.WriteString(blog.Context)
	if err != nil {
		log.Printf("write file error:%s\n", err.Error())
		return err
	}

	// 保存记录信息
	file, err = os.OpenFile(constants.BlogJsonFileName, os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	blog.Id = common.GetAndIncrIndex()
	path, err := filepath.Abs(filepath.Dir(filePath))
	if err != nil {
		log.Fatal("get abs path error", err)
	}
	blog.Location = path + "/" + fileName

	blogBase := BlogBase{
		Id:         blog.Id,
		Title:      blog.Title,
		Location:   blog.Location,
		CreateTime: blog.CreateTime,
		UpdateTime: blog.UpdateTime,
		Author:     blog.Author,
		Categories: blog.Categories,
	}

	bytes, err := json.Marshal(blogBase)
	if err != nil {
		log.Fatal("json Marshal error", err.Error())
	}
	context := string(bytes)
	_, err = file.WriteString(context + "\n")
	if err != nil {
		log.Fatal("write file error", err)
	}
	return nil
}

// 更新博客
func (blog Blog) UpdateBlog() error {
	var blogId = blog.Id
	if blogId <= 0 {
		log.Fatal("blog Id is zero")
	}

	// 保存文件
	var fileName = blog.Title + ".md"
	filePath := constants.BlogPath + "/" + fileName
	if !util.Exists(filePath) {
		log.Printf("filePath:%s does not exist, create.", filePath)
		_, err := os.Create(filePath)
		if err != nil {
			log.Fatalf("create file error", err)
		}

	}
	blogFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		log.Printf("open file error: %s\n", err.Error())
		return err
	}
	defer blogFile.Close()
	_, err = blogFile.WriteString(blog.Context)
	if err != nil {
		log.Printf("write file error:%s\n", err.Error())
		return err
	}

	// 更新记录信息
	recordFile, err := os.OpenFile(constants.BlogJsonFileName, os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	defer recordFile.Close()
	scanner := bufio.NewScanner(recordFile)
	var blogsList = make([]BlogBase, 0)
	for scanner.Scan() {
		blogJson := scanner.Text()
		var b BlogBase
		err = json.Unmarshal([]byte(blogJson), &b)
		if err != nil {
			log.Fatal("Unmarshal error", err)
		}

		if blogId == b.Id {
			path, err := filepath.Abs(filepath.Dir(filePath))
			if err != nil {
				log.Fatal("get abs path error", err)
			}
			b.Location = path + "/" + fileName
			b.Title = blog.Title
			b.CreateTime = blog.CreateTime
			b.UpdateTime = blog.UpdateTime
			b.Author = blog.Author
			b.Categories = blog.Categories
		}
		blogsList = append(blogsList, b)
	}

	log.Printf("blogs : %v\n", blogsList)
	recordFile, err = os.OpenFile(constants.BlogJsonFileName, os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	for _, blogBase := range blogsList {
		bytes, err := json.Marshal(blogBase)
		if err != nil {
			log.Fatal("Marshal error", err)
		}
		_, err = recordFile.WriteString(string(bytes) + "\n")
		if err != nil {
			log.Fatal("write blog record file error", err)
		}
	}

	return nil
}

// 查询博客列表
func BlogList() []BlogBase {
	file, e := os.OpenFile(constants.BlogJsonFileName, os.O_RDONLY, 777)
	if e != nil {
		log.Fatal("open file error", e)
	}
	var blogList = make([]BlogBase, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		var blogBase BlogBase
		e := json.Unmarshal([]byte(text), &blogBase)
		if e != nil {
			log.Fatal("Unmarshal error", e)
		}
		blogList = append(blogList, blogBase)
	}
	return blogList
}
