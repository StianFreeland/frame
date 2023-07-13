package kernel

import (
	"frame/services/zlog"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var kernelPort int64

func loadConfig() {
	zlog.Warn(moduleName, "load config begin")
	kernelPort = viper.GetInt64("kernel.port")
	if kernelPort == 0 {
		zlog.Fatal(moduleName, "load config", zap.Int64("port", kernelPort))
	}
	zlog.Warn(moduleName, "load config", zap.Int64("port", kernelPort))
	zlog.Warn(moduleName, "load config end")
}
