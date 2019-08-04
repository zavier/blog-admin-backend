package constants

const (
	// 认证token
	TokenName = "token"
	// 博客目录
	BlogPath = "blog"
	// 博客管理文件
	BlogJsonFileName = "blogList.json"
)

var AuthExcludeURI = [3]string{"/login", "/logout", "/register"}
