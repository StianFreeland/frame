package kernel

import (
	"context"
	"frame/services/zlog"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

var chanSig = make(chan os.Signal)

func handleSignal() {
	zlog.Warn(moduleName, "handle signal begin")
	defer waitGroup.Done()
	signal.Notify(chanSig, syscall.SIGINT, syscall.SIGTERM)
	sig := <-chanSig
	zlog.Warn(moduleName, "handle signal", zap.Stringer("sig", sig))
	close(chanServices)
	if err := httpServer.Shutdown(context.Background()); err != nil {
		zlog.Error(moduleName, "handle signal", zap.Error(err))
		return
	}
	zlog.Warn(moduleName, "handle signal end")
}
