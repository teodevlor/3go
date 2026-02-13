package usecase

import (
	"context"

	"go-structure/internal/setting"
)

type ISettingUsecase interface {
	GetResendConfig(ctx context.Context) (*setting.ResendOTPConfig, error)
}
