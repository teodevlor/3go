package global

import (
	"go-structure/config"
	applogger "go-structure/internal/helper/logger"
)

var (
	Config         config.Config
	Logger         *applogger.LoggerZap
	ChannelLoggers map[string]*applogger.LoggerZap
)

func GetChannelLogger(channel string) *applogger.LoggerZap {
	if ChannelLoggers != nil {
		if l, ok := ChannelLoggers[channel]; ok && l != nil {
			return l
		}
	}
	return Logger
}
