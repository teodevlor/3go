package registry

import (
	account_repo "go-structure/internal/repository"
	user_profile_repo "go-structure/internal/repository/app_user"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sarulabs/di"
)

const (
	UserProfileRepoDIName = "user_profile_repo_di"
	AccountRepoDIName     = "account_repo_di"
	OTPRepoDIName         = "otp_repo_di"
	OTPAuditRepoDIName    = "otp_audit_repo_di"
	SettingRepoDIName     = "setting_repo_di"
	DeviceRepoDIName           = "device_repo_di"
	AccountAppDeviceRepoDIName = "account_app_device_repo_di"
	SessionRepoDIName          = "session_repo_di"
	LoginHistoryRepoDIName     = "login_history_repo_di"
)

func buildRepositories() error {
	userProfileDef := di.Def{
		Name:  UserProfileRepoDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			pool := ctn.Get(PostgresPoolDIName).(*pgxpool.Pool)
			return user_profile_repo.NewUserProfileRepository(pool), nil
		},
	}

	accountDef := di.Def{
		Name:  AccountRepoDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			pool := ctn.Get(PostgresPoolDIName).(*pgxpool.Pool)
			return account_repo.NewAccountRepository(pool), nil
		},
	}

	otpDef := di.Def{
		Name:  OTPRepoDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			pool := ctn.Get(PostgresPoolDIName).(*pgxpool.Pool)
			return account_repo.NewOTPRepository(pool), nil
		},
	}

	otpAuditDef := di.Def{
		Name:  OTPAuditRepoDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			pool := ctn.Get(PostgresPoolDIName).(*pgxpool.Pool)
			return account_repo.NewOTPAuditRepository(pool), nil
		},
	}

	settingDef := di.Def{
		Name:  SettingRepoDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			pool := ctn.Get(PostgresPoolDIName).(*pgxpool.Pool)
			return account_repo.NewSettingRepository(pool), nil
		},
	}

	deviceDef := di.Def{
		Name:  DeviceRepoDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			pool := ctn.Get(PostgresPoolDIName).(*pgxpool.Pool)
			return account_repo.NewDeviceRepository(pool), nil
		},
	}

	accountAppDeviceDef := di.Def{
		Name:  AccountAppDeviceRepoDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			pool := ctn.Get(PostgresPoolDIName).(*pgxpool.Pool)
			return account_repo.NewAccountAppDeviceRepository(pool), nil
		},
	}

	sessionDef := di.Def{
		Name:  SessionRepoDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			pool := ctn.Get(PostgresPoolDIName).(*pgxpool.Pool)
			return account_repo.NewSessionRepository(pool), nil
		},
	}

	loginHistoryDef := di.Def{
		Name:  LoginHistoryRepoDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			pool := ctn.Get(PostgresPoolDIName).(*pgxpool.Pool)
			return account_repo.NewLoginHistoryRepository(pool), nil
		},
	}

	return builder.Add(
		userProfileDef,
		accountDef,
		otpDef,
		otpAuditDef,
		settingDef,
		deviceDef,
		accountAppDeviceDef,
		sessionDef,
		loginHistoryDef,
	)
}
