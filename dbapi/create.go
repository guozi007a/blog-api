package dbapi

import (
	"blog-api/global"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name   string
	Age    int
	Gender string
	City   string
}

func Create() {
	db := global.GlobalDB

	db.AutoMigrate(&User{})
}
