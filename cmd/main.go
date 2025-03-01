package main

import database "todo/db"

func main() {
	database.InitDatabase()
	bootstrap()
}
