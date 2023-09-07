package dbserver

import (
	"fmt"
)

var (
	user     = "root" //用户名
	password = "123456"
	// password = "jL9@xZ2!tQ6*kT1)>D5+nA^|lN4~zF5$" // 密码
	// host     = "127.0.0.1"                                // 域名
	host   = "mysql"
	port   = "3306"                                     // 数据库端口
	dbname = "vm50"                                     //数据库名称
	extra  = "charset=utf8mb4&parseTime=True&loc=Local" // 其他配置
)

// "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
var dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", user, password, host, port, dbname, extra)
