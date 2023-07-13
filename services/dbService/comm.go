package dbService

import (
	"frame/services/zlog"
)

const moduleName = "db service"

func Init() {
	zlog.Warn(moduleName, "init ...")
	initMgmt()
	initMenu()
}
