package app_driver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	dto "go-structure/internal/dto/app_driver"
	"go-structure/internal/helper/database"
	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository"
	appdriverrepo "go-structure/internal/repository/app_driver"
	"go-structure/internal/repository/model"
	appdrivermodel "go-structure/internal/repository/model/app_driver"
	appdrivertransformer "go-structure/internal/transformer/app_driver"
	"go-structure/internal/usecase"
	"go-structure/internal/utils/generate"
	jwtutil "go-structure/internal/utils/jwt"
	"go-structure/pkg/validator"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	appTypeDriver   = "driver"
	refreshTokenTTL = 15 * 24 * time.Hour

	driverOnlineTTL = 15 * time.Second

	driverPingLimitWindow = 5 * time.Second
	driverPingLimitCount  = 3
)

var (
	ErrDriverAlreadyRegistered = errors.New("số điện thoại đã được đăng ký làm tài xế")
	ErrDriverNotFound          = errors.New("không tìm thấy hồ sơ tài xế")
	ErrDriverInvalidOTP        = errors.New("mã OTP không hợp lệ")
	ErrDriverAlreadyVerified   = errors.New("tài khoản tài xế đã được xác thực")
	ErrDriverRequireVerifyOtp  = errors.New("vui lòng xác thực OTP trước khi đăng nhập")
	ErrDriverInvalidPassword   = errors.New("mật khẩu không chính xác")
	ErrDriverNotActive         = errors.New("tài xế chưa được kích hoạt, vui lòng liên hệ CSKH")
	ErrDriverPingTooFrequent   = errors.New("bạn đang gửi yêu cầu quá nhanh, vui lòng thử lại sau")
)

type (
	IDriverProfileUsecase interface {
		RegisterDriver(ctx context.Context, req *dto.DriverRegisterRequestDto) (*dto.DriverRegisterResponseDto, error)
		VerifyDriverOtp(ctx context.Context, phone, code, ip, userAgent string) (*dto.DriverVerifyOtpResponseDto, error)
		LoginDriver(ctx context.Context, req *dto.DriverLoginRequestDto, ip, userAgent string) (*dto.DriverLoginResponseDto, error)

		GoOnline(ctx context.Context, accountID uuid.UUID, req *dto.DriverLocationStatusRequestDto) error
		GoOffline(ctx context.Context, accountID uuid.UUID) error
		PingOnline(ctx context.Context, accountID uuid.UUID, req *dto.DriverLocationStatusRequestDto) error
	}

	driverProfileUsecase struct {
		driverProfileRepo    appdriverrepo.IDriverProfileRepository
		accountRepo          repository.IAccountRepository
		deviceRepo           repository.IDeviceRepository
		accountAppDeviceRepo repository.IAccountAppDeviceRepository
		sessionRepo          repository.ISessionRepository
		loginHistoryRepo     repository.ILoginHistoryRepository
		otpUsecase           usecase.IOTPUsecase
		notifyUsecase        usecase.INotifyUsecase
		txManager            database.TransactionManager

		redisClient *redis.Client
	}
)

func NewDriverProfileUsecase(
	driverProfileRepo appdriverrepo.IDriverProfileRepository,
	accountRepo repository.IAccountRepository,
	deviceRepo repository.IDeviceRepository,
	accountAppDeviceRepo repository.IAccountAppDeviceRepository,
	sessionRepo repository.ISessionRepository,
	loginHistoryRepo repository.ILoginHistoryRepository,
	otpUsecase usecase.IOTPUsecase,
	notifyUsecase usecase.INotifyUsecase,
	txManager database.TransactionManager,
	redisClient *redis.Client,
) IDriverProfileUsecase {
	return &driverProfileUsecase{
		driverProfileRepo:    driverProfileRepo,
		accountRepo:          accountRepo,
		deviceRepo:           deviceRepo,
		accountAppDeviceRepo: accountAppDeviceRepo,
		sessionRepo:          sessionRepo,
		loginHistoryRepo:     loginHistoryRepo,
		otpUsecase:           otpUsecase,
		notifyUsecase:        notifyUsecase,
		txManager:            txManager,
		redisClient:          redisClient,
	}
}

type driverOnlineValue struct {
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
	Status string  `json:"status"`
}

func (u *driverProfileUsecase) RegisterDriver(ctx context.Context, req *dto.DriverRegisterRequestDto) (*dto.DriverRegisterResponseDto, error) {
	otpCode, err := database.WithTransaction(
		u.txManager,
		ctx,
		func(txCtx context.Context) (string, error) {
			account, err := u.getOrCreateAccount(txCtx, req)
			if err != nil {
				return "", err
			}
			if err := u.ensureDriverNotAlreadyRegistered(txCtx, account.ID); err != nil {
				return "", err
			}
			if err := u.createDriverProfileRecord(txCtx, account.ID, req.FullName); err != nil {
				return "", err
			}
			if u.otpUsecase != nil {
				code, err := u.otpUsecase.CreateOTP(txCtx, req.Phone, usecase.OTPPurposeDriverRegister)
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
		msg := fmt.Sprintf("Đăng ký tài khoản tài xế thành công. Mã OTP xác thực: %s", otpCode)
		_ = u.notifyUsecase.SendOtp(ctx, msg)
	}

	return &dto.DriverRegisterResponseDto{
		UserMessage: "Đăng ký tài khoản tài xế thành công, vui lòng kiểm tra điện thoại để nhận mã OTP",
	}, nil
}

func (u *driverProfileUsecase) getOrCreateAccount(ctx context.Context, req *dto.DriverRegisterRequestDto) (*model.Account, error) {
	acc, err := u.accountRepo.GetByPhone(ctx, req.Phone)
	if err != nil {
		return nil, err
	}
	if acc != nil {
		return acc, nil
	}
	hashedPassword, err := validator.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	account, err := u.accountRepo.CreateAccount(ctx, &model.Account{
		Phone:        req.Phone,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (u *driverProfileUsecase) ensureDriverNotAlreadyRegistered(ctx context.Context, accountID uuid.UUID) error {
	existing, err := u.driverProfileRepo.GetByAccountID(ctx, accountID)
	if err != nil {
		return err
	}
	if existing != nil {
		return ErrDriverAlreadyRegistered
	}
	return nil
}

func (u *driverProfileUsecase) createDriverProfileRecord(ctx context.Context, accountID uuid.UUID, fullName string) error {
	_, err := u.driverProfileRepo.Create(ctx, accountID, fullName)
	return err
}

func (u *driverProfileUsecase) VerifyDriverOtp(ctx context.Context, phone, code, ip, userAgent string) (*dto.DriverVerifyOtpResponseDto, error) {
	account, err := u.accountRepo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, ErrDriverNotFound
	}

	driverProfile, err := u.driverProfileRepo.GetByAccountID(ctx, account.ID)
	if err != nil {
		return nil, err
	}
	if driverProfile == nil {
		return nil, ErrDriverNotFound
	}

	verified, err := u.otpUsecase.VerifyOTP(ctx, phone, code, usecase.OTPPurposeDriverRegister, ip, userAgent)
	if err != nil {
		return nil, err
	}
	if !verified {
		return nil, ErrDriverInvalidOTP
	}

	fromStatus := pgdb.NullDriverProfileStatus{
		DriverProfileStatus: pgdb.DriverProfileStatus(driverProfile.GlobalStatus),
		Valid:               true,
	}
	_, err = database.WithTransaction(
		u.txManager,
		ctx,
		func(txCtx context.Context) (struct{}, error) {
			if err := u.driverProfileRepo.CreateStatusHistory(txCtx, driverProfile.ID, fromStatus, pgdb.DriverProfileStatusDOCUMENTINCOMPLETE, nil, nil); err != nil {
				return struct{}{}, err
			}
			_, err := u.driverProfileRepo.UpdateStatus(txCtx, driverProfile.ID, pgdb.DriverProfileStatusDOCUMENTINCOMPLETE)
			return struct{}{}, err
		},
	)
	if err != nil {
		return nil, err
	}

	return &dto.DriverVerifyOtpResponseDto{
		UserMessage: "Xác thực OTP thành công. Bạn có thể đăng nhập và upload tài liệu.",
		DriverID:    driverProfile.ID,
	}, nil
}

func (u *driverProfileUsecase) LoginDriver(ctx context.Context, req *dto.DriverLoginRequestDto, ip, userAgent string) (*dto.DriverLoginResponseDto, error) {
	device, err := u.ensureDriverDevice(ctx, &req.Device)
	if err != nil {
		return nil, err
	}

	type txResult struct {
		account       *model.Account
		driverProfile *appdrivermodel.DriverProfile
		refreshToken  string
	}

	res, err := database.WithTransaction(
		u.txManager,
		ctx,
		func(txCtx context.Context) (txResult, error) {
			account, driverProfile, err := u.getAccountAndDriverForLogin(txCtx, req.Phone, req.Password)
			if err != nil {
				if u.loginHistoryRepo != nil {
					_ = u.logDriverLoginHistory(txCtx, uuid.Nil, device.ID, ip, userAgent, err, req.Phone)
				}
				return txResult{}, err
			}

			accountAppDevice, err := u.ensureDriverAccountAppDevice(txCtx, account, device, &req.Device)
			if err != nil {
				return txResult{}, err
			}

			refreshToken, err := u.createDriverLoginSession(txCtx, accountAppDevice.ID, ip, userAgent)
			if err != nil {
				return txResult{}, err
			}

			if u.loginHistoryRepo != nil {
				_ = u.logDriverLoginHistory(txCtx, account.ID, device.ID, ip, userAgent, nil, "")
			}

			return txResult{account: account, driverProfile: driverProfile, refreshToken: refreshToken}, nil
		},
	)
	if err != nil {
		if errors.Is(err, ErrDriverRequireVerifyOtp) {
			return &dto.DriverLoginResponseDto{
				RequireVerifyOtp: true,
				Message:          ErrDriverRequireVerifyOtp.Error(),
			}, nil
		}
		return nil, err
	}

	accessToken, err := jwtutil.GenerateAccessToken(res.account.ID)
	if err != nil {
		return nil, err
	}

	profileDto := appdrivertransformer.ToDriverProfileItemDto(res.account, res.driverProfile)
	return &dto.DriverLoginResponseDto{
		RequireVerifyOtp: false,
		AccessToken:      accessToken,
		RefreshToken:     res.refreshToken,
		DriverProfile:    &profileDto,
	}, nil
}

func (u *driverProfileUsecase) getAccountAndDriverForLogin(ctx context.Context, phone, password string) (*model.Account, *appdrivermodel.DriverProfile, error) {
	account, err := u.accountRepo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, nil, err
	}
	if account == nil {
		return nil, nil, ErrDriverNotFound
	}
	if !validator.CheckPassword(password, account.PasswordHash) {
		return nil, nil, ErrDriverInvalidPassword
	}

	driverProfile, err := u.driverProfileRepo.GetByAccountID(ctx, account.ID)
	if err != nil {
		return nil, nil, err
	}
	if driverProfile == nil {
		return nil, nil, ErrDriverNotFound
	}
	if driverProfile.GlobalStatus == string(pgdb.DriverProfileStatusPENDINGVERIFICATION) {
		return nil, nil, ErrDriverRequireVerifyOtp
	}

	return account, driverProfile, nil
}

func (u *driverProfileUsecase) ensureDriverDevice(ctx context.Context, d *dto.DriverDeviceDto) (*model.Device, error) {
	device, err := u.deviceRepo.GetDeviceByUID(ctx, d.DeviceUID)
	if err != nil {
		return nil, err
	}
	if device != nil {
		return device, nil
	}
	return u.deviceRepo.CreateDevice(ctx, &model.Device{
		DeviceUID:  d.DeviceUID,
		Platform:   d.Platform,
		DeviceName: d.DeviceName,
		OsVersion:  d.OsVersion,
		AppVersion: d.AppVersion,
	})
}

func (u *driverProfileUsecase) ensureDriverAccountAppDevice(ctx context.Context, account *model.Account, device *model.Device, d *dto.DriverDeviceDto) (*model.AccountAppDevice, error) {
	aad, err := u.accountAppDeviceRepo.GetByAccountDeviceAndAppType(ctx, account.ID, device.ID, appTypeDriver)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	if aad == nil {
		return u.accountAppDeviceRepo.CreateAccountAppDevice(ctx, &model.AccountAppDevice{
			AccountID:  account.ID,
			DeviceID:   device.ID,
			AppType:    appTypeDriver,
			FcmToken:   d.FCMToken,
			IsActive:   true,
			LastUsedAt: &now,
		})
	}
	aad.IsActive = true
	aad.LastUsedAt = &now
	if d.FCMToken != "" {
		aad.FcmToken = d.FCMToken
	}
	return u.accountAppDeviceRepo.UpdateAccountAppDevice(ctx, aad)
}

func (u *driverProfileUsecase) createDriverLoginSession(ctx context.Context, accountAppDeviceID uuid.UUID, ip, userAgent string) (string, error) {
	refreshToken := generate.GenerateRandomString(64)
	expiresAt := time.Now().Add(refreshTokenTTL)
	_, err := u.sessionRepo.CreateSession(ctx, &model.Session{
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

func (u *driverProfileUsecase) GoOnline(ctx context.Context, accountID uuid.UUID, req *dto.DriverLocationStatusRequestDto) error {
	driverProfile, err := u.driverProfileRepo.GetByAccountID(ctx, accountID)
	if err != nil {
		return err
	}
	if driverProfile == nil {
		return ErrDriverNotFound
	}
	if driverProfile.GlobalStatus != string(pgdb.DriverProfileStatusACTIVE) {
		return ErrDriverNotActive
	}

	return u.setDriverOnline(ctx, driverProfile.ID, req)
}

func (u *driverProfileUsecase) GoOffline(ctx context.Context, accountID uuid.UUID) error {
	driverProfile, err := u.driverProfileRepo.GetByAccountID(ctx, accountID)
	if err != nil {
		return err
	}
	if driverProfile == nil {
		return ErrDriverNotFound
	}

	key := fmt.Sprintf("driver:%s:online", driverProfile.ID.String())
	return u.redisClient.Del(ctx, key).Err()
}

func (u *driverProfileUsecase) PingOnline(ctx context.Context, accountID uuid.UUID, req *dto.DriverLocationStatusRequestDto) error {
	driverProfile, err := u.driverProfileRepo.GetByAccountID(ctx, accountID)
	if err != nil {
		return err
	}
	if driverProfile == nil {
		return ErrDriverNotFound
	}
	if driverProfile.GlobalStatus != string(pgdb.DriverProfileStatusACTIVE) {
		return ErrDriverNotActive
	}

	if u.redisClient == nil {
		return errors.New("redis client is not configured")
	}

	pingKey := fmt.Sprintf("driver:%s:ping_count", driverProfile.ID.String())
	count, err := u.redisClient.Incr(ctx, pingKey).Result()
	if err != nil {
		return err
	}
	if count == 1 {
		_ = u.redisClient.Expire(ctx, pingKey, driverPingLimitWindow).Err()
	}
	if count > driverPingLimitCount {
		return ErrDriverPingTooFrequent
	}

	return u.setDriverOnline(ctx, driverProfile.ID, req)
}

func (u *driverProfileUsecase) setDriverOnline(ctx context.Context, driverID uuid.UUID, req *dto.DriverLocationStatusRequestDto) error {
	if u.redisClient == nil {
		return errors.New("redis client is not configured")
	}

	val := driverOnlineValue{
		Lat:    req.Lat,
		Lng:    req.Lng,
		Status: req.Status,
	}
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("driver:%s:online", driverID.String())
	return u.redisClient.Set(ctx, key, data, driverOnlineTTL).Err()
}

func (u *driverProfileUsecase) logDriverLoginHistory(ctx context.Context, accountID uuid.UUID, deviceID uuid.UUID, ip, userAgent string, loginErr error, phone string) error {
	var result, reason string
	if loginErr == nil {
		result = "success"
	} else {
		if accountID == uuid.Nil && phone != "" {
			acc, _ := u.accountRepo.GetByPhone(ctx, phone)
			if acc != nil {
				accountID = acc.ID
			}
		}
		switch {
		case errors.Is(loginErr, ErrDriverNotFound):
			result = "failed_not_found"
			reason = "Tài khoản không tồn tại"
		case errors.Is(loginErr, ErrDriverInvalidPassword):
			result = "failed_password"
			reason = "Mật khẩu không chính xác"
		case errors.Is(loginErr, ErrDriverRequireVerifyOtp):
			result = "failed_require_verify_otp"
			reason = ErrDriverRequireVerifyOtp.Error()
		default:
			result = "failed_unknown"
			reason = loginErr.Error()
		}
	}
	_, err := u.loginHistoryRepo.CreateLoginHistory(ctx, &model.AppLoginHistory{
		AccountID:     accountID,
		DeviceID:      deviceID,
		AppType:       appTypeDriver,
		Result:        result,
		FailureReason: reason,
		IpAddress:     ip,
		UserAgent:     userAgent,
	})
	return err
}
