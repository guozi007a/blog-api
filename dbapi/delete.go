package dbapi

import (
	"blog-api/global"
)

func DeleteOne() {
	db := global.GlobalDB

	var user User
	db.Table("users").First(&user, 3)
	db.Unscoped().Delete(&user)
	// db.Where("id = ?", 1).Unscoped().Delete(&User{})
}
