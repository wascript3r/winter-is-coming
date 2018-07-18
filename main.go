package main

import (
	"runtime"

	"github.com/wascript3r/winter-is-coming/server"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	server.Run(&server.Config{
		BX:         10,
		BY:         30,
		ZombieName: "night-king",
		Port:       3000,
	})
}
