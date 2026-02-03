package registry

import (
	"go-structure/config"
	"go-structure/internal/adapter"
	"go-structure/pkg/telegram"

	"github.com/sarulabs/di"
)

const (
	TelegramAdapterDIName = "telegram_adapter_di"
)

func buildAdapters() error {
	def := di.Def{
		Name:  TelegramAdapterDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get(ConfigDIName).(*config.Config)
			client := telegram.NewClient(cfg.TelegramConfig.BotToken, cfg.TelegramConfig.ChatID)
			return adapter.NewTelegramAdapter(client), nil
		},
	}
	return builder.Add(def)
}
