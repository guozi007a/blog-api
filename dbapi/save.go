package dbapi

import (
	"blog-api/global"
)

func Save() {
	db := global.GlobalDB.Model(&User{})

	var user User

	db.First(&user, 2)

	user.Name = "world"

	db.Save(&user)
}
