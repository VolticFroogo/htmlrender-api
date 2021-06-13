package main

import (
	"github.com/VolticFroogo/htmlrender-api/api"
	"github.com/VolticFroogo/htmlrender-api/rabbit"
)

func main() {
	rabbit.Init()

	api.Init()
}
