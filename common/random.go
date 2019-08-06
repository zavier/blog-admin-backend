package common

import (
	"math/rand"
	"strconv"
	"time"
)

// 生产随机数字字符串 长度为24
func Random24NumberString() string {
	nano := time.Now().UnixNano()
	prefix := strconv.FormatInt(nano, 10)

	i := randInt(10000, 99999)
	suffix := strconv.Itoa(i)

	return prefix + suffix
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
