package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zavier/blog-admin-backend/common"
	"net/http"
	"strings"
)

func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		origin := context.Request.Header.Get("Origin")
		var headerKeys []string
		for k, _ := range context.Request.Header {
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
			context.Header("Access-Control-Allow-Origin", "*")
			// 表示是否允许发送cookie
			context.Header("Access-Control-Allow-Credentials", "false")
			// 支持跨越请求的方法
			context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			// 服务器支持的所有头信息字段
			context.Header("Access-Control-Allow-Headers", headerStr)
			// 浏览器能拿到的扩展字段
			context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			// 指定本次预检请求的有效期，单位为秒  20天（1728000秒）
			context.Header("Access-Control-Max-Age", "172800")
		}

		// 如果是options方法， 则直接返回成功-200
		if context.Request.Method == "OPTIONS" {
			context.AbortWithStatusJSON(http.StatusOK, common.SuccessResult(true))
			return
		}

		context.Next()
	}
}
