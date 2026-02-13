package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-structure/internal/common"
	"go-structure/internal/dto"
	"go-structure/internal/helper/database"
	"go-structure/internal/repository"
	"go-structure/internal/setting"
	"go-structure/internal/utils/generate"

	"github.com/google/uuid"
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
		settingUsecase     ISettingUsecase
		notifyUsecase      INotifyUsecase
		txManager          database.TransactionManager
	}
)

const (
	OTPPurposeRegister       = "register"
	OTPPurposeResetPassword  = "reset_password"
	OTPLength                = 6
	OTPMaxAttempt            = 5
	OTPExpireMinutes         = 5
	FailureReasonInvalidCode = "invalid_code"
	FailureReasonMaxAttempt  = "max_attempt"
	ResultSuccess            = "success"
	ResultFailed             = "failed"
	ResultLocked             = "locked"
)

func NewOTPUsecase(
	otpRepository repository.IOTPRepository,
	otpAuditRepository repository.IOTPAuditRepository,
	settingUsecase ISettingUsecase,
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
	purpose = u.defaultPurposeIfEmpty(purpose)
	return u.createOTPWithPurpose(ctx, target, purpose, nil)
}

func (u *otpUsecase) CreateForgotPasswordOTP(ctx context.Context, target string) (string, error) {
	purpose := OTPPurposeResetPassword

	cfg, err := u.settingUsecase.GetResendConfig(ctx)
	if err != nil {
		return "", err
	}

	if err := u.ensureResendCooldown(ctx, target, purpose, cfg); err != nil {
		return "", err
	}
	if err := u.ensureResendRateLimit(ctx, target, purpose, cfg); err != nil {
		return "", err
	}

	return u.createOTPWithPurpose(ctx, target, purpose, cfg)
}

func (u *otpUsecase) createOTPWithPurpose(ctx context.Context, target string, purpose string, cfg *setting.ResendOTPConfig) (string, error) {
	code := generate.GenerateOTPCode(OTPLength)
	if err := u.otpRepository.ExpireOldOTPs(ctx); err != nil {
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
		return "", err
	}
	return code, nil
}

func (u *otpUsecase) ResendOTP(ctx context.Context, target string, purpose string) (string, error) {
	purpose = u.defaultPurposeIfEmpty(purpose)

	cfg, err := u.settingUsecase.GetResendConfig(ctx)
	if err != nil {
		return "", err
	}

	if purpose != OTPPurposeRegister {
		if err := u.ensureHasPreviousOTP(ctx, target, purpose); err != nil {
			return "", err
		}
	}

	if err := u.ensureResendCooldown(ctx, target, purpose, cfg); err != nil {
		return "", err
	}
	if err := u.ensureResendRateLimit(ctx, target, purpose, cfg); err != nil {
		return "", err
	}

	code, err := u.createOTPWithPurpose(ctx, target, purpose, cfg)
	if err != nil {
		return "", err
	}

	if u.notifyUsecase != nil {
		msg := fmt.Sprintf("Mã OTP cho mục đích %s của bạn là: %s", purpose, code)
		_ = u.notifyUsecase.SendOtp(ctx, msg)
	}

	return code, nil
}

func (u *otpUsecase) VerifyOTP(ctx context.Context, target string, code string, purpose string, ip string, userAgent string) (bool, error) {
	var verified bool

	err := u.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		row, err := u.fetchActiveOTP(txCtx, target, purpose)
		if err != nil {
			return err
		}

		if row == nil {
			verified = false
			return nil
		}

		attemptNumber := int(row.AttemptCount) + 1
		if !u.isCodeMatched(row, code) {
			verified = false
			return u.handleInvalidOTP(txCtx, row, attemptNumber, target, purpose, ip, userAgent)
		}

		verified = true
		return u.handleValidOTP(txCtx, row, attemptNumber, target, purpose, ip, userAgent)
	})

	if err != nil {
		return false, err
	}

	return verified, nil
}

// Helpers
func (u *otpUsecase) resolveExpireDuration(cfg *setting.ResendOTPConfig) time.Duration {
	expireSec := int64(OTPExpireMinutes * 60)
	if cfg != nil && cfg.TimeOutExpired > 0 {
		expireSec = int64(cfg.TimeOutExpired)
	}
	return time.Duration(expireSec) * time.Second
}

// Dùng để chặn spam resend OTP cho những mục tiêu chưa từng yêu cầu OTP trước đó
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
		return OTPPurposeRegister
	}

	switch purpose {
	case "register":
		return OTPPurposeRegister
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
		return fmt.Errorf("%w"+common.BaseMessageResendOTPTooSoonWaitSeconds, common.ErrResendTooSoon, waitSec)
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
