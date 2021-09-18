package main

import (
	"chat/db"
	"chat/server"
)

func main() {
	db.ConnectScylladb("localhost:9042")
	server.Router()
}
