package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-structure/global"
	"go-structure/internal/common"
	"go-structure/internal/constants"
	"go-structure/internal/dto"
	"go-structure/internal/helper/database"
	"go-structure/internal/middleware"
	"go-structure/internal/repository"
	"go-structure/internal/setting"
	"go-structure/internal/usecase/web_system"
	"go-structure/internal/utils/generate"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	IOTPUsecase interface {
		CreateOTP(ctx context.Context, target string, purpose string) (string, error)
		VerifyOTP(ctx context.Context, target string, code string, purpose string, ip string, userAgent string) (bool, error)
		ResendOTP(ctx context.Context, target string, purpose string) (string, error)
		CreateForgotPasswordOTP(ctx context.Context, target string) (string, error)
	}

	otpUsecase struct {
		otpRepository      repository.IOTPRepository
		otpAuditRepository repository.IOTPAuditRepository
		settingUsecase     web_system.ISettingUsecase
		notifyUsecase      INotifyUsecase
		txManager          database.TransactionManager
	}
)

const (
	OTPPurposeUserRegister   = common.OTPPurposeUserRegister
	OTPPurposeDriverRegister = common.OTPPurposeDriverRegister
	OTPPurposeResetPassword  = common.OTPPurposeResetPassword
	OTPLength                = 6
	OTPMaxAttempt            = 5
	OTPExpireMinutes         = 5
	FailureReasonInvalidCode = common.OTPFailureReasonInvalidCode
	FailureReasonMaxAttempt  = common.OTPFailureReasonMaxAttempt
	ResultSuccess            = common.OTPResultSuccess
	ResultFailed             = common.OTPResultFailed
	ResultLocked             = common.OTPResultLocked
)

func NewOTPUsecase(
	otpRepository repository.IOTPRepository,
	otpAuditRepository repository.IOTPAuditRepository,
	settingUsecase web_system.ISettingUsecase,
	notifyUsecase INotifyUsecase,
	txManager database.TransactionManager,
) IOTPUsecase {
	return &otpUsecase{
		otpRepository:      otpRepository,
		otpAuditRepository: otpAuditRepository,
		settingUsecase:     settingUsecase,
		notifyUsecase:      notifyUsecase,
		txManager:          txManager,
	}
}

func (u *otpUsecase) CreateOTP(ctx context.Context, target string, purpose string) (string, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("CreateOTP: start", zap.String(global.KeyCorrelationID, cid), zap.String("target", target), zap.String("purpose", purpose))

	purpose = u.defaultPurposeIfEmpty(purpose)
	code, err := u.createOTPWithPurpose(ctx, target, purpose, nil)
	if err != nil {
		global.Logger.Error("CreateOTP: failed to create OTP", zap.String(global.KeyCorrelationID, cid), zap.String("target", target), zap.Error(err))
		return "", err
	}
	global.Logger.Info("CreateOTP: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("target", target))
	return code, nil
}

func (u *otpUsecase) CreateForgotPasswordOTP(ctx context.Context, target string) (string, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("CreateForgotPasswordOTP: start", zap.String(global.KeyCorrelationID, cid), zap.String("target", target))

	purpose := OTPPurposeResetPassword

	cfg, err := u.settingUsecase.GetResendConfig(ctx)
	if err != nil {
		global.Logger.Error("CreateForgotPasswordOTP: failed to get resend config", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return "", err
	}

	if err := u.ensureResendCooldown(ctx, target, purpose, cfg); err != nil {
		global.Logger.Error("CreateForgotPasswordOTP: resend cooldown check failed", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return "", err
	}
	if err := u.ensureResendRateLimit(ctx, target, purpose, cfg); err != nil {
		global.Logger.Error("CreateForgotPasswordOTP: resend rate limit exceeded", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return "", err
	}

	code, err := u.createOTPWithPurpose(ctx, target, purpose, cfg)
	if err != nil {
		global.Logger.Error("CreateForgotPasswordOTP: failed to create OTP", zap.String(global.KeyCorrelationID, cid), zap.String("target", target), zap.Error(err))
		return "", err
	}
	global.Logger.Info("CreateForgotPasswordOTP: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("target", target))
	return code, nil
}

func (u *otpUsecase) createOTPWithPurpose(ctx context.Context, target string, purpose string, cfg *setting.ResendOTPConfig) (string, error) {
	cid := middleware.CorrelationIDFromContext(ctx)

	code := generate.GenerateOTPCode(OTPLength)
	if err := u.otpRepository.ExpireOldOTPs(ctx); err != nil {
		global.Logger.Error("createOTPWithPurpose: failed to expire old OTPs", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return "", err
	}

	expiresAt := time.Now().Add(u.resolveExpireDuration(cfg))
	input := dto.CreateOTPRequestData{
		Target:     target,
		OtpCode:    code,
		Purpose:    purpose,
		MaxAttempt: OTPMaxAttempt,
		ExpiresAt:  expiresAt,
	}
	if err := u.otpRepository.CreateOTP(ctx, input); err != nil {
		global.Logger.Error("createOTPWithPurpose: failed to create OTP", zap.String(global.KeyCorrelationID, cid), zap.String("target", target), zap.Error(err))
		return "", err
	}
	return code, nil
}

func (u *otpUsecase) ResendOTP(ctx context.Context, target string, purpose string) (string, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("ResendOTP: start", zap.String(global.KeyCorrelationID, cid), zap.String("target", target), zap.String("purpose", purpose))

	purpose = u.defaultPurposeIfEmpty(purpose)

	cfg, err := u.settingUsecase.GetResendConfig(ctx)
	if err != nil {
		global.Logger.Error("ResendOTP: failed to get resend config", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return "", err
	}

	if purpose != OTPPurposeUserRegister {
		if err := u.ensureHasPreviousOTP(ctx, target, purpose); err != nil {
			global.Logger.Error("ResendOTP: no previous OTP to resend", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
			return "", err
		}
	}

	if err := u.ensureResendCooldown(ctx, target, purpose, cfg); err != nil {
		global.Logger.Error("ResendOTP: resend cooldown check failed", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return "", err
	}
	if err := u.ensureResendRateLimit(ctx, target, purpose, cfg); err != nil {
		global.Logger.Error("ResendOTP: resend rate limit exceeded", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return "", err
	}

	code, err := u.createOTPWithPurpose(ctx, target, purpose, cfg)
	if err != nil {
		global.Logger.Error("ResendOTP: failed to create OTP", zap.String(global.KeyCorrelationID, cid), zap.String("target", target), zap.Error(err))
		return "", err
	}

	if u.notifyUsecase != nil {
		msg := fmt.Sprintf("Mã OTP cho mục đích %s của bạn là: %s", purpose, code)
		_ = u.notifyUsecase.SendOtp(ctx, msg)
	}

	global.Logger.Info("ResendOTP: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("target", target))
	return code, nil
}

func (u *otpUsecase) VerifyOTP(ctx context.Context, target string, code string, purpose string, ip string, userAgent string) (bool, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("VerifyOTP: start", zap.String(global.KeyCorrelationID, cid), zap.String("target", target), zap.String("purpose", purpose))

	ok, err := database.WithTransaction(
		u.txManager,
		ctx,
		func(txCtx context.Context) (bool, error) {
			row, err := u.fetchActiveOTP(txCtx, target, purpose)
			if err != nil {
				return false, err
			}
			if row == nil {
				return false, nil
			}

			attemptNumber := int(row.AttemptCount) + 1
			if !u.isCodeMatched(row, code) {
				return false, u.handleInvalidOTP(txCtx, row, attemptNumber, target, purpose, ip, userAgent)
			}

			return true, u.handleValidOTP(txCtx, row, attemptNumber, target, purpose, ip, userAgent)
		},
	)
	if err != nil {
		global.Logger.Error("VerifyOTP: failed to verify OTP", zap.String(global.KeyCorrelationID, cid), zap.String("target", target), zap.Error(err))
		return false, err
	}
	global.Logger.Info("VerifyOTP: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("target", target), zap.Bool("verified", ok))
	return ok, nil
}

// Helpers
func (u *otpUsecase) resolveExpireDuration(cfg *setting.ResendOTPConfig) time.Duration {
	expireSec := int64(OTPExpireMinutes * 60)
	if cfg != nil && cfg.TimeOutExpired > 0 {
		expireSec = int64(cfg.TimeOutExpired)
	}
	return time.Duration(expireSec) * time.Second
}

func (u *otpUsecase) ensureHasPreviousOTP(ctx context.Context, target, purpose string) error {
	lastCreated, err := u.otpRepository.GetLastOTPCreatedAt(ctx, target, purpose)
	if err != nil {
		return err
	}
	if lastCreated == nil {
		return fmt.Errorf("không có yêu cầu OTP hợp lệ để resend cho mục đích %s", purpose)
	}
	return nil
}

func (u *otpUsecase) defaultPurposeIfEmpty(purpose string) string {
	purpose = strings.TrimSpace(strings.ToLower(purpose))

	if purpose == "" {
		return OTPPurposeUserRegister
	}

	switch purpose {
	case "register":
		return OTPPurposeUserRegister
	case "forgot-password":
		return OTPPurposeResetPassword
	default:
		return purpose
	}
}

func (u *otpUsecase) ensureResendCooldown(ctx context.Context, target, purpose string, cfg *setting.ResendOTPConfig) error {
	lastCreated, err := u.otpRepository.GetLastOTPCreatedAt(ctx, target, purpose)
	if err != nil {
		return err
	}
	if lastCreated == nil {
		return nil
	}

	elapsed := time.Since(*lastCreated)
	if int32(elapsed.Seconds()) < cfg.TimeOutResent {
		waitSec := cfg.TimeOutResent - int32(elapsed.Seconds())
		return fmt.Errorf("%w"+constants.BaseMessageResendOTPTooSoonWaitSeconds, common.ErrResendTooSoon, waitSec)
	}
	return nil
}

func (u *otpUsecase) ensureResendRateLimit(ctx context.Context, target, purpose string, cfg *setting.ResendOTPConfig) error {
	if cfg.TrackingTTL <= 0 || cfg.MaxCount <= 0 {
		return nil
	}

	since := time.Now().Add(-time.Duration(cfg.TrackingTTL) * time.Second)
	count, err := u.otpRepository.CountOTPsCreatedSince(ctx, target, purpose, since)
	if err != nil {
		return err
	}
	if count < cfg.MaxCount {
		return nil
	}

	retryAfterSec := u.computeRetryAfter(ctx, target, purpose, cfg, since)
	return &common.ErrorWithRetryAfter{
		Err:               common.ErrResendMaxExceeded,
		RetryAfterSeconds: retryAfterSec,
	}
}

func (u *otpUsecase) computeRetryAfter(ctx context.Context, target, purpose string, cfg *setting.ResendOTPConfig, since time.Time) int64 {
	oldest, _ := u.otpRepository.GetOldestOTPCreatedAtSince(ctx, target, purpose, since)
	if oldest == nil {
		return 0
	}
	retryAfterSec := oldest.Unix() + cfg.TrackingTTL - time.Now().Unix()
	if retryAfterSec < 0 {
		return 0
	}
	return retryAfterSec
}

func (u *otpUsecase) fetchActiveOTP(ctx context.Context, target, purpose string) (*dto.ActiveOTPResponseData, error) {
	return u.otpRepository.GetOTP(ctx, target, purpose)
}

func (u *otpUsecase) isCodeMatched(row *dto.ActiveOTPResponseData, code string) bool {
	return row.OtpCode == code
}

func (u *otpUsecase) handleInvalidOTP(ctx context.Context, row *dto.ActiveOTPResponseData, attemptNumber int, target, purpose string, ip string, userAgent string) error {
	if err := u.otpRepository.IncrementOTPAttempt(ctx, row.ID); err != nil {
		return err
	}

	failureReason := FailureReasonInvalidCode
	result := ResultFailed

	if attemptNumber >= int(row.MaxAttempt) {
		if err := u.otpRepository.LockOTP(ctx, row.ID); err != nil {
			return err
		}
		failureReason = FailureReasonMaxAttempt
		result = ResultLocked
	}

	return u.logOTPAudit(ctx, row.ID, target, purpose, attemptNumber, result, failureReason, ip, userAgent)
}

func (u *otpUsecase) handleValidOTP(ctx context.Context, row *dto.ActiveOTPResponseData, attemptNumber int, target, purpose string, ip string, userAgent string) error {
	if err := u.otpRepository.MarkOTPAsUsed(ctx, row.ID); err != nil {
		return err
	}
	return u.logOTPAudit(ctx, row.ID, target, purpose, attemptNumber, ResultSuccess, "", ip, userAgent)
}

func (u *otpUsecase) logOTPAudit(ctx context.Context, otpID uuid.UUID, target, purpose string, attemptNumber int, result, failureReason string, ip string, userAgent string) error {
	return u.otpAuditRepository.CreateOTPAudit(ctx, dto.CreateOTPAuditRequestData{
		OTPId:         otpID,
		Target:        target,
		Purpose:       purpose,
		AttemptNumber: attemptNumber,
		Result:        result,
		FailureReason: failureReason,
		IPAddress:     ip,
		UserAgent:     userAgent,
	})
}
