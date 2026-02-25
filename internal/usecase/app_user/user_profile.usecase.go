package appuser

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-structure/global"
	dto "go-structure/internal/dto/app_user"
	"go-structure/internal/helper/database"
	"go-structure/internal/repository"
	appuserrepo "go-structure/internal/repository/app_user"
	"go-structure/internal/repository/model"
	appusertransformer "go-structure/internal/transformer/app_user"
	"go-structure/internal/usecase"
	"go-structure/internal/utils/generate"
	jwtutil "go-structure/internal/utils/jwt"
	"go-structure/pkg/validator"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	IUserProfileUsecase interface {
		RegisterUserProfile(ctx context.Context, req *dto.UserRegisterRequestDto) (*dto.UserRegisterResponseDto, error)
		ActiveUserProfile(ctx context.Context, phone string, code string, ip string, userAgent string) (bool, error)
		LoginUserProfile(ctx context.Context, req *dto.UserLoginRequestDto, ip string, userAgent string) (*dto.UserLoginResponseDto, error)
		GetUserProfile(ctx context.Context, accountID uuid.UUID) (*dto.UserProfileResponseDto, error)
		RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshTokenResponseDto, error)
		Logout(ctx context.Context, accountID uuid.UUID) error
		UpdateUserProfile(ctx context.Context, accountID uuid.UUID, req *dto.UpdateUserProfile) (*dto.UpdateUserProfileResponseDto, error)
		ForgotPassword(ctx context.Context, phone string) (*dto.ForgotPasswordResponseDto, error)
		ResetPassword(ctx context.Context, phone string, code string, newPassword string, ip string, userAgent string) (*dto.ResetPasswordResponseDto, error)
	}

	userProfileUsecase struct {
		userProfileRepository appuserrepo.IUserProfileRepository
		accountRepository     repository.IAccountRepository
		deviceRepository      repository.IDeviceRepository
		accountAppDeviceRepo  repository.IAccountAppDeviceRepository
		sessionRepository     repository.ISessionRepository
		loginHistoryRepo      repository.ILoginHistoryRepository
		notifyUsecase         usecase.INotifyUsecase
		otpUsecase            usecase.IOTPUsecase
		txManager             database.TransactionManager
	}
)

var (
	ErrUserAlreadyRegistered = errors.New("tài khoản đã được đăng ký")
	ErrInvalidOTP            = errors.New("mã OTP không hợp lệ")
	ErrUserNotFound          = errors.New("không tìm thấy người dùng")
	ErrInvalidPassword       = errors.New("mật khẩu không chính xác")
	ErrUserNotActive         = errors.New("tài khoản chưa được kích hoạt")
	ErrUserAlreadyActive     = errors.New("tài khoản đã được kích hoạt")
	ErrInvalidRefreshToken   = errors.New("refresh token không hợp lệ hoặc đã hết hạn")
)

const (
	refreshTokenTTL = 15 * 24 * time.Hour
	appTypeUser     = "user"
)

func NewUserProfileUsecase(
	userProfileRepo appuserrepo.IUserProfileRepository,
	accountRepo repository.IAccountRepository,
	deviceRepo repository.IDeviceRepository,
	accountAppDeviceRepo repository.IAccountAppDeviceRepository,
	sessionRepo repository.ISessionRepository,
	loginHistoryRepo repository.ILoginHistoryRepository,
	notifyUc usecase.INotifyUsecase,
	otpUc usecase.IOTPUsecase,
	txManager database.TransactionManager,
) IUserProfileUsecase {
	return &userProfileUsecase{
		userProfileRepository: userProfileRepo,
		accountRepository:     accountRepo,
		deviceRepository:      deviceRepo,
		accountAppDeviceRepo:  accountAppDeviceRepo,
		sessionRepository:     sessionRepo,
		loginHistoryRepo:      loginHistoryRepo,
		notifyUsecase:         notifyUc,
		otpUsecase:            otpUc,
		txManager:             txManager,
	}
}

func (u *userProfileUsecase) RegisterUserProfile(ctx context.Context, req *dto.UserRegisterRequestDto) (*dto.UserRegisterResponseDto, error) {
	otpCode, err := database.WithTransaction(
		u.txManager,
		ctx,
		func(txCtx context.Context) (string, error) {
			account, err := u.getOrCreateAccount(txCtx, req)
			if err != nil {
				return "", err
			}
			if err := u.ensureUserNotAlreadyRegistered(txCtx, account.ID); err != nil {
				return "", err
			}
			if err := u.createUserProfileRecord(txCtx, account.ID, req.FullName); err != nil {
				return "", err
			}
			if u.notifyUsecase != nil {
				code, err := u.otpUsecase.CreateOTP(txCtx, req.Phone, usecase.OTPPurposeUserRegister)
				if err != nil {
					return "", err
				}
				return code, nil
			}
			return "", nil
		},
	)
	if err != nil {
		return nil, err
	}

	if u.notifyUsecase != nil && otpCode != "" {
		msg := fmt.Sprintf("Đăng ký tài khoản thành công, vui lòng kiểm tra điện thoại để nhận mã OTP: %s", otpCode)
		_ = u.notifyUsecase.SendOtp(ctx, msg)
	}

	return &dto.UserRegisterResponseDto{
		UserMessage: "Đăng ký tài khoản thành công, vui lòng kiểm tra điện thoại để nhận mã OTP",
	}, nil
}

func (u *userProfileUsecase) getOrCreateAccount(ctx context.Context, req *dto.UserRegisterRequestDto) (*model.Account, error) {
	hashedPassword, err := validator.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	account, err := u.accountRepository.GetByPhone(ctx, req.Phone)
	if err != nil {
		return nil, err
	}

	if account != nil {
		return account, nil
	}

	account, err = u.accountRepository.CreateAccount(ctx, &model.Account{
		Phone:        req.Phone,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		global.Logger.Logger.Info("create account failed", zap.String("phone", req.Phone), zap.Error(err))
		return nil, err
	}
	return account, nil
}

func (u *userProfileUsecase) ensureUserNotAlreadyRegistered(ctx context.Context, accountID uuid.UUID) error {
	existingProfile, err := u.userProfileRepository.GetByAccountID(ctx, accountID)
	if err != nil {
		return err
	}
	if existingProfile != nil {
		return ErrUserAlreadyRegistered
	}
	return nil
}

func (u *userProfileUsecase) createUserProfileRecord(ctx context.Context, accountID uuid.UUID, fullName string) error {
	userProfile := &model.UserProfile{
		AccountID: accountID,
		FullName:  fullName,
		IsActive:  false,
	}
	_, err := u.userProfileRepository.RegisterUserProfile(ctx, userProfile)
	if err != nil {
		global.Logger.Info("register user profile failed", zap.String("accountID", accountID.String()), zap.Error(err))
		return err
	}
	return nil
}

func (u *userProfileUsecase) ActiveUserProfile(ctx context.Context, phone string, code string, ip string, userAgent string) (bool, error) {
	account, err := u.accountRepository.GetByPhone(ctx, phone)
	if err != nil {
		return false, err
	}
	if account == nil {
		return false, ErrUserNotFound
	}

	userProfile, err := u.userProfileRepository.GetByAccountID(ctx, account.ID)
	if err != nil {
		return false, err
	}
	if userProfile == nil {
		return false, ErrUserNotFound
	}
	if userProfile.IsActive {
		return false, ErrUserAlreadyActive
	}

	verified, err := u.otpUsecase.VerifyOTP(ctx, phone, code, usecase.OTPPurposeUserRegister, ip, userAgent)
	if err != nil {
		return false, err
	}
	if !verified {
		return false, ErrInvalidOTP
	}

	_, err = database.WithTransaction(
		u.txManager,
		ctx,
		func(txCtx context.Context) (struct{}, error) {
			profileInTx, err := u.userProfileRepository.GetByAccountID(txCtx, account.ID)
			if err != nil {
				return struct{}{}, err
			}
			if profileInTx == nil {
				return struct{}{}, ErrUserNotFound
			}
			if profileInTx.IsActive {
				return struct{}{}, ErrUserAlreadyActive
			}
			profileInTx.IsActive = true
			profileInTx.UpdatedAt = time.Now()
			_, err = u.userProfileRepository.UpdateUserProfile(txCtx, profileInTx)
			return struct{}{}, err
		},
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *userProfileUsecase) LoginUserProfile(ctx context.Context, req *dto.UserLoginRequestDto, ip string, userAgent string) (*dto.UserLoginResponseDto, error) {
	device, err := u.ensureDevice(ctx, req.Device)
	global.GetChannelLogger("auth").Info("ensureDevice", zap.String("device", device.DeviceUID), zap.Error(err))
	if err != nil {
		return nil, err
	}

	type txResult struct {
		account      *model.Account
		userProfile  *model.UserProfile
		refreshToken string
	}

	res, err := database.WithTransaction(
		u.txManager,
		ctx,
		func(txCtx context.Context) (txResult, error) {
			account, userProfile, err := u.getAccountAndProfileForLogin(txCtx, req.Phone, req.Password)
			if err != nil {
				_ = u.logLoginHistory(ctx, uuid.Nil, device.ID, appTypeUser, ip, userAgent, err, req.Phone)
				return txResult{}, err
			}

			accountAppDevice, err := u.ensureAccountAppDevice(txCtx, account, device, appTypeUser, req.Device.FCMToken)
			if err != nil {
				return txResult{}, err
			}

			refreshToken, err := u.createLoginSession(txCtx, accountAppDevice.ID, ip, userAgent)
			if err != nil {
				return txResult{}, err
			}

			if err := u.logLoginHistory(txCtx, account.ID, device.ID, appTypeUser, ip, userAgent, nil, ""); err != nil {
				return txResult{}, err
			}

			return txResult{account: account, userProfile: userProfile, refreshToken: refreshToken}, nil
		},
	)
	if err != nil {
		if errors.Is(err, ErrUserNotActive) {
			return &dto.UserLoginResponseDto{
				RequireVerifyOtp: true,
				Message:          "vui lòng xác thực OTP trước khi đăng nhập",
			}, nil
		}
		return nil, err
	}

	accessToken, err := jwtutil.GenerateAccessToken(res.account.ID)
	if err != nil {
		return nil, err
	}

	resp := appusertransformer.ToLoginResponseDto(accessToken, res.refreshToken, res.account, res.userProfile)
	return &resp, nil
}

func (u *userProfileUsecase) getAccountAndProfileForLogin(ctx context.Context, phone string, password string) (*model.Account, *model.UserProfile, error) {
	account, err := u.accountRepository.GetByPhone(ctx, phone)
	if err != nil {
		return nil, nil, err
	}
	if account == nil {
		return nil, nil, ErrUserNotFound
	}

	if ok := validator.CheckPassword(password, account.PasswordHash); !ok {
		return nil, nil, ErrInvalidPassword
	}

	userProfile, err := u.userProfileRepository.GetByAccountID(ctx, account.ID)
	if err != nil {
		return nil, nil, err
	}
	if userProfile == nil {
		return nil, nil, ErrUserNotFound
	}
	if !userProfile.IsActive {
		return nil, nil, ErrUserNotActive
	}

	return account, userProfile, nil
}

func (u *userProfileUsecase) ensureDevice(ctx context.Context, deviceReq dto.Device) (*model.Device, error) {
	device, err := u.deviceRepository.GetDeviceByUID(ctx, deviceReq.DeviceUID)
	if err != nil {
		return nil, err
	}
	if device != nil {
		return device, nil
	}

	return u.deviceRepository.CreateDevice(ctx, &model.Device{
		DeviceUID:  deviceReq.DeviceUID,
		Platform:   deviceReq.Platform,
		DeviceName: deviceReq.DeviceName,
		OsVersion:  deviceReq.OsVersion,
		AppVersion: deviceReq.AppVersion,
	})
}

func (u *userProfileUsecase) ensureAccountAppDevice(
	ctx context.Context,
	account *model.Account,
	device *model.Device,
	appType string,
	fcmToken string,
) (*model.AccountAppDevice, error) {
	accountAppDevice, err := u.accountAppDeviceRepo.GetByAccountDeviceAndAppType(ctx, account.ID, device.ID, appType)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	if accountAppDevice == nil {
		return u.accountAppDeviceRepo.CreateAccountAppDevice(ctx, &model.AccountAppDevice{
			AccountID:  account.ID,
			DeviceID:   device.ID,
			AppType:    appType,
			FcmToken:   fcmToken,
			IsActive:   true,
			LastUsedAt: &now,
		})
	}

	accountAppDevice.IsActive = true
	accountAppDevice.LastUsedAt = &now
	if fcmToken != "" {
		accountAppDevice.FcmToken = fcmToken
	}

	updated, err := u.accountAppDeviceRepo.UpdateAccountAppDevice(ctx, accountAppDevice)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (u *userProfileUsecase) createLoginSession(
	ctx context.Context,
	accountAppDeviceID uuid.UUID,
	ip string,
	userAgent string,
) (string, error) {
	refreshToken := generate.GenerateRandomString(64)
	expiresAt := time.Now().Add(refreshTokenTTL)

	_, err := u.sessionRepository.CreateSession(ctx, &model.Session{
		AccountAppDeviceID: accountAppDeviceID,
		RefreshTokenHash:   refreshToken,
		ExpiresAt:          expiresAt,
		IpAddress:          ip,
		UserAgent:          userAgent,
	})
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

func (u *userProfileUsecase) logLoginHistory(
	ctx context.Context,
	accountID uuid.UUID,
	deviceID uuid.UUID,
	appType string,
	ip string,
	userAgent string,
	loginErr error,
	phone string,
) error {
	var result string
	var reason string
	if loginErr == nil {
		result = "success"
	} else {
		if accountID == uuid.Nil && phone != "" {
			account, _ := u.accountRepository.GetByPhone(ctx, phone)
			if account != nil {
				accountID = account.ID
			}
		}
		switch {
		case errors.Is(loginErr, ErrUserNotFound):
			result = "failed_not_found"
			reason = "Tài khoản không tồn tại"
		case errors.Is(loginErr, ErrInvalidPassword):
			result = "failed_password"
			reason = "Mật khẩu không chính xác"
		case errors.Is(loginErr, ErrUserNotActive):
			result = "failed_inactive"
			reason = "Tài khoản chưa được kích hoạt"
		default:
			result = "failed_unknown"
			reason = loginErr.Error()
		}
	}
	_, err := u.loginHistoryRepo.CreateLoginHistory(ctx, &model.AppLoginHistory{
		AccountID:     accountID,
		DeviceID:      deviceID,
		AppType:       appType,
		Result:        result,
		FailureReason: reason,
		IpAddress:     ip,
		UserAgent:     userAgent,
	})
	return err
}

func (u *userProfileUsecase) GetUserProfile(ctx context.Context, accountID uuid.UUID) (*dto.UserProfileResponseDto, error) {
	userProfile, err := u.userProfileRepository.GetByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if userProfile == nil {
		return nil, ErrUserNotFound
	}

	account, err := u.accountRepository.GetById(ctx, accountID.String())
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, ErrUserNotFound
	}

	resp := appusertransformer.ToUserProfileResponse(account, userProfile)
	return &resp, nil
}

func (u *userProfileUsecase) RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshTokenResponseDto, error) {
	type txResult struct {
		account     *model.Account
		userProfile *model.UserProfile
		newRefresh  string
		newAccess   string
	}

	res, err := database.WithTransaction(
		u.txManager,
		ctx,
		func(txCtx context.Context) (txResult, error) {
			session, err := u.sessionRepository.GetSessionByRefreshToken(txCtx, refreshToken)
			if err != nil || session == nil {
				return txResult{}, ErrInvalidRefreshToken
			}

			now := time.Now()
			if session.IsRevoked || now.After(session.ExpiresAt) {
				_ = u.sessionRepository.RevokeSessionByRefreshToken(txCtx, refreshToken, "expired_or_revoked")
				return txResult{}, ErrInvalidRefreshToken
			}

			accountAppDevice, err := u.accountAppDeviceRepo.GetByID(txCtx, session.AccountAppDeviceID)
			if err != nil || accountAppDevice == nil {
				return txResult{}, ErrInvalidRefreshToken
			}

			userProfile, err := u.userProfileRepository.GetByAccountID(txCtx, accountAppDevice.AccountID)
			if err != nil || userProfile == nil {
				return txResult{}, ErrUserNotFound
			}
			if !userProfile.IsActive {
				return txResult{}, ErrUserNotActive
			}

			account, err := u.accountRepository.GetById(txCtx, accountAppDevice.AccountID.String())
			if err != nil || account == nil {
				return txResult{}, ErrUserNotFound
			}

			if err := u.sessionRepository.RevokeSessionByRefreshToken(txCtx, refreshToken, "rotated"); err != nil {
				return txResult{}, err
			}

			newRefresh := generate.GenerateRandomString(64)
			_, err = u.sessionRepository.CreateSession(txCtx, &model.Session{
				AccountAppDeviceID: accountAppDevice.ID,
				RefreshTokenHash:   newRefresh,
				ExpiresAt:          now.Add(refreshTokenTTL),
				IpAddress:          session.IpAddress,
				UserAgent:          session.UserAgent,
			})
			if err != nil {
				return txResult{}, err
			}

			newAccess, err := jwtutil.GenerateAccessToken(account.ID)
			if err != nil {
				return txResult{}, err
			}

			return txResult{account: account, userProfile: userProfile, newRefresh: newRefresh, newAccess: newAccess}, nil
		},
	)
	if err != nil {
		return nil, err
	}

	resp := appusertransformer.ToRefreshTokenResponseDto(res.newAccess, res.newRefresh)
	return &resp, nil
}

func (u *userProfileUsecase) Logout(ctx context.Context, accountID uuid.UUID) error {
	return u.sessionRepository.RevokeAllSessionsByAccount(ctx, accountID, "logout")
}

func (u *userProfileUsecase) UpdateUserProfile(ctx context.Context, accountID uuid.UUID, req *dto.UpdateUserProfile) (*dto.UpdateUserProfileResponseDto, error) {
	type txResult struct {
		account        *model.Account
		updatedProfile *model.UserProfile
	}

	res, err := database.WithTransaction(
		u.txManager,
		ctx,
		func(txCtx context.Context) (txResult, error) {
			userProfile, err := u.userProfileRepository.GetByAccountID(txCtx, accountID)
			if err != nil {
				return txResult{}, err
			}
			if userProfile == nil {
				return txResult{}, ErrUserNotFound
			}

			account, err := u.accountRepository.GetById(txCtx, accountID.String())
			if err != nil {
				return txResult{}, err
			}
			if account == nil {
				return txResult{}, ErrUserNotFound
			}

			userProfile.FullName = req.FullName
			userProfile.AvatarURL = req.AvatarURL
			userProfile.UpdatedAt = time.Now()

			updatedProfile, err := u.userProfileRepository.UpdateUserProfile(txCtx, userProfile)
			if err != nil {
				return txResult{}, err
			}

			return txResult{account: account, updatedProfile: updatedProfile}, nil
		},
	)
	if err != nil {
		return nil, err
	}

	resp := appusertransformer.ToUpdateUserProfileResponseDto("Cập nhật thông tin thành công", res.account, res.updatedProfile)
	return &resp, nil
}

func (u *userProfileUsecase) ForgotPassword(ctx context.Context, phone string) (*dto.ForgotPasswordResponseDto, error) {
	account, err := u.accountRepository.GetByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, ErrUserNotFound
	}

	code, err := u.otpUsecase.CreateForgotPasswordOTP(ctx, phone)
	if err != nil {
		return nil, err
	}

	if u.notifyUsecase != nil && code != "" {
		msg := fmt.Sprintf("Mã OTP để đặt lại mật khẩu của bạn là: %s", code)
		_ = u.notifyUsecase.SendOtp(ctx, msg)
	}

	return &dto.ForgotPasswordResponseDto{
		UserMessage: "Vui lòng kiểm tra điện thoại để nhận mã OTP đặt lại mật khẩu",
	}, nil
}

func (u *userProfileUsecase) ResetPassword(ctx context.Context, phone string, code string, newPassword string, ip string, userAgent string) (*dto.ResetPasswordResponseDto, error) {
	account, err := u.accountRepository.GetByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, ErrUserNotFound
	}

	verified, err := u.otpUsecase.VerifyOTP(ctx, phone, code, usecase.OTPPurposeResetPassword, ip, userAgent)
	if err != nil {
		return nil, err
	}
	if !verified {
		return nil, ErrInvalidOTP
	}

	hashedPassword, err := validator.HashPassword(newPassword)
	if err != nil {
		return nil, err
	}

	_, err = database.WithTransaction(
		u.txManager,
		ctx,
		func(txCtx context.Context) (struct{}, error) {
			if err := u.accountRepository.UpdatePassword(txCtx, account.ID, hashedPassword); err != nil {
				return struct{}{}, err
			}
			return struct{}{}, u.sessionRepository.RevokeAllSessionsByAccount(txCtx, account.ID, "password_reset")
		},
	)
	if err != nil {
		return nil, err
	}

	return &dto.ResetPasswordResponseDto{
		UserMessage: "Đặt lại mật khẩu thành công",
	}, nil
}
