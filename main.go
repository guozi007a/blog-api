package main

import (
	"blog-api/dbapi"
	"blog-api/mysql"
)

func main() {
	mysql.InitMySQL()
	dbapi.DeleteOne()
}
