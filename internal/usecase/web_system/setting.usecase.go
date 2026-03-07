package web_system

import (
	"context"
	"encoding/json"

	"go-structure/global"
	settingRepository "go-structure/internal/repository/web_system"
	"go-structure/internal/setting"
	"go-structure/internal/middleware"

	"go.uber.org/zap"
)

const (
	SettingKeyResendOTP = "config_resent_otp"
)

type (
	ISettingUsecase interface {
		GetResendConfig(ctx context.Context) (*setting.ResendOTPConfig, error)
	}

	settingUsecase struct {
		settingRepository settingRepository.ISettingRepository
	}
)

func NewSettingUsecase(settingRepository settingRepository.ISettingRepository) ISettingUsecase {
	return &settingUsecase{settingRepository: settingRepository}
}

func defaultResendConfig() *setting.ResendOTPConfig {
	return &setting.ResendOTPConfig{
		MaxCount:       3,
		TimeOutExpired: 300,
		TimeOutResent:  60,
		TrackingTTL:    900,
		BlockDurations: setting.ResendOTPBlockDurs{
			Violation1:     300,
			Violation2:     900,
			Violation3:     3600,
			Violation4Plus: 86400,
		},
	}
}

func (u *settingUsecase) GetResendConfig(ctx context.Context) (*setting.ResendOTPConfig, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("GetResendConfig: start", zap.String(global.KeyCorrelationID, cid))

	raw := u.GetSettingValueByKey(ctx, SettingKeyResendOTP)
	if len(raw) == 0 {
		global.Logger.Info("GetResendConfig: completed successfully (using default)", zap.String(global.KeyCorrelationID, cid))
		return defaultResendConfig(), nil
	}
	var cfg setting.ResendOTPConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		global.Logger.Error("GetResendConfig: failed to unmarshal config", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		global.Logger.Info("GetResendConfig: completed (using default due to parse error)", zap.String(global.KeyCorrelationID, cid))
		return defaultResendConfig(), nil
	}
	global.Logger.Info("GetResendConfig: completed successfully", zap.String(global.KeyCorrelationID, cid))
	return &cfg, nil
}

func (u *settingUsecase) GetSettingValueByKey(ctx context.Context, key string) []byte {
	if u.settingRepository == nil {
		return nil
	}
	s, err := u.settingRepository.GetSettingByKey(ctx, key)
	if err != nil || s == nil || len(s.Value) == 0 {
		return nil
	}
	return s.Value
}
