package utils

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	highLogP = "logs/high/log.log"
	lowLogP  = "logs/low/log.log"
)

var (
	core zapcore.Core
)

func init() {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	nonePriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return true
	})

	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	highLogger := &lumberjack.Logger{
		Filename:   highLogP,
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   false,
	}

	lowLogger := &lumberjack.Logger{
		Filename:   lowLogP,
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   false,
	}

	core = zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, os.Stdout, nonePriority),

		zapcore.NewCore(fileEncoder, zapcore.AddSync(highLogger), highPriority),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(lowLogger), lowPriority),
	)
}

// GetLog 获取日志输出对象
func GetLog() *zap.Logger {
	return zap.New(core)
}

// GetSugaredLogger 获取包装日志输出对象
func GetSugaredLogger() *zap.SugaredLogger {
	return zap.New(core).Sugar()
}
