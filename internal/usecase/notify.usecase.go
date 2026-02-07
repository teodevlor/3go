package usecase

import (
	"context"
	telegramadapter "go-structure/internal/adapter"
)

type (
	INotifyUsecase interface {
		SendOtp(ctx context.Context, message string) error
	}

	notifyUsecase struct {
		notifyAdapter telegramadapter.INotifyAdapter
	}
)

func NewNotifyUsecase(notifyAdapter telegramadapter.INotifyAdapter) INotifyUsecase {
	return &notifyUsecase{notifyAdapter: notifyAdapter}
}

func (u *notifyUsecase) SendOtp(ctx context.Context, message string) error {
	return u.notifyAdapter.SendOtp(ctx, message)
}
