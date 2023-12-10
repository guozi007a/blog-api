package dbserver

import (
	"blog-api/global"
	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	user     string                                       //用户名
	password string                                       //密码
	host     string                                       // 域名
	port     = "3306"                                     // 数据库端口
	dbname   string                                       //数据库名称
	extra    = "charset=utf8mb4&parseTime=True&loc=Local" // 其他配置，如果使用世界标准时间，就改为loc=UTC
)

func GeneratorDSN() string {
	mode := gin.Mode()

	// fmt.Printf("当前环境: %s\n", mode)

	switch mode {
	case "release":
		password = global.MYSQL_ROOT_PASSWORD_PROD
		host = global.MYSQL_HOST_PROD
	default:
		password = global.MYSQL_ROOT_PASSWORD_DEV
		host = global.MYSQL_HOST_DEV
	}

	user = global.MYSQL_ROOT_USER
	dbname = global.MYSQL_DATABASE

	// fmt.Printf("user: %s\n password: %s\n host: %s\n dbname: %s\n", user, password, host, dbname)

	// "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", user, password, host, port, dbname, extra)

	return dsn
}
