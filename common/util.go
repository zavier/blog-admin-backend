package common

import (
	"os"
	"time"
)

// 文件是否存在
func ExistsPath(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

// 开启定时器
func StartTimer(f func(), d time.Duration) {
	go func() {
		c := time.Tick(d)
		for {
			<-c
			go f()
		}
	}()
}
