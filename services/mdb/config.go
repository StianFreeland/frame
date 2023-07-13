package mdb

import (
	"frame/services/zlog"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var mongoScheme string
var mongoAddress string
var mongoUsername string
var mongoPassword string

func loadConfig() {
	zlog.Warn(moduleName, "load config ...")
	mongoScheme = viper.GetString("mongo.scheme")
	if mongoScheme == "" {
		zlog.Fatal(moduleName, "load config", zap.String("scheme", mongoScheme))
	}
	zlog.Warn(moduleName, "load config", zap.String("scheme", mongoScheme))
	mongoAddress = viper.GetString("mongo.address")
	if mongoAddress == "" {
		zlog.Fatal(moduleName, "load config", zap.String("address", mongoAddress))
	}
	zlog.Warn(moduleName, "load config", zap.String("address", mongoAddress))
	mongoUsername = viper.GetString("mongo.username")
	zlog.Warn(moduleName, "load config", zap.String("username", mongoUsername))
	mongoPassword = viper.GetString("mongo.password")
	zlog.Warn(moduleName, "load config", zap.String("password", mongoPassword))
}
