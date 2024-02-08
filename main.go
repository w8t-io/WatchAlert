package main

import (
	"watchAlert/alert/eval"
	"watchAlert/initialize"
)

func main() {

	initialize.InitConfig()
	initialize.InitLogger()
	initialize.InitClient()
	eval.Initialize()
	initialize.InitRoute()

}
