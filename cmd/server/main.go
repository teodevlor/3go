package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-structure/config"
	"go-structure/global"
	"go-structure/internal/helper/logger"
	"go-structure/internal/registry"

	"go.uber.org/zap"
)

func main() {
	// 1. Build DI
	registry.BuildDependencyInjectContainer()

	// 2. Resolve dependencies
	cfg := registry.GetDependency(registry.ConfigDIName).(*config.Config)
	router := registry.GetDependency(registry.ApiDIName).(http.Handler)

	// 3. Global config & logger
	global.Config = *cfg
	global.Logger = logger.NewLoggerApplication(cfg.Logger)
	zap.ReplaceGlobals(global.Logger.Logger)

	// Logger theo channel (auth -> auth.log, http -> http.log, ...)
	global.ChannelLoggers = make(map[string]*logger.LoggerZap)
	for name, filePath := range cfg.Logger.Channels {
		if name == "" || filePath == "" {
			continue
		}
		global.ChannelLoggers[name] = logger.NewChannelLogger(cfg.Logger, name, filePath)
	}

	// 4. HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	// 5. Run server
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.S().Panic("HTTP server failed", zap.Error(err))
		}
	}()

	// 6. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	zap.S().Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		zap.S().Error("Server shutdown failed", zap.Error(err))
	}

	zap.S().Info("Server exited properly")
}
