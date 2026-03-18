package log

import (
	"fmt"
	"os"
	"path/filepath"
	"practice/models"
	"practice/viper"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

func InitLogger() error {
	cfg := viper.GetConfig()
	if cfg == nil {
		return fmt.Errorf("viper init fail")
	}
	logConfig := cfg.Log
	logDir := filepath.Dir(logConfig.Output)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("shibai:%w", logDir, err)
	}
	core := buildLogCore(logConfig)
	options := []zap.Option{zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)}
	if viper.IsDev() {
		options = append(options, zap.AddCallerSkip(1))
	}
	Logger = zap.New(core, options...)
	zap.ReplaceGlobals(Logger)
	Logger.Info("logger init success", zap.String("log_level", logConfig.Level))
	return nil
}
func buildLogCore(logConfig models.LogConfig) zapcore.Core {
	level := zapcore.InfoLevel
	switch logConfig.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "panic":
		level = zapcore.PanicLevel
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
	}
	var encoder zapcore.Encoder
	if logConfig.Format == "json" || !viper.IsDev() {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	writers := []zapcore.WriteSyncer{}
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename: logConfig.Output,
		MaxSize:  logConfig.MaxSize,
		MaxAge:   logConfig.MaxAge,
		Compress: logConfig.Compress,
	})
	writers = append(writers, fileWriter)
	if viper.IsDev() {
		writers = append(writers, zapcore.AddSync(os.Stdout))
	}
	return zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writers...), level)
}
func custonTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}
func Debug(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Debug(msg, fields...)
	}
}
func Info(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Info(msg, fields...)
	}
}
func Warn(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Warn(msg, fields...)
	}
}

func Error(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Error(msg, fields...)
	}
}
func Panic(msg string, fields ...zap.Field) {
	if Logger != nil {
		Logger.Panic(msg, fields...)
	}
}
