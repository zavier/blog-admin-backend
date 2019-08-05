package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zavier/blog-admin-backend/constants"
	"github.com/zavier/blog-admin-backend/handler"
	"log"
	"net/http"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			log.Printf("method is options")
			context.Next()
		}

		uri := context.Request.RequestURI
		log.Printf("uri is : %s\n", uri)
		var excludeURI = constants.AuthExcludeURI
		for _, item := range excludeURI {
			if strings.HasPrefix(uri, item) {
				context.Next()
				return
			}
		}

		token, err := context.Request.Cookie(constants.TokenName)
		log.Printf("token:%v token name:%s, token value:%s", token, token.Name, token.Value)
		if err != nil || token == nil || token.Value == "" {
			context.Abort()
			context.JSON(http.StatusOK, handler.ErrorResult(handler.StatusInternalServerError, "认证错误"))
			return
		}

		context.Next()
	}
}
