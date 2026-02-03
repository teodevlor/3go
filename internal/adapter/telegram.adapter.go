package adapter

import (
	"context"

	"go-structure/pkg/telegram"
)

type (
	INotifyAdapter interface {
		SendMessage(ctx context.Context, message string) error
	}

	notifyAdapter struct {
		telegramClient *telegram.Client
	}
)

func NewTelegramAdapter(telegramClient *telegram.Client) INotifyAdapter {
	return &notifyAdapter{telegramClient: telegramClient}
}

func (a *notifyAdapter) SendMessage(ctx context.Context, message string) error {
	return a.telegramClient.SendMessage(ctx, message)
}
