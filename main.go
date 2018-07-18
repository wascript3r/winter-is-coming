package main

import (
	"runtime"

	"github.com/wascript3r/winter-is-coming/server"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	server.Run()
}
