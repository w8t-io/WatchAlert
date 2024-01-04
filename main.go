package main

import (
	"watchAlert/initialize"
)

func main() {

	initialize.InitConfig()
	initialize.InitLogger()
	initialize.InitClient()
	initialize.InitRoute()

}
