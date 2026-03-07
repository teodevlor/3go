package usecase

import (
	"context"

	"go-structure/global"
	telegramadapter "go-structure/internal/adapter"
	"go-structure/internal/middleware"

	"go.uber.org/zap"
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
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("SendOtp: start", zap.String(global.KeyCorrelationID, cid))
	if err := u.notifyAdapter.SendOtp(ctx, message); err != nil {
		global.Logger.Error("SendOtp: failed to send OTP notification", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return err
	}
	global.Logger.Info("SendOtp: completed successfully", zap.String(global.KeyCorrelationID, cid))
	return nil
}
