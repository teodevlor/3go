package global

import (
	"go-structure/config"
	applogger "go-structure/internal/helper/logger"
)

var (
	Config config.Config
	Logger *applogger.LoggerZap
)
