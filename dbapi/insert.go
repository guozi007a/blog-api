package dbapi

import (
	"blog-api/global"
)

func InsertOne() {
	db := global.GlobalDB

	user := User{
		Name:   "迪丽热巴",
		Age:    18,
		Gender: "女",
		City:   "上海",
	}

	result := db.Create(&user)
	if result.Error != nil {
		panic(result.Error)
	}
}

func InsertMany() {
	db := global.GlobalDB

	users := []*User{
		{
			Name:   "小明",
			Age:    14,
			Gender: "男",
			City:   "北京",
		},
		{
			Name:   "张三",
			Age:    30,
			Gender: "男",
			City:   "钝角",
		},
	}

	result := db.Create(users)
	if result.Error != nil {
		panic(result.Error)
	}
}

func InsertMore() {
	db := global.GlobalDB

	users := []User{
		{
			Name:   "cxh",
			Age:    31,
			Gender: "女",
			City:   "杭州",
		},
		{
			Name:   "lbn",
			Age:    38,
			Gender: "女",
			City:   "杭州",
		},
	}

	db.Create(&users)
}
