package zlog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

const moduleName = "zlog"

var zapLogger *zap.Logger

func Start() {
	loadConfig()
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05.000")
	fileCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(config),
		zapcore.Lock(&logWriter{}),
		logLevel,
	)
	cores := []zapcore.Core{fileCore}
	if isDebug {
		stdoutCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(config),
			zapcore.Lock(&writerWrapper{os.Stdout}),
			logLevel,
		)
		cores = append(cores, stdoutCore)
	}
	zapLogger = zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.PanicLevel),
	)
}

func Stop() {
	if err := zapLogger.Sync(); err != nil {
		fmt.Println(moduleName, "sync, error:", err)
	}
}

func GetLevel() int64 {
	return int64(logLevel.Level())
}

func UpdateLevel(level int64) {
	logLevel.SetLevel(zapcore.Level(level))
}

func Fatal(prefix string, msg string, fields ...zap.Field) {
	zapLogger.Fatal(getLogMsg(prefix, msg), fields...)
}

func Error(prefix string, msg string, fields ...zap.Field) {
	zapLogger.Error(getLogMsg(prefix, msg), fields...)
}

func Warn(prefix string, msg string, fields ...zap.Field) {
	zapLogger.Warn(getLogMsg(prefix, msg), fields...)
}

func Info(prefix string, msg string, fields ...zap.Field) {
	zapLogger.Info(getLogMsg(prefix, msg), fields...)
}

func Debug(prefix string, msg string, fields ...zap.Field) {
	zapLogger.Debug(getLogMsg(prefix, msg), fields...)
}

func getLogMsg(prefix string, msg string) string {
	return "[ " + prefix + " ] " + msg
}
