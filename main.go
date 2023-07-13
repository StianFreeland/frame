package main

import (
	"frame/kernel"
	"frame/services/config"
	"frame/services/mdb"
	"frame/services/zlog"
)

func main() {
	if !config.Init() {
		return
	}
	zlog.Start()
	defer zlog.Stop()
	zlog.Warn("program", "start ...")
	defer zlog.Warn("program", "stop ...")
	mdb.Start()
	defer mdb.Stop()
	kernel.Run()
}
