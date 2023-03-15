package main

import (
	"go-midterm/router"
)

func main() {
	r := router.SetupRouter()
	r.Run()
}
