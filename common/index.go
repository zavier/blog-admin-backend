package common

import (
	"github.com/zavier/blog-admin-backend/constants"
	"github.com/zavier/blog-admin-backend/util"
	"log"
	"os"
	"strconv"
)

// 初始化博客的索引
func InitIndex(index string) error {
	log.Printf("init index %s\n", index)
	exists, err := util.Exists(constants.IndexFile)
	if err != nil {
		return err
	}
	if !exists {
		if _, e := os.Create(constants.IndexFile); e != nil {
			return e
		}
	}

	file, e := os.OpenFile(constants.IndexFile, os.O_WRONLY|os.O_TRUNC, 777)
	if e != nil {
		return e
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal("close file error")
		}
	}()

	if _, e = file.WriteString(index); e != nil {
		return e
	}
	return nil
}

// 获取并生生下一个索引值
func GetAndIncrIndex() (index int, e error) {
	file, e := os.OpenFile(constants.IndexFile, os.O_RDWR, 777)
	if e != nil {
		return 0, e
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal("close file error")
		}
	}()

	// 读取索引并加1
	buffer := make([]byte, 10)
	var n int
	n, e = file.Read(buffer)
	if e != nil {
		return 0, e
	}
	number, e := strconv.Atoi(string(buffer[:n]))
	if e != nil {
		return 0, e
	}
	newNumber := strconv.Itoa(number + 1)

	// 清空数据，重新写入
	if _, e = file.WriteAt([]byte(newNumber), 0); e != nil {
		return 0, e
	}

	log.Printf("write new index:%s success\n", newNumber)
	return number + 1, nil
}
