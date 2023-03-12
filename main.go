package main

import (
	"go-midterm/database"
	"go-midterm/router"
)

func main() {
	r := router.SetupRouter()

	database.Connect()
	r.Run()
}
