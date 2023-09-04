package main

import (
	dbserver "blog-api/db_server"
)

func main() {
	dbserver.InitMySQL()
}
