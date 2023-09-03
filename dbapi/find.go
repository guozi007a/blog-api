package dbapi

import (
	"blog-api/global"
	"fmt"
)

type Msg struct {
	Name   string
	Age    int
	Gender string
	City   string
}

func Find() {
	db := global.GlobalDB

	var user Msg

	// result := db.Table("users").Limit(1).Find(&user)
	result := db.Table("users").Order("id desc").First(&user)
	if result.Error != nil {
		panic(result.Error)
	}

	fmt.Printf("%+v", user)
}

func FindAll() {
	db := global.GlobalDB

	var user []Msg

	result := db.Table("users").Find(&user)
	if result.Error != nil {
		panic(result.Error)
	}

	fmt.Println(user)
}
