package kernel

import (
	"frame/routers"
	"frame/services/zlog"
	"net/http"
	"strconv"
	"sync"
)

const moduleName = "kernel"

var waitGroup *sync.WaitGroup
var httpServer *http.Server
var chanServices = make(chan bool)

func Run() {
	zlog.Warn(moduleName, "run begin")
	loadConfig()
	waitGroup = &sync.WaitGroup{}
	addr := ":" + strconv.FormatInt(kernelPort, 10)
	httpServer = &http.Server{Addr: addr, Handler: routers.GetEngine()}
	waitGroup.Add(1)
	go handleServer()
	waitGroup.Add(1)
	go handleSignal()
	startServices()
	<-chanServices
	stopServices()
	waitGroup.Wait()
	zlog.Warn(moduleName, "run end")
}
