package logger

import (
	"go-structure/config"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerZap struct {
	*zap.Logger
}

func NewLoggerApplication(config config.Logger) *LoggerZap {
	level := parseLogLevel(config.LogLevel)
	encoder := getEncoderLog()
	hook := lumberjack.Logger{
		Filename:   config.FileLog,
		MaxSize:    config.MaxSize, // megabytes
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge, // days
		Compress:   config.Compress,
	}

	writeSyncer := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(&hook),
	)
	core := zapcore.NewCore(encoder, writeSyncer, level)

	return &LoggerZap{zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))}
}

func NewChannelLogger(config config.Logger, channelName string, filePath string) *LoggerZap {
	level := parseLogLevel(config.LogLevel)
	encoder := getEncoderLog()
	hook := lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}
	writeSyncer := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(&hook),
	)
	core := zapcore.NewCore(encoder, writeSyncer, level)
	log := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel), zap.Fields(zap.String("channel", channelName)))
	return &LoggerZap{log}
}

func parseLogLevel(logLevel string) zapcore.Level {
	switch logLevel {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func getEncoderLog() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()

	// format time
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// caller
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.CallerKey = "file"

	// message
	encoderConfig.MessageKey = "message"
	return zapcore.NewJSONEncoder(encoderConfig)
}
