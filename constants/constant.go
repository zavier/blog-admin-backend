package constants

const (
	// 认证token
	TokenName = "token"
	// 博客目录
	BlogPath = "blog"
	// 博客管理文件
	BlogJsonFileName = "blogList.json"
	// 域名
	Domain = "127.0.0.1"
	// 密码文件
	PwdFilePath string = "pwd.txt"

	// cookie有效时间
	CookieMaxAge = 2592000
)

var AuthExcludeURI = [3]string{"/login", "/logout", "/register"}
