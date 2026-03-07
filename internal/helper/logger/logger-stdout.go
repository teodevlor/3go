package logger

import (
	"go-structure/config"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewStdoutLoggerApplication(cfg config.Logger) *LoggerZap {
	level := parseLogLevel(cfg.LogLevel)
	encoder := getEncoderLog()

	writeSyncer := zapcore.AddSync(os.Stdout)
	core := zapcore.NewCore(encoder, writeSyncer, level)

	return &LoggerZap{zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))}
}

func NewStdoutChannelLogger(cfg config.Logger, channelName string) *LoggerZap {
	level := parseLogLevel(cfg.LogLevel)
	encoder := getEncoderLog()

	writeSyncer := zapcore.AddSync(os.Stdout)
	core := zapcore.NewCore(encoder, writeSyncer, level)

	log := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.Fields(zap.String("channel", channelName)),
	)

	return &LoggerZap{log}
}
