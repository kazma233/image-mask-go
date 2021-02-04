package utils

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	logLogP  = "logs/high/log.log"
	highLogP = "logs/low/log.log"
)

var (
	// Llogger log
	Llogger *zap.Logger
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
		Filename:   logLogP,
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   false,
	}

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, os.Stdout, nonePriority),

		zapcore.NewCore(fileEncoder, zapcore.AddSync(highLogger), highPriority),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(lowLogger), lowPriority),
	)

	Llogger = zap.New(core)

	defer Llogger.Sync()
}
