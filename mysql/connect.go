package mysql

import (
	"database/sql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"blog-api/global"
)

func InitMySQL() {
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	mysqlDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	global.GlobalDB = mysqlDB
}
