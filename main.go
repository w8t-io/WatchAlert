package main

import (
	"watchAlert/initialization"
	"watchAlert/internal/global"
)

var Version string

func main() {

	global.Version = Version
	initialization.InitBasic()
	initialization.InitRoute()

}
