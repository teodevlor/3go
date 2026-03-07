package registry

import (
	"go-structure/config"
	"go-structure/internal/adapter"
	"go-structure/internal/adapter/storage"
	"go-structure/pkg/telegram"

	"github.com/sarulabs/di"
)

const (
	TelegramAdapterDIName = "telegram_adapter_di"
	StorageAdapterDIName  = "storage_adapter_di"
)

func buildAdapters() error {
	telegramDef := di.Def{
		Name:  TelegramAdapterDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get(ConfigDIName).(*config.Config)
			client := telegram.NewClient(cfg.TelegramConfig.BotToken, cfg.TelegramConfig.ChatID)
			return adapter.NewTelegramAdapter(client), nil
		},
	}
	storageDef := di.Def{
		Name:  StorageAdapterDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get(ConfigDIName).(*config.Config)
			s := cfg.Storage
			minioCfg := storage.MinIOConfig{
				Endpoint:      s.Endpoint,
				AccessKey:     s.AccessKey,
				SecretKey:     s.SecretKey,
				BucketPublic:  s.BucketPublic,
				BucketPrivate: s.BucketPrivate,
			}
			return storage.NewMinIOAdapter(minioCfg)
		},
	}
	return builder.Add(telegramDef, storageDef)
}
