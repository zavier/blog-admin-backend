package middleware

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/zavier/blog-admin-backend/common"
	"github.com/zavier/blog-admin-backend/server"
	"log"
	"net/http"
	"time"
)

var jwtMiddleWare *jwt.GinJWTMiddleware

/**
iss (issuer)：签发人
exp (expiration time)：过期时间
sub (subject)：主题
aud (audience)：受众
nbf (Not Before)：生效时间
iat (Issued At)：签发时间
jti (JWT ID)：编号
*/
func init() {
	middleWare, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "login",
		Key:         []byte(common.JwtSecretKey),
		Timeout:     time.Hour * 12,
		MaxRefresh:  time.Hour,
		IdentityKey: common.JwtIdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*server.User); ok {
				return jwt.MapClaims{
					"aud": v.Name,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return claims["aud"].(string)
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var user server.User
			if err := c.ShouldBind(&user); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			correct, ex := user.CheckPassword()
			if ex != nil {
				return nil, ex
			}

			if correct {
				return &user, nil
			} else {
				return nil, jwt.ErrFailedAuthentication
			}
		},
		LoginResponse: func(c *gin.Context, code int, message string, t time.Time) {
			c.JSON(http.StatusOK, common.SuccessResult(message))
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			log.Printf("Unauthorized code:%d message%s\n", code, message)
			c.JSON(http.StatusOK, common.ErrorResult(common.StatusUnauthorized, "登录失效"))
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:g
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	if err != nil {
		log.Fatal("Jwt new error", err)
	}
	jwtMiddleWare = middleWare
}

func JwtMiddleware() *jwt.GinJWTMiddleware {
	return jwtMiddleWare
}
