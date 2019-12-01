package server

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/zavier/blog-admin-backend/common"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)
	if err := prepareData(); err != nil {
		log.Fatal("prepareData error", err)
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
	Content string `form:"content" json:"content" binding:"required"`
}

// 初始化准备数据
func prepareData() error {
	// 初始化博客路径
	exist, e := common.ExistsPath(common.BlogPath)
	if e != nil {
		return e
	}
	if !exist {
		err := os.Mkdir(common.BlogPath, 0777)
		if err != nil {
			log.Printf("create blog path error:%s\n", err.Error())
			return err
		}
	}

	// 初始化博客管理文件
	exist, e = common.ExistsPath(common.BlogManageFileName)
	if e != nil {
		return e
	}
	if !exist {
		_, err := os.Create(common.BlogManageFileName)
		if err != nil {
			return err
		}
	}

	// 初始化索引
	if e = common.InitBlogIndex(strconv.Itoa(0)); e != nil {
		return e
	}
	return nil
}

func isExistBlogTitle(blog Blog) (bool, error) {
	blogList, e := BlogList()
	if e != nil {
		return true, e
	}
	if len(blogList) == 0 {
		return false, nil
	}

	for _, b := range blogList {
		if b.Title == blog.Title && b.Author == blog.Author {
			return true, nil
		}
	}
	return false, nil
}

// 保存博客
func (blog *Blog) SaveBlog() error {
	// 创建要保存的文件
	fileName := common.Random24NumberString() + ".md"
	filePath := common.BlogPath + "/" + fileName
	exist, e := isExistBlogTitle(*blog)
	if e != nil {
		return e
	}
	if exist {
		log.Printf("filePath:%s has existed", filePath)
		return errors.New("博客名称已存在")
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
	_, err = file.WriteString(blog.Content)
	if err != nil {
		log.Printf("write file error:%s\n", err.Error())
		return err
	}

	// 保存记录信息
	jsonFile, err := os.OpenFile(common.BlogManageFileName, os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	defer func() {
		err := jsonFile.Close()
		if err != nil {
			log.Fatal("close file error")
		}
	}()
	blog.Id, err = common.GetAndIncrBlogIndex()
	if err != nil {
		return err
	}
	path, err := filepath.Abs(filepath.Dir(filePath))
	if err != nil {
		return err
	}
	blog.Location = path + "/" + fileName
	blog.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	blog.UpdateTime = time.Now().Format("2006-01-02 15:04:05")

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
	_, err = jsonFile.WriteString(context + "\n")
	if err != nil {
		return err
	}

	// 将更新状态置为已更新
	common.BlogUpdated = true
	return nil
}

// 更新博客
func (blog *Blog) UpdateBlog() error {
	// 更新记录信息
	recordFile, err := os.OpenFile(common.BlogManageFileName, os.O_RDWR, 0777)
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

		if blog.Id == b.Id {
			// 找到文件进行写入
			location := b.Location
			blogFile, err := os.OpenFile(location, os.O_WRONLY|os.O_TRUNC, 0777)
			if err != nil {
				log.Printf("open file error: %s\n", err.Error())
				return err
			}
			defer blogFile.Close()

			if _, err = blogFile.WriteString(blog.Content); err != nil {
				log.Printf("write file error:%s\n", err.Error())
				return err
			}

			// 更新信息
			b.Title = blog.Title
			b.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
			b.Author = blog.Author
			b.Categories = blog.Categories
		}
		blogList = append(blogList, b)
	}

	if _, err := recordFile.Seek(0, 0); err != nil {
		return err
	}

	err = recordFile.Truncate(0)
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

	// 将更新状态置为已更新
	common.BlogUpdated = true
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
				Content:  string(bytes),
			}
			return blog, nil
		}
	}
	return Blog{}, errors.New("此博客不存在")
}

// 删除博客
func DelBlog(id int) error {
	// 只删除记录信息
	recordFile, err := os.OpenFile(common.BlogManageFileName, os.O_RDWR, 0777)
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

		if id != b.Id {
			blogList = append(blogList, b)
		}
	}

	if _, err := recordFile.Seek(0, 0); err != nil {
		return err
	}
	err = recordFile.Truncate(0)
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

// 查询博客列表
func BlogList() ([]BlogBase, error) {
	file, e := os.OpenFile(common.BlogManageFileName, os.O_RDONLY, 777)
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
	sortByIdAndTotoPrimary(blogList)
	return blogList, nil
}

// 按照ID倒序，但是todo标签的优先级最高
func sortByIdAndTotoPrimary(blogList []BlogBase) {
	// 按照ID倒序
	if len(blogList) > 0 {
		sort.Slice(blogList, func(i, j int) bool {
			return blogList[i].Id > blogList[j].Id
		})
		// todo标签优先
		for i := 0; i < len(blogList); i++ {
			blog := blogList[i]
			if strings.Contains(blog.Categories, "todo") {
				for j := i; i > 0; i-- {
					blogList[j] = blogList[j-1]
				}
				blogList[0] = blog
			}
		}
	}
}
