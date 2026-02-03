package usecase

import (
	"context"
	"go-structure/internal/adapter"
)

type (
	INotifyUsecase interface {
		SendMessage(ctx context.Context, message string) error
	}

	notifyUsecase struct {
		notifyAdapter adapter.INotifyAdapter
	}
)

func NewNotifyUsecase(notifyAdapter adapter.INotifyAdapter) INotifyUsecase {
	return &notifyUsecase{notifyAdapter: notifyAdapter}
}

func (u *notifyUsecase) SendMessage(ctx context.Context, message string) error {
	return u.notifyAdapter.SendMessage(ctx, message)
}
