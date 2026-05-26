package logger

import (
	"AIM/internal/common/config"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init() {
	LogConfig := config.Config.Log
	//编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	var level zapcore.Level
	switch LogConfig.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	//同时输出到文件和输出台
	file, _ := os.OpenFile(LogConfig.Filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(file)),
		level)
	log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	zap.ReplaceGlobals(log)

}
func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

// Info 普通信息
func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

// Warn 警告
func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

// Error 错误
func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

// WithField 携带字段
func WithField(key string, value interface{}) *zap.Logger {
	return log.With(zap.Any(key, value))
}
