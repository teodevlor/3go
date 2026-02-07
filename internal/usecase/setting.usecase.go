package usecase

import (
	"context"
	"encoding/json"

	"go-structure/internal/repository"
	"go-structure/internal/setting"
)

const (
	SettingKeyResendOTP = "config_resent_otp"
)

type (
	ISettingUsecase interface {
		GetResendConfig(ctx context.Context) (*setting.ResendOTPConfig, error)
	}

	settingUsecase struct {
		settingRepository repository.ISettingRepository
	}
)

func NewSettingUsecase(settingRepository repository.ISettingRepository) ISettingUsecase {
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
	raw := u.GetSettingValueByKey(ctx, SettingKeyResendOTP)
	if len(raw) == 0 {
		return defaultResendConfig(), nil
	}
	var cfg setting.ResendOTPConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return defaultResendConfig(), nil
	}
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
