package server

import (
	"bufio"
	"encoding/json"
	"errors"
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
	if err := SyncBlogJsonFile(); err != nil {
		log.Fatal("SyncBlogJsonFile error", err)
	}

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
func SyncBlogJsonFile() error {
	// 初始化创建路径及必须文件
	exist, e := util.Exists(constants.BlogPath)
	if e != nil {
		return e
	}
	if !exist {
		err := os.Mkdir(constants.BlogPath, 0777)
		if err != nil {
			log.Printf("create blog path error:%s\n", err.Error())
			return err
		}
	}

	firstStart := false

	exist, e = util.Exists(constants.BlogJsonFileName)
	if e != nil {
		return e
	}
	if !exist {
		_, err := os.Create(constants.BlogJsonFileName)
		if err != nil {
			return err
		}
		firstStart = true
	}
	jsonFile, err := os.OpenFile(constants.BlogJsonFileName, os.O_RDWR, 777)
	if err != nil {
		return err
	}
	defer func() {
		if err := jsonFile.Close(); err != nil {
			log.Fatal("close file error")
		}
	}()

	// 是否是第一次启动，第一次启动时才需要同步文件夹中的内容
	// 否则直接返回，不进行后续操作
	if !firstStart {
		return nil
	}

	// 初始读取文件夹下所有文件
	infos, err := ioutil.ReadDir(constants.BlogPath)
	if err != nil {
		return err
	}
	var index = 0
	for _, file := range infos {
		name := file.Name()
		path, err := filepath.Abs(filepath.Dir(file.Name()))
		if err != nil {
			return err
		}
		index++
		blogBase := BlogBase{
			Id:       index,
			Title:    name,
			Location: path,
		}
		bytes, err := json.Marshal(blogBase)
		if err != nil {
			return err
		}
		context := string(bytes)
		_, err = jsonFile.WriteString(context + "\n")
		if err != nil {
			return err
		}
	}

	// 初始化索引
	if e = common.InitIndex(strconv.Itoa(index)); e != nil {
		return e
	}
	return nil
}

// 保存博客
func (blog Blog) SaveBlog() error {
	// 创建要保存的文件
	fileName := blog.Title + ".md"
	filePath := constants.BlogPath + "/" + fileName
	exist, e := util.Exists(filePath)
	if e != nil {
		return e
	}
	if exist {
		log.Printf("filePath:%s has existed", filePath)
		return errors.New("文件路径已存在")
	}
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("create file error: %s\n", err.Error())
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal("close file error")
		}
	}()

	// 文件写入
	_, err = file.WriteString(blog.Context)
	if err != nil {
		log.Printf("write file error:%s\n", err.Error())
		return err
	}

	// 保存记录信息
	jsonFile, err := os.OpenFile(constants.BlogJsonFileName, os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	defer func() {
		err := jsonFile.Close()
		if err != nil {
			log.Fatal("close file error")
		}
	}()
	blog.Id, err = common.GetAndIncrIndex()
	if err != nil {
		return err
	}
	path, err := filepath.Abs(filepath.Dir(filePath))
	if err != nil {
		return err
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
		return err
	}
	context := string(bytes)
	_, err = file.WriteString(context + "\n")
	if err != nil {
		return err
	}
	return nil
}

// 更新博客
func (blog Blog) UpdateBlog() error {
	var blogId = blog.Id
	if blogId <= 0 {
		return errors.New("ID不能小于等于0")
	}

	// 保存文件
	var fileName = blog.Title + ".md"
	filePath := constants.BlogPath + "/" + fileName
	exists, e := util.Exists(filePath)
	if e != nil {
		return e
	}
	if !exists {
		log.Printf("filePath:%s does not exist, create.", filePath)
		if _, err := os.Create(filePath); err != nil {
			return err
		}
	}
	blogFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		log.Printf("open file error: %s\n", err.Error())
		return err
	}
	defer func() {
		if err := blogFile.Close(); err != nil {
			log.Fatal("close file error")
		}
	}()
	if _, err = blogFile.WriteString(blog.Context); err != nil {
		log.Printf("write file error:%s\n", err.Error())
		return err
	}

	// 更新记录信息
	recordFile, err := os.OpenFile(constants.BlogJsonFileName, os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	defer func() {
		if err := recordFile.Close(); err != nil {
			log.Fatal("close file error")
		}
	}()
	scanner := bufio.NewScanner(recordFile)
	var blogList = make([]BlogBase, 0)
	for scanner.Scan() {
		blogJson := scanner.Text()
		var b BlogBase
		if err := json.Unmarshal([]byte(blogJson), &b); err != nil {
			return err
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
		blogList = append(blogList, b)
	}

	log.Printf("blogs : %v\n", blogList)
	recordFile, err = os.OpenFile(constants.BlogJsonFileName, os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	for _, blogBase := range blogList {
		bytes, err := json.Marshal(blogBase)
		if err != nil {
			return err
		}
		if _, err = recordFile.WriteString(string(bytes) + "\n"); err != nil {
			return err
		}
	}

	return nil
}

// 通过ID查询博客信息
func GetBlog(id int) (Blog, error) {
	bases, e := BlogList()
	if e != nil {
		return Blog{}, e
	}
	list := bases
	for _, blogBase := range list {
		if blogBase.Id == id {
			location := blogBase.Location
			file, e := os.OpenFile(location, os.O_RDONLY, 777)
			if e != nil {
				return Blog{}, e
			}
			bytes, e := ioutil.ReadAll(file)
			if e != nil {
				return Blog{}, e
			}
			blog := Blog{
				BlogBase: blogBase,
				Context:  string(bytes),
			}
			return blog, nil
		}
	}
	return Blog{}, errors.New("此博客不存在")
}

// 查询博客列表
func BlogList() ([]BlogBase, error) {
	file, e := os.OpenFile(constants.BlogJsonFileName, os.O_RDONLY, 777)
	if e != nil {
		return nil, e
	}
	var blogList = make([]BlogBase, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		var blogBase BlogBase
		e := json.Unmarshal([]byte(text), &blogBase)
		if e != nil {
			return nil, e
		}
		blogList = append(blogList, blogBase)
	}
	return blogList, nil
}
