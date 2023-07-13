package kernel

import (
	"frame/services/zlog"
	"go.uber.org/zap"
	"net/http"
	"syscall"
)

func handleServer() {
	defer waitGroup.Done()
	zlog.Warn(moduleName, "handle server begin")
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		zlog.Error(moduleName, "handle server", zap.Error(err))
		chanSig <- syscall.SIGTERM
		return
	}
	zlog.Warn(moduleName, "handle server end")
}
