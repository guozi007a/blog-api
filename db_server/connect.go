package dbserver

import (
	"database/sql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"blog-api/db_server/tables"
	"blog-api/global"
)

func InitMySQL() {
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	mysqlDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		// 禁止自动创建外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		panic(err)
	}

	mysqlDB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&tables.DevLogs{},
		&tables.DateLogs{},
	)

	global.GlobalDB = mysqlDB
}
