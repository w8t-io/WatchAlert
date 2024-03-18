package main

import (
	"watchAlert/alert/eval"
	"watchAlert/initialize"
)

func main() {

	initialize.InitConfig()
	initialize.InitLogger()
	initialize.InitClient()
	initialize.InitResource()
	eval.Initialize()
	initialize.InitRoute()

}
