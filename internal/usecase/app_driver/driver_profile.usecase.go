package app_driver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go-structure/global"
	common "go-structure/internal/common"
	"go-structure/internal/constants"
	dto "go-structure/internal/dto/app_driver"
	dto_common "go-structure/internal/dto/common"
	"go-structure/internal/helper/database"
	"go-structure/internal/middleware"
	"go-structure/internal/repository"
	appdriverrepo "go-structure/internal/repository/app_driver"
	"go-structure/internal/repository/model"
	appdrivermodel "go-structure/internal/repository/model/app_driver"
	appdrivertransformer "go-structure/internal/transformer/app_driver"
	"go-structure/internal/usecase"
	"go-structure/internal/utils/generate"
	jwtutil "go-structure/internal/utils/jwt"
	pgdb "go-structure/orm/db/postgres"
	"go-structure/pkg/validator"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	appTypeDriver   = common.AppTypeDriver
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
	DriverProfileListFilter struct {
		Page         int
		Limit        int
		Search       string
		GlobalStatus string
	}

	driverOnlineValue struct {
		Lat    float64 `json:"lat"`
		Lng    float64 `json:"lng"`
		Status string  `json:"status"`
	}

	IDriverProfileUsecase interface {
		RegisterDriver(ctx context.Context, req *dto.DriverRegisterRequestDto) (*dto.DriverRegisterResponseDto, error)
		VerifyDriverOtp(ctx context.Context, phone, code, ip, userAgent string) (*dto.DriverVerifyOtpResponseDto, error)
		LoginDriver(ctx context.Context, req *dto.DriverLoginRequestDto, ip, userAgent string) (*dto.DriverLoginResponseDto, error)

		AdminCreateDriverProfile(ctx context.Context, req *dto.AdminCreateDriverProfileRequestDto) (*dto.DriverProfileItemDto, error)

		GoOnline(ctx context.Context, accountID uuid.UUID, req *dto.DriverLocationStatusRequestDto) error
		GoOffline(ctx context.Context, accountID uuid.UUID) error
		PingOnline(ctx context.Context, accountID uuid.UUID, req *dto.DriverLocationStatusRequestDto) error

		GetByID(ctx context.Context, id uuid.UUID) (*dto.DriverProfileItemDto, error)
		List(ctx context.Context, filter DriverProfileListFilter) (*dto.ListDriverProfilesResponseDto, error)
		UpdateProfile(ctx context.Context, id uuid.UUID, req *dto.UpdateDriverProfileRequestDto) (*dto.DriverProfileItemDto, error)
		DeleteProfile(ctx context.Context, id uuid.UUID) error
	}

	driverProfileUsecase struct {
		driverProfileRepo    appdriverrepo.IDriverProfileRepository
		driverServiceRepo    appdriverrepo.IDriverServiceRepository
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
	driverServiceRepo appdriverrepo.IDriverServiceRepository,
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
		driverServiceRepo:    driverServiceRepo,
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

func (u *driverProfileUsecase) RegisterDriver(ctx context.Context, req *dto.DriverRegisterRequestDto) (*dto.DriverRegisterResponseDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("RegisterDriver: start", zap.String(global.KeyCorrelationID, cid), zap.String("phone", req.Phone))

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
			driverID, err := u.createDriverProfileRecord(txCtx, account.ID, req.FullName)
			if err != nil {
				return "", err
			}
			if u.driverServiceRepo != nil && len(req.ServiceIDs) > 0 {
				global.Logger.Info("RegisterDriver: assigning services", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", driverID.String()), zap.Int("service_count", len(req.ServiceIDs)))
				if err := u.driverServiceRepo.SetDriverServices(txCtx, driverID, req.ServiceIDs); err != nil {
					global.Logger.Error("RegisterDriver: failed to set driver services", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", driverID.String()), zap.Error(err))
					return "", err
				}
			}
			if u.otpUsecase != nil {
				code, err := u.otpUsecase.CreateOTP(txCtx, req.Phone, usecase.OTPPurposeDriverRegister)
				if err != nil {
					global.Logger.Error("RegisterDriver: failed to create OTP", zap.String(global.KeyCorrelationID, cid), zap.String("phone", req.Phone), zap.Error(err))
					return "", err
				}
				return code, nil
			}
			return "", nil
		},
	)
	if err != nil {
		global.Logger.Error("RegisterDriver: transaction failed", zap.String(global.KeyCorrelationID, cid), zap.String("phone", req.Phone), zap.Error(err))
		return nil, err
	}

	if u.notifyUsecase != nil && otpCode != "" {
		global.Logger.Info("RegisterDriver: sending OTP notification", zap.String(global.KeyCorrelationID, cid), zap.String("phone", req.Phone))
		msg := fmt.Sprintf("Đăng ký tài khoản tài xế thành công. Mã OTP xác thực: %s", otpCode)
		_ = u.notifyUsecase.SendOtp(ctx, msg)
	}

	global.Logger.Info("RegisterDriver: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("phone", req.Phone))
	return &dto.DriverRegisterResponseDto{
		UserMessage: "Đăng ký tài khoản tài xế thành công, vui lòng kiểm tra điện thoại để nhận mã OTP",
	}, nil
}

func (u *driverProfileUsecase) getOrCreateAccount(ctx context.Context, req *dto.DriverRegisterRequestDto) (*model.Account, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("getOrCreateAccount: looking up existing account", zap.String(global.KeyCorrelationID, cid), zap.String("phone", req.Phone))

	acc, err := u.accountRepo.GetByPhone(ctx, req.Phone)
	if err != nil {
		global.Logger.Error("getOrCreateAccount: failed to query account by phone", zap.String(global.KeyCorrelationID, cid), zap.String("phone", req.Phone), zap.Error(err))
		return nil, err
	}
	if acc != nil {
		global.Logger.Info("getOrCreateAccount: existing account found, reusing", zap.String(global.KeyCorrelationID, cid), zap.String("phone", req.Phone), zap.String("account_id", acc.ID.String()))
		return acc, nil
	}

	hashedPassword, err := validator.HashPassword(req.Password)
	if err != nil {
		global.Logger.Error("getOrCreateAccount: failed to hash password", zap.String(global.KeyCorrelationID, cid), zap.String("phone", req.Phone), zap.Error(err))
		return nil, err
	}
	account, err := u.accountRepo.CreateAccount(ctx, &model.Account{
		Phone:        req.Phone,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		global.Logger.Error("getOrCreateAccount: failed to create account", zap.String(global.KeyCorrelationID, cid), zap.String("phone", req.Phone), zap.Error(err))
		return nil, err
	}

	global.Logger.Info("getOrCreateAccount: new account created", zap.String(global.KeyCorrelationID, cid), zap.String("phone", req.Phone), zap.String("account_id", account.ID.String()))
	return account, nil
}

func (u *driverProfileUsecase) ensureDriverNotAlreadyRegistered(ctx context.Context, accountID uuid.UUID) error {
	cid := middleware.CorrelationIDFromContext(ctx)

	existing, err := u.driverProfileRepo.GetByAccountID(ctx, accountID)
	if err != nil {
		global.Logger.Error("ensureDriverNotAlreadyRegistered: failed to query driver profile", zap.String(global.KeyCorrelationID, cid), zap.String("account_id", accountID.String()), zap.Error(err))
		return err
	}
	if existing != nil {
		global.Logger.Error("ensureDriverNotAlreadyRegistered: driver already registered", zap.String(global.KeyCorrelationID, cid), zap.String("account_id", accountID.String()), zap.String("driver_id", existing.ID.String()))
		return ErrDriverAlreadyRegistered
	}
	return nil
}

func (u *driverProfileUsecase) createDriverProfileRecord(ctx context.Context, accountID uuid.UUID, fullName string) (uuid.UUID, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("createDriverProfileRecord: creating driver profile", zap.String(global.KeyCorrelationID, cid), zap.String("account_id", accountID.String()))

	profile, err := u.driverProfileRepo.Create(ctx, accountID, fullName)
	if err != nil {
		global.Logger.Error("createDriverProfileRecord: failed to create driver profile", zap.String(global.KeyCorrelationID, cid), zap.String("account_id", accountID.String()), zap.Error(err))
		return uuid.Nil, err
	}

	global.Logger.Info("createDriverProfileRecord: driver profile created", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", profile.ID.String()))
	return profile.ID, nil
}

func (u *driverProfileUsecase) VerifyDriverOtp(ctx context.Context, phone, code, ip, userAgent string) (*dto.DriverVerifyOtpResponseDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("VerifyDriverOtp: start", zap.String(global.KeyCorrelationID, cid), zap.String("phone", phone))

	account, err := u.accountRepo.GetByPhone(ctx, phone)
	if err != nil {
		global.Logger.Error("VerifyDriverOtp: failed to get account by phone", zap.String(global.KeyCorrelationID, cid), zap.String("phone", phone), zap.Error(err))
		return nil, err
	}
	if account == nil {
		global.Logger.Error("VerifyDriverOtp: account not found", zap.String(global.KeyCorrelationID, cid), zap.String("phone", phone))
		return nil, ErrDriverNotFound
	}

	driverProfile, err := u.driverProfileRepo.GetByAccountID(ctx, account.ID)
	if err != nil {
		global.Logger.Error("VerifyDriverOtp: failed to get driver profile", zap.String(global.KeyCorrelationID, cid), zap.String("account_id", account.ID.String()), zap.Error(err))
		return nil, err
	}
	if driverProfile == nil {
		global.Logger.Error("VerifyDriverOtp: driver profile not found", zap.String(global.KeyCorrelationID, cid), zap.String("account_id", account.ID.String()))
		return nil, ErrDriverNotFound
	}

	global.Logger.Info("VerifyDriverOtp: verifying OTP code", zap.String(global.KeyCorrelationID, cid), zap.String("phone", phone), zap.String("driver_id", driverProfile.ID.String()))
	verified, err := u.otpUsecase.VerifyOTP(ctx, phone, code, usecase.OTPPurposeDriverRegister, ip, userAgent)
	if err != nil {
		global.Logger.Error("VerifyDriverOtp: OTP verification error", zap.String(global.KeyCorrelationID, cid), zap.String("phone", phone), zap.Error(err))
		return nil, err
	}
	if !verified {
		global.Logger.Error("VerifyDriverOtp: invalid OTP code", zap.String(global.KeyCorrelationID, cid), zap.String("phone", phone))
		return nil, ErrDriverInvalidOTP
	}

	fromStatus := pgdb.NullDriverProfileStatus{
		DriverProfileStatus: pgdb.DriverProfileStatus(driverProfile.GlobalStatus),
		Valid:               true,
	}
	global.Logger.Info("VerifyDriverOtp: updating driver status to DOCUMENT_INCOMPLETE", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", driverProfile.ID.String()), zap.String("from_status", driverProfile.GlobalStatus))
	_, err = database.WithTransaction(
		u.txManager,
		ctx,
		func(txCtx context.Context) (struct{}, error) {
			if err := u.driverProfileRepo.CreateStatusHistory(txCtx, driverProfile.ID, fromStatus, pgdb.DriverProfileStatusDOCUMENTINCOMPLETE, nil, nil); err != nil {
				global.Logger.Error("VerifyDriverOtp: failed to create status history", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", driverProfile.ID.String()), zap.Error(err))
				return struct{}{}, err
			}
			_, err := u.driverProfileRepo.UpdateStatus(txCtx, driverProfile.ID, pgdb.DriverProfileStatusDOCUMENTINCOMPLETE)
			if err != nil {
				global.Logger.Error("VerifyDriverOtp: failed to update driver status", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", driverProfile.ID.String()), zap.Error(err))
			}
			return struct{}{}, err
		},
	)
	if err != nil {
		return nil, err
	}

	global.Logger.Info("VerifyDriverOtp: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("phone", phone), zap.String("driver_id", driverProfile.ID.String()))
	return &dto.DriverVerifyOtpResponseDto{
		UserMessage: "Xác thực OTP thành công. Bạn có thể đăng nhập và upload tài liệu.",
		DriverID:    driverProfile.ID,
	}, nil
}

func (u *driverProfileUsecase) LoginDriver(ctx context.Context, req *dto.DriverLoginRequestDto, ip, userAgent string) (*dto.DriverLoginResponseDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("LoginDriver: start", zap.String(global.KeyCorrelationID, cid), zap.String("phone", req.Phone), zap.String("ip", ip))

	device, err := u.ensureDriverDevice(ctx, &req.Device)
	if err != nil {
		global.Logger.Error("LoginDriver: failed to ensure device", zap.String(global.KeyCorrelationID, cid), zap.String("phone", req.Phone), zap.String("device_uid", req.Device.DeviceUID), zap.Error(err))
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
				global.Logger.Error("LoginDriver: authentication failed", zap.String(global.KeyCorrelationID, cid), zap.String("phone", req.Phone), zap.Error(err))
				if u.loginHistoryRepo != nil {
					_ = u.logDriverLoginHistory(txCtx, uuid.Nil, device.ID, ip, userAgent, err, req.Phone)
				}
				return txResult{}, err
			}

			global.Logger.Info("LoginDriver: auth passed, ensuring account-app-device link", zap.String(global.KeyCorrelationID, cid), zap.String("account_id", account.ID.String()))
			accountAppDevice, err := u.ensureDriverAccountAppDevice(txCtx, account, device, &req.Device)
			if err != nil {
				global.Logger.Error("LoginDriver: failed to ensure account app device", zap.String(global.KeyCorrelationID, cid), zap.String("account_id", account.ID.String()), zap.Error(err))
				return txResult{}, err
			}

			refreshToken, err := u.createDriverLoginSession(txCtx, accountAppDevice.ID, ip, userAgent)
			if err != nil {
				global.Logger.Error("LoginDriver: failed to create login session", zap.String(global.KeyCorrelationID, cid), zap.String("account_id", account.ID.String()), zap.Error(err))
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
			global.Logger.Info("LoginDriver: OTP verification required", zap.String(global.KeyCorrelationID, cid), zap.String("phone", req.Phone))
			return &dto.DriverLoginResponseDto{
				RequireVerifyOtp: true,
				Message:          ErrDriverRequireVerifyOtp.Error(),
			}, nil
		}
		global.Logger.Error("LoginDriver: transaction failed", zap.String(global.KeyCorrelationID, cid), zap.String("phone", req.Phone), zap.Error(err))
		return nil, err
	}

	accessToken, err := jwtutil.GenerateAccessToken(res.account.ID)
	if err != nil {
		global.Logger.Error("LoginDriver: failed to generate access token", zap.String(global.KeyCorrelationID, cid), zap.String("account_id", res.account.ID.String()), zap.Error(err))
		return nil, err
	}

	var serviceIDs []uuid.UUID
	if u.driverServiceRepo != nil {
		if ids, err := u.driverServiceRepo.GetServiceIDsByDriverID(ctx, res.driverProfile.ID); err != nil {
			global.Logger.Error("LoginDriver: failed to load driver services", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", res.driverProfile.ID.String()), zap.Error(err))
		} else {
			serviceIDs = ids
		}
	}

	global.Logger.Info("LoginDriver: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("phone", req.Phone), zap.String("account_id", res.account.ID.String()))
	profileDto := appdrivertransformer.ToDriverProfileItemDto(res.account, res.driverProfile)
	profileDto.ServiceIDs = serviceIDs
	return &dto.DriverLoginResponseDto{
		RequireVerifyOtp: false,
		AccessToken:      accessToken,
		RefreshToken:     res.refreshToken,
		DriverProfile:    &profileDto,
	}, nil
}

func (u *driverProfileUsecase) getAccountAndDriverForLogin(ctx context.Context, phone, password string) (*model.Account, *appdrivermodel.DriverProfile, error) {
	cid := middleware.CorrelationIDFromContext(ctx)

	account, err := u.accountRepo.GetByPhone(ctx, phone)
	if err != nil {
		global.Logger.Error("getAccountAndDriverForLogin: failed to query account by phone", zap.String(global.KeyCorrelationID, cid), zap.String("phone", phone), zap.Error(err))
		return nil, nil, err
	}
	if account == nil {
		global.Logger.Error("getAccountAndDriverForLogin: account not found", zap.String(global.KeyCorrelationID, cid), zap.String("phone", phone))
		return nil, nil, ErrDriverNotFound
	}
	if !validator.CheckPassword(password, account.PasswordHash) {
		global.Logger.Error("getAccountAndDriverForLogin: invalid password", zap.String(global.KeyCorrelationID, cid), zap.String("phone", phone))
		return nil, nil, ErrDriverInvalidPassword
	}

	driverProfile, err := u.driverProfileRepo.GetByAccountID(ctx, account.ID)
	if err != nil {
		global.Logger.Error("getAccountAndDriverForLogin: failed to query driver profile", zap.String(global.KeyCorrelationID, cid), zap.String("account_id", account.ID.String()), zap.Error(err))
		return nil, nil, err
	}
	if driverProfile == nil {
		global.Logger.Error("getAccountAndDriverForLogin: driver profile not found", zap.String(global.KeyCorrelationID, cid), zap.String("account_id", account.ID.String()))
		return nil, nil, ErrDriverNotFound
	}
	if driverProfile.GlobalStatus == string(pgdb.DriverProfileStatusPENDINGVERIFICATION) {
		global.Logger.Info("getAccountAndDriverForLogin: driver requires OTP verification", zap.String(global.KeyCorrelationID, cid), zap.String("account_id", account.ID.String()))
		return nil, nil, ErrDriverRequireVerifyOtp
	}

	return account, driverProfile, nil
}

func (u *driverProfileUsecase) ensureDriverDevice(ctx context.Context, d *dto.DriverDeviceDto) (*model.Device, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	device, err := u.deviceRepo.GetDeviceByUID(ctx, d.DeviceUID)
	if err != nil {
		global.Logger.Error("ensureDriverDevice: failed to get device by UID", zap.String(global.KeyCorrelationID, cid), zap.String("device_uid", d.DeviceUID), zap.Error(err))
		return nil, err
	}
	if device != nil {
		return device, nil
	}
	device, err = u.deviceRepo.CreateDevice(ctx, &model.Device{
		DeviceUID:  d.DeviceUID,
		Platform:   d.Platform,
		DeviceName: d.DeviceName,
		OsVersion:  d.OsVersion,
		AppVersion: d.AppVersion,
	})
	if err != nil {
		global.Logger.Error("ensureDriverDevice: failed to create device", zap.String(global.KeyCorrelationID, cid), zap.String("device_uid", d.DeviceUID), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("ensureDriverDevice: new device created", zap.String(global.KeyCorrelationID, cid), zap.String("device_uid", d.DeviceUID), zap.String("device_id", device.ID.String()))
	return device, nil
}

func (u *driverProfileUsecase) ensureDriverAccountAppDevice(ctx context.Context, account *model.Account, device *model.Device, d *dto.DriverDeviceDto) (*model.AccountAppDevice, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	aad, err := u.accountAppDeviceRepo.GetByAccountDeviceAndAppType(ctx, account.ID, device.ID, appTypeDriver)
	if err != nil {
		global.Logger.Error("ensureDriverAccountAppDevice: failed to get account app device", zap.String(global.KeyCorrelationID, cid), zap.String("account_id", account.ID.String()), zap.String("device_id", device.ID.String()), zap.Error(err))
		return nil, err
	}
	now := time.Now()
	if aad == nil {
		aad, err = u.accountAppDeviceRepo.CreateAccountAppDevice(ctx, &model.AccountAppDevice{
			AccountID:  account.ID,
			DeviceID:   device.ID,
			AppType:    appTypeDriver,
			FcmToken:   d.FCMToken,
			IsActive:   true,
			LastUsedAt: &now,
		})
		if err != nil {
			global.Logger.Error("ensureDriverAccountAppDevice: failed to create account app device", zap.String(global.KeyCorrelationID, cid), zap.String("account_id", account.ID.String()), zap.Error(err))
			return nil, err
		}
		global.Logger.Info("ensureDriverAccountAppDevice: created account app device link", zap.String(global.KeyCorrelationID, cid), zap.String("account_id", account.ID.String()), zap.String("device_id", device.ID.String()))
		return aad, nil
	}
	aad.IsActive = true
	aad.LastUsedAt = &now
	if d.FCMToken != "" {
		aad.FcmToken = d.FCMToken
	}
	aad, err = u.accountAppDeviceRepo.UpdateAccountAppDevice(ctx, aad)
	if err != nil {
		global.Logger.Error("ensureDriverAccountAppDevice: failed to update account app device", zap.String(global.KeyCorrelationID, cid), zap.String("account_id", account.ID.String()), zap.Error(err))
		return nil, err
	}
	return aad, nil
}

func (u *driverProfileUsecase) createDriverLoginSession(ctx context.Context, accountAppDeviceID uuid.UUID, ip, userAgent string) (string, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
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
		global.Logger.Error("createDriverLoginSession: failed to create session", zap.String(global.KeyCorrelationID, cid), zap.String("account_app_device_id", accountAppDeviceID.String()), zap.Error(err))
		return "", err
	}
	return refreshToken, nil
}

func (u *driverProfileUsecase) GoOnline(ctx context.Context, accountID uuid.UUID, req *dto.DriverLocationStatusRequestDto) error {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("GoOnline: start", zap.String(global.KeyCorrelationID, cid), zap.String("account_id", accountID.String()))
	driverProfile, err := u.driverProfileRepo.GetByAccountID(ctx, accountID)
	if err != nil {
		global.Logger.Error("GoOnline: failed to get driver profile", zap.String(global.KeyCorrelationID, cid), zap.String("account_id", accountID.String()), zap.Error(err))
		return err
	}
	if driverProfile == nil {
		global.Logger.Error("GoOnline: driver not found", zap.String(global.KeyCorrelationID, cid), zap.String("account_id", accountID.String()))
		return ErrDriverNotFound
	}
	if driverProfile.GlobalStatus != string(pgdb.DriverProfileStatusACTIVE) {
		global.Logger.Error("GoOnline: driver not active", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", driverProfile.ID.String()), zap.String("status", driverProfile.GlobalStatus))
		return ErrDriverNotActive
	}
	if err := u.setDriverOnline(ctx, driverProfile.ID, req); err != nil {
		global.Logger.Error("GoOnline: failed to set driver online", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", driverProfile.ID.String()), zap.Error(err))
		return err
	}
	global.Logger.Info("GoOnline: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", driverProfile.ID.String()))
	return nil
}

func (u *driverProfileUsecase) GoOffline(ctx context.Context, accountID uuid.UUID) error {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("GoOffline: start", zap.String(global.KeyCorrelationID, cid), zap.String("account_id", accountID.String()))
	driverProfile, err := u.driverProfileRepo.GetByAccountID(ctx, accountID)
	if err != nil {
		global.Logger.Error("GoOffline: failed to get driver profile", zap.String(global.KeyCorrelationID, cid), zap.String("account_id", accountID.String()), zap.Error(err))
		return err
	}
	if driverProfile == nil {
		global.Logger.Error("GoOffline: driver not found", zap.String(global.KeyCorrelationID, cid), zap.String("account_id", accountID.String()))
		return ErrDriverNotFound
	}
	key := fmt.Sprintf("driver:%s:online", driverProfile.ID.String())
	if err := u.redisClient.Del(ctx, key).Err(); err != nil {
		global.Logger.Error("GoOffline: failed to remove driver from online", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", driverProfile.ID.String()), zap.Error(err))
		return err
	}
	global.Logger.Info("GoOffline: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", driverProfile.ID.String()))
	return nil
}

func (u *driverProfileUsecase) PingOnline(ctx context.Context, accountID uuid.UUID, req *dto.DriverLocationStatusRequestDto) error {
	cid := middleware.CorrelationIDFromContext(ctx)
	driverProfile, err := u.driverProfileRepo.GetByAccountID(ctx, accountID)
	if err != nil {
		global.Logger.Error("PingOnline: failed to get driver profile", zap.String(global.KeyCorrelationID, cid), zap.String("account_id", accountID.String()), zap.Error(err))
		return err
	}
	if driverProfile == nil {
		global.Logger.Error("PingOnline: driver not found", zap.String(global.KeyCorrelationID, cid), zap.String("account_id", accountID.String()))
		return ErrDriverNotFound
	}
	if driverProfile.GlobalStatus != string(pgdb.DriverProfileStatusACTIVE) {
		global.Logger.Error("PingOnline: driver not active", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", driverProfile.ID.String()), zap.String("status", driverProfile.GlobalStatus))
		return ErrDriverNotActive
	}
	if u.redisClient == nil {
		global.Logger.Error("PingOnline: redis client not configured", zap.String(global.KeyCorrelationID, cid))
		return errors.New("redis client is not configured")
	}
	pingKey := fmt.Sprintf("driver:%s:ping_count", driverProfile.ID.String())
	count, err := u.redisClient.Incr(ctx, pingKey).Result()
	if err != nil {
		global.Logger.Error("PingOnline: failed to increment ping count", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", driverProfile.ID.String()), zap.Error(err))
		return err
	}
	if count == 1 {
		_ = u.redisClient.Expire(ctx, pingKey, driverPingLimitWindow).Err()
	}
	if count > driverPingLimitCount {
		global.Logger.Error("PingOnline: ping too frequent", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", driverProfile.ID.String()), zap.Int64("count", count))
		return ErrDriverPingTooFrequent
	}
	if err := u.setDriverOnline(ctx, driverProfile.ID, req); err != nil {
		global.Logger.Error("PingOnline: failed to set driver online", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", driverProfile.ID.String()), zap.Error(err))
		return err
	}
	return nil
}

func (u *driverProfileUsecase) setDriverOnline(ctx context.Context, driverID uuid.UUID, req *dto.DriverLocationStatusRequestDto) error {
	cid := middleware.CorrelationIDFromContext(ctx)
	if u.redisClient == nil {
		global.Logger.Error("setDriverOnline: redis client not configured", zap.String(global.KeyCorrelationID, cid))
		return errors.New("redis client is not configured")
	}
	val := driverOnlineValue{
		Lat:    req.Lat,
		Lng:    req.Lng,
		Status: req.Status,
	}
	data, err := json.Marshal(val)
	if err != nil {
		global.Logger.Error("setDriverOnline: failed to marshal location data", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", driverID.String()), zap.Error(err))
		return err
	}
	key := fmt.Sprintf("driver:%s:online", driverID.String())
	if err := u.redisClient.Set(ctx, key, data, driverOnlineTTL).Err(); err != nil {
		global.Logger.Error("setDriverOnline: failed to set redis key", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", driverID.String()), zap.Error(err))
		return err
	}
	return nil
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

func (u *driverProfileUsecase) GetByID(ctx context.Context, id uuid.UUID) (*dto.DriverProfileItemDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("GetByID: start", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", id.String()))
	profile, err := u.driverProfileRepo.GetByID(ctx, id)
	if err != nil {
		global.Logger.Error("GetByID: failed to get driver profile", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", id.String()), zap.Error(err))
		return nil, err
	}
	if profile == nil {
		global.Logger.Error("GetByID: driver not found", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", id.String()))
		return nil, ErrDriverNotFound
	}
	account, err := u.accountRepo.GetById(ctx, profile.AccountID.String())
	if err != nil {
		global.Logger.Error("GetByID: failed to get account", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", id.String()), zap.String("account_id", profile.AccountID.String()), zap.Error(err))
		return nil, err
	}
	item := appdrivertransformer.ToDriverProfileItemDto(account, profile)
	global.Logger.Info("GetByID: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", id.String()))
	return &item, nil
}

func (u *driverProfileUsecase) List(ctx context.Context, filter DriverProfileListFilter) (*dto.ListDriverProfilesResponseDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("List: start", zap.String(global.KeyCorrelationID, cid), zap.Int("page", filter.Page), zap.Int("limit", filter.Limit), zap.String("search", filter.Search), zap.String("global_status", filter.GlobalStatus))
	page := filter.Page
	limit := filter.Limit
	if page < 1 {
		page = constants.DefaultPage
	}
	if limit < 1 || limit > constants.MaxLimit {
		limit = constants.DefaultLimit
	}
	offset := int32((page - 1) * limit)
	limit32 := int32(limit)
	total, err := u.driverProfileRepo.Count(ctx, filter.Search, filter.GlobalStatus)
	if err != nil {
		global.Logger.Error("List: failed to count driver profiles", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	profiles, err := u.driverProfileRepo.List(ctx, filter.Search, filter.GlobalStatus, limit32, offset)
	if err != nil {
		global.Logger.Error("List: failed to list driver profiles", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	items := make([]dto.DriverProfileItemDto, 0, len(profiles))
	for _, p := range profiles {
		account, err := u.accountRepo.GetById(ctx, p.AccountID.String())
		if err != nil {
			global.Logger.Error("List: failed to get account for driver", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", p.ID.String()), zap.Error(err))
			return nil, err
		}
		item := appdrivertransformer.ToDriverProfileItemDto(account, p)
		items = append(items, item)
	}
	global.Logger.Info("List: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.Int64("total", total), zap.Int("returned", len(items)))
	return &dto.ListDriverProfilesResponseDto{
		Items: items,
		Pagination: dto_common.PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil
}

func parseDriverProfileStatus(s string) pgdb.NullDriverProfileStatus {
	if s == "" {
		return pgdb.NullDriverProfileStatus{}
	}

	status := pgdb.DriverProfileStatusPENDINGPROFILE
	switch s {
	case "PENDING_PROFILE":
		status = pgdb.DriverProfileStatusPENDINGPROFILE
	case "DOCUMENT_INCOMPLETE":
		status = pgdb.DriverProfileStatusDOCUMENTINCOMPLETE
	case "PENDING_VERIFICATION":
		status = pgdb.DriverProfileStatusPENDINGVERIFICATION
	case "ACTIVE":
		status = pgdb.DriverProfileStatusACTIVE
	case "SUSPENDED":
		status = pgdb.DriverProfileStatusSUSPENDED
	case "REJECTED":
		status = pgdb.DriverProfileStatusREJECTED
	}

	return pgdb.NullDriverProfileStatus{
		DriverProfileStatus: status,
		Valid:               true,
	}
}

func (u *driverProfileUsecase) AdminCreateDriverProfile(ctx context.Context, req *dto.AdminCreateDriverProfileRequestDto) (*dto.DriverProfileItemDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	if req == nil {
		global.Logger.Error("AdminCreateDriverProfile: request is required", zap.String(global.KeyCorrelationID, cid))
		return nil, errors.New("request is required")
	}
	global.Logger.Info("AdminCreateDriverProfile: start", zap.String(global.KeyCorrelationID, cid), zap.String("phone", req.Phone), zap.String("full_name", req.FullName))
	createdProfile, err := database.WithTransaction(
		u.txManager,
		ctx,
		func(txCtx context.Context) (*appdrivermodel.DriverProfile, error) {
			// Tạo hoặc lấy account theo phone (không đổi mật khẩu nếu account đã tồn tại)
			accReq := &dto.DriverRegisterRequestDto{
				Phone:    req.Phone,
				FullName: req.FullName,
				Password: req.Password,
			}
			account, err := u.getOrCreateAccount(txCtx, accReq)
			if err != nil {
				return nil, err
			}

			// Đảm bảo chưa có hồ sơ driver cho account này
			if err := u.ensureDriverNotAlreadyRegistered(txCtx, account.ID); err != nil {
				return nil, err
			}

			profile, err := u.driverProfileRepo.Create(txCtx, account.ID, req.FullName)
			if err != nil {
				return nil, err
			}

			updateArg := pgdb.UpdateDriverProfileParams{
				ID: profile.ID,
			}
			if req.FullName != "" {
				updateArg.FullName = pgtype.Text{String: req.FullName, Valid: true}
			}
			updateArg.DateOfBirth = common.ParseYYYYMMDDToPgDate(req.DateOfBirth)
			if req.Gender != "" {
				updateArg.Gender = pgtype.Text{String: req.Gender, Valid: true}
			}
			if req.Address != "" {
				updateArg.Address = pgtype.Text{String: req.Address, Valid: true}
			}

			updated, err := u.driverProfileRepo.Update(txCtx, updateArg)
			if err != nil {
				return nil, err
			}
			if updated == nil {
				return nil, ErrDriverNotFound
			}

			// Đăng ký các dịch vụ mà driver chọn
			if u.driverServiceRepo != nil && len(req.ServiceIDs) > 0 {
				if err := u.driverServiceRepo.SetDriverServices(txCtx, updated.ID, req.ServiceIDs); err != nil {
					return nil, err
				}
			}

			return updated, nil
		},
	)
	if err != nil {
		global.Logger.Error("AdminCreateDriverProfile: transaction failed", zap.String(global.KeyCorrelationID, cid), zap.String("phone", req.Phone), zap.Error(err))
		return nil, err
	}
	account, err := u.accountRepo.GetById(ctx, createdProfile.AccountID.String())
	if err != nil {
		global.Logger.Error("AdminCreateDriverProfile: failed to get account", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", createdProfile.ID.String()), zap.Error(err))
		return nil, err
	}
	item := appdrivertransformer.ToDriverProfileItemDto(account, createdProfile)
	global.Logger.Info("AdminCreateDriverProfile: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", createdProfile.ID.String()), zap.String("phone", req.Phone))
	return &item, nil
}

func (u *driverProfileUsecase) UpdateProfile(ctx context.Context, id uuid.UUID, req *dto.UpdateDriverProfileRequestDto) (*dto.DriverProfileItemDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("UpdateProfile: start", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", id.String()))
	existing, err := u.driverProfileRepo.GetByID(ctx, id)
	if err != nil {
		global.Logger.Error("UpdateProfile: failed to get driver profile", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", id.String()), zap.Error(err))
		return nil, err
	}
	if existing == nil {
		global.Logger.Error("UpdateProfile: driver not found", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", id.String()))
		return nil, ErrDriverNotFound
	}
	arg := pgdb.UpdateDriverProfileParams{
		ID: id,
	}

	if req.FullName != "" {
		arg.FullName = pgtype.Text{String: req.FullName, Valid: true}
	}
	arg.DateOfBirth = common.ParseYYYYMMDDToPgDate(req.DateOfBirth)
	if req.Gender != "" {
		arg.Gender = pgtype.Text{String: req.Gender, Valid: true}
	}
	if req.Address != "" {
		arg.Address = pgtype.Text{String: req.Address, Valid: true}
	}
	if req.GlobalStatus != "" {
		arg.GlobalStatus = parseDriverProfileStatus(req.GlobalStatus)
	}

	updated, err := database.WithTransaction(
		u.txManager,
		ctx,
		func(txCtx context.Context) (*appdrivermodel.DriverProfile, error) {
			return u.driverProfileRepo.Update(txCtx, arg)
		},
	)
	if err != nil {
		global.Logger.Error("UpdateProfile: transaction failed", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", id.String()), zap.Error(err))
		return nil, err
	}
	if updated == nil {
		global.Logger.Error("UpdateProfile: driver not found after update", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", id.String()))
		return nil, ErrDriverNotFound
	}
	account, err := u.accountRepo.GetById(ctx, updated.AccountID.String())
	if err != nil {
		global.Logger.Error("UpdateProfile: failed to get account", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", id.String()), zap.Error(err))
		return nil, err
	}
	item := appdrivertransformer.ToDriverProfileItemDto(account, updated)
	global.Logger.Info("UpdateProfile: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", id.String()))
	return &item, nil
}

func (u *driverProfileUsecase) DeleteProfile(ctx context.Context, id uuid.UUID) error {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("DeleteProfile: start", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", id.String()))
	existing, err := u.driverProfileRepo.GetByID(ctx, id)
	if err != nil {
		global.Logger.Error("DeleteProfile: failed to get driver profile", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", id.String()), zap.Error(err))
		return err
	}
	if existing == nil {
		global.Logger.Error("DeleteProfile: driver not found", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", id.String()))
		return ErrDriverNotFound
	}
	_, err = database.WithTransaction(
		u.txManager,
		ctx,
		func(txCtx context.Context) (struct{}, error) {
			return struct{}{}, u.driverProfileRepo.Delete(txCtx, id)
		},
	)
	if err != nil {
		global.Logger.Error("DeleteProfile: transaction failed", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", id.String()), zap.Error(err))
		return err
	}
	global.Logger.Info("DeleteProfile: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", id.String()))
	return nil
}
