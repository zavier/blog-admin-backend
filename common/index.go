package common

import (
	"github.com/zavier/blog-admin-backend/util"
	"log"
	"os"
	"strconv"
)

const indexFile = "index"

// 初始化博客的索引
func InitIndex(index string) {
	log.Printf("init index %s\n", index)
	exists := util.Exists(indexFile)
	if exists {
		err := os.Remove(indexFile)
		if err != nil {
			log.Fatal("remove file error", err)
		}
	}

	file, e := os.Create(indexFile)
	if e != nil {
		log.Fatal("create index file Error", e)
	}

	file, e = os.OpenFile(indexFile, os.O_RDWR, 777)
	if e != nil {
		log.Fatal("open index file Error", e)
	}
	defer file.Close()
	_, e = file.WriteString(index)
	if e != nil {
		log.Fatal("write index file Error", e)
	}
}

// 获取并生生下一个索引值
func GetAndIncrIndex() int {
	file, e := os.OpenFile(indexFile, os.O_RDWR, 777)
	if e != nil {
		log.Fatal("open index file Error", e)
	}
	defer file.Close()

	// 读取索引并加1
	buffer := make([]byte, 10)
	var n int
	n, e = file.Read(buffer)
	if e != nil {
		log.Fatal("read index file error", e)
	}
	number, e := strconv.Atoi(string(buffer[:n]))
	if e != nil {
		log.Fatal("atoi error", e)
	}
	newNumber := strconv.Itoa(number + 1)

	// 清空数据，重新写入
	_, e = file.WriteAt([]byte(newNumber), 0)
	if e != nil {
		log.Fatal("write index file error", e)
	}

	log.Printf("write new index:%s success\n", newNumber)
	return number + 1
}
