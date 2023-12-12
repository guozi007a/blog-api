package dbserver

import (
	"database/sql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"blog-api/db_server/tables"
	"blog-api/global"
)

func InitMySQL() {
	sqlDB, err := sql.Open("mysql", GeneratorDSN())
	if err != nil {
		panic(err)
	}
	mysqlDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	mysqlDB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&tables.DevLogs{},
		&tables.DateLogs{},
		&tables.SourceInfo{},
		&tables.IdInfo{},
		&tables.ActivityListInfo{},
		&tables.ChargeInfo{},
		&tables.Play_2399_Sign_List{},
		&tables.Play_2399_Turn_Cards{},
		&tables.Play_2399_Card_Prize{},
	)

	global.GlobalDB = mysqlDB
	global.PreDirs()
}
