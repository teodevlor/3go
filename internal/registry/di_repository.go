package registry

import (
	account_repo "go-structure/internal/repository"
	user_profile_repo "go-structure/internal/repository/app_user"
	setting_repo "go-structure/internal/repository/web_system"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sarulabs/di"
)

const (
	UserProfileRepoDIName         = "user_profile_repo_di"
	AccountRepoDIName             = "account_repo_di"
	OTPRepoDIName                 = "otp_repo_di"
	OTPAuditRepoDIName            = "otp_audit_repo_di"
	SettingRepoDIName             = "setting_repo_di"
	DeviceRepoDIName              = "device_repo_di"
	AccountAppDeviceRepoDIName    = "account_app_device_repo_di"
	SessionRepoDIName             = "session_repo_di"
	LoginHistoryRepoDIName        = "login_history_repo_di"
	ZoneRepoDIName                = "zone_repo_di"
	SidebarRepoDIName             = "sidebar_repo_di"
	ServiceRepoDIName             = "service_repo_di"
	ServiceZoneRepoDIName         = "service_zone_repo_di"
	DistancePricingRuleRepoDIName = "distance_pricing_rule_repo_di"
	SurchargeRuleRepoDIName       = "surcharge_rule_repo_di"
	PackageSizePricingRepoDIName  = "package_size_pricing_repo_di"
	SystemAdminRepoDIName                = "system_admin_repo_di"
	SystemLoginHistoryRepoDIName         = "system_login_history_repo_di"
	SystemAdminRefreshTokenRepoDIName    = "system_admin_refresh_token_repo_di"
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
			return setting_repo.NewSettingRepository(pool), nil
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

	zoneDef := di.Def{
		Name:  ZoneRepoDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			pool := ctn.Get(PostgresPoolDIName).(*pgxpool.Pool)
			return account_repo.NewZoneRepository(pool), nil
		},
	}

	sidebarDef := di.Def{
		Name:  SidebarRepoDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			pool := ctn.Get(PostgresPoolDIName).(*pgxpool.Pool)
			return setting_repo.NewSidebarRepository(pool), nil
		},
	}

	serviceDef := di.Def{
		Name:  ServiceRepoDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			pool := ctn.Get(PostgresPoolDIName).(*pgxpool.Pool)
			return setting_repo.NewServiceRepository(pool), nil
		},
	}

	serviceZoneDef := di.Def{
		Name:  ServiceZoneRepoDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			pool := ctn.Get(PostgresPoolDIName).(*pgxpool.Pool)
			return setting_repo.NewServiceZoneRepository(pool), nil
		},
	}

	distancePricingRuleDef := di.Def{
		Name:  DistancePricingRuleRepoDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			pool := ctn.Get(PostgresPoolDIName).(*pgxpool.Pool)
			return setting_repo.NewDistancePricingRuleRepository(pool), nil
		},
	}

	surchargeRuleDef := di.Def{
		Name:  SurchargeRuleRepoDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			pool := ctn.Get(PostgresPoolDIName).(*pgxpool.Pool)
			return setting_repo.NewSurchargeRuleRepository(pool), nil
		},
	}

	packageSizePricingDef := di.Def{
		Name:  PackageSizePricingRepoDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			pool := ctn.Get(PostgresPoolDIName).(*pgxpool.Pool)
			return setting_repo.NewPackageSizePricingRepository(pool), nil
		},
	}

	systemAdminDef := di.Def{
		Name:  SystemAdminRepoDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			pool := ctn.Get(PostgresPoolDIName).(*pgxpool.Pool)
			return setting_repo.NewSystemAdminRepository(pool), nil
		},
	}

	systemLoginHistoryDef := di.Def{
		Name:  SystemLoginHistoryRepoDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			pool := ctn.Get(PostgresPoolDIName).(*pgxpool.Pool)
			return setting_repo.NewSystemLoginHistoryRepository(pool), nil
		},
	}

	systemAdminRefreshTokenDef := di.Def{
		Name:  SystemAdminRefreshTokenRepoDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			pool := ctn.Get(PostgresPoolDIName).(*pgxpool.Pool)
			return setting_repo.NewSystemAdminRefreshTokenRepository(pool), nil
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
		zoneDef,
		sidebarDef,
		serviceDef,
		serviceZoneDef,
		distancePricingRuleDef,
		surchargeRuleDef,
		packageSizePricingDef,
		systemAdminDef,
		systemLoginHistoryDef,
		systemAdminRefreshTokenDef,
	)
}
