package main

import (
	"blog-api/api"
	dbserver "blog-api/db_server"
)

func main() {
	dbserver.InitMySQL()
	api.Server()
}
