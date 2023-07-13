package zlog

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var isDebug bool
var logLevel zap.AtomicLevel

func loadConfig() {
	isDebug = viper.GetBool("zap.debug")
	logLevel = zap.NewAtomicLevelAt(zapcore.Level(viper.GetInt64("zap.level")))
}
