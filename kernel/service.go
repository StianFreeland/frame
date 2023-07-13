package kernel

import (
	"frame/services/cryptoService"
	"frame/services/dbService"
	"frame/services/zlog"
)

func startServices() {
	zlog.Warn(moduleName, "start services begin")
	cryptoService.Init()
	dbService.Init()
	zlog.Warn(moduleName, "start services end")
}

func stopServices() {
	zlog.Warn(moduleName, "stop services begin")
	zlog.Warn(moduleName, "stop services end")
}
