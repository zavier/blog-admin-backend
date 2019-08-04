package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		var headerKeys []string
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			// 表示接受的域名
			c.Header("Access-Control-Allow-Origin", "http://114.67.66.81:9090")
			// 表示是否允许发送cookie
			c.Header("Access-Control-Allow-Credentials", "true")
			// 支持跨越请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			// 服务器支持的所有头信息字段
			c.Header("Access-Control-Allow-Headers", headerStr)
			// 浏览器能拿到的扩展字段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			// 指定本次预检请求的有效期，单位为秒  20天（1728000秒）
			c.Header("Access-Control-Max-Age", "172800")
		}
		c.Next()
	}
}
