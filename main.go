package main

import (
	"prometheus-manager/initialize"
)

func main() {

	initialize.InitConfig()
	initialize.InitLogger()
	initialize.InitRoute()

}
