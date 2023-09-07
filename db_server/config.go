package dbserver

import (
	"fmt"
)

var (
	user     = "root" //用户名
	password = "123456"
	host     = "localhost" // dev域名
	// host   = "mysql" // prod域名 跟yml中mysql名称保持一致
	port   = "3306"                                     // 数据库端口
	dbname = "vm50"                                     //数据库名称
	extra  = "charset=utf8mb4&parseTime=True&loc=Local" // 其他配置
)

// "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
var dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", user, password, host, port, dbname, extra)
