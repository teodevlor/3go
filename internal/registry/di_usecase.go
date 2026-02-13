package registry

import (
	"go-structure/internal/adapter"
	"go-structure/internal/helper/database"
	account_repo "go-structure/internal/repository"
	user_profile_repo "go-structure/internal/repository/app_user"
	settingRepository "go-structure/internal/repository/web_system"
	usecase_pkg "go-structure/internal/usecase"
	user_profile_usecase "go-structure/internal/usecase/app_user"
	websystem_usecase "go-structure/internal/usecase/web_system"

	"github.com/sarulabs/di"
)

const (
	UserProfileUsecaseDIName         = "user_profile_usecase_di"
	NotifyUsecaseDIName              = "notify_usecase_di"
	SettingUsecaseDIName             = "setting_usecase_di"
	OTPUsecaseDIName                 = "otp_usecase_di"
	ZoneUsecaseDIName                = "zone_usecase_di"
	SidebarUsecaseDIName             = "sidebar_usecase_di"
	ServiceUsecaseDIName             = "service_usecase_di"
	ServiceZoneUsecaseDIName         = "service_zone_usecase_di"
	DistancePricingRuleUsecaseDIName = "distance_pricing_rule_usecase_di"
	SurchargeRuleUsecaseDIName       = "surcharge_rule_usecase_di"
	PackageSizePricingUsecaseDIName  = "package_size_pricing_usecase_di"
	AuthAdminUsecaseDIName           = "auth_admin_usecase_di"
)

func buildUsecases() error {
	userProfileDef := di.Def{
		Name:  UserProfileUsecaseDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			usrProfileRepo := ctn.Get(UserProfileRepoDIName).(user_profile_repo.IUserProfileRepository)
			accountRepo := ctn.Get(AccountRepoDIName).(account_repo.IAccountRepository)
			deviceRepo := ctn.Get(DeviceRepoDIName).(account_repo.IDeviceRepository)
			accountAppDeviceRepo := ctn.Get(AccountAppDeviceRepoDIName).(account_repo.IAccountAppDeviceRepository)
			sessionRepo := ctn.Get(SessionRepoDIName).(account_repo.ISessionRepository)
			loginHistoryRepo := ctn.Get(LoginHistoryRepoDIName).(account_repo.ILoginHistoryRepository)
			txManager := ctn.Get(TransactionManagerDIName).(database.TransactionManager)
			notifyUcAny := ctn.Get(NotifyUsecaseDIName)

			var notifyUc usecase_pkg.INotifyUsecase
			if notifyUcAny != nil {
				notifyUc = notifyUcAny.(usecase_pkg.INotifyUsecase)
			}
			otpUcAny := ctn.Get(OTPUsecaseDIName)
			var otpUc usecase_pkg.IOTPUsecase
			if otpUcAny != nil {
				otpUc = otpUcAny.(usecase_pkg.IOTPUsecase)
			}
			return user_profile_usecase.NewUserProfileUsecase(
				usrProfileRepo,
				accountRepo,
				deviceRepo,
				accountAppDeviceRepo,
				sessionRepo,
				loginHistoryRepo,
				notifyUc,
				otpUc,
				txManager,
			), nil
		},
	}

	notifyDef := di.Def{
		Name:  NotifyUsecaseDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			ad := ctn.Get(TelegramAdapterDIName)
			if ad == nil {
				return nil, nil
			}
			return usecase_pkg.NewNotifyUsecase(ad.(adapter.INotifyAdapter)), nil
		},
	}

	settingDef := di.Def{
		Name:  SettingUsecaseDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			settingRepo := ctn.Get(SettingRepoDIName).(settingRepository.ISettingRepository)
			return websystem_usecase.NewSettingUsecase(settingRepo), nil
		},
	}

	otpDef := di.Def{
		Name:  OTPUsecaseDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			otpRepo := ctn.Get(OTPRepoDIName).(account_repo.IOTPRepository)
			otpAuditRepo := ctn.Get(OTPAuditRepoDIName).(account_repo.IOTPAuditRepository)
			settingUc := ctn.Get(SettingUsecaseDIName).(usecase_pkg.ISettingUsecase)
			txManager := ctn.Get(TransactionManagerDIName).(database.TransactionManager)
			notifyUcAny := ctn.Get(NotifyUsecaseDIName)

			var notifyUc usecase_pkg.INotifyUsecase
			if notifyUcAny != nil {
				notifyUc = notifyUcAny.(usecase_pkg.INotifyUsecase)
			}

			return usecase_pkg.NewOTPUsecase(otpRepo, otpAuditRepo, settingUc, notifyUc, txManager), nil
		},
	}

	zoneDef := di.Def{
		Name:  ZoneUsecaseDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			zoneRepo := ctn.Get(ZoneRepoDIName).(account_repo.IZoneRepository)
			txManager := ctn.Get(TransactionManagerDIName).(database.TransactionManager)
			return websystem_usecase.NewZoneUsecase(zoneRepo, txManager), nil
		},
	}

	sidebarDef := di.Def{
		Name:  SidebarUsecaseDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			sidebarRepo := ctn.Get(SidebarRepoDIName).(settingRepository.ISidebarRepository)
			txManager := ctn.Get(TransactionManagerDIName).(database.TransactionManager)
			return websystem_usecase.NewSidebarUsecase(sidebarRepo, txManager), nil
		},
	}

	serviceDef := di.Def{
		Name:  ServiceUsecaseDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			serviceRepo := ctn.Get(ServiceRepoDIName).(settingRepository.IServiceRepository)
			serviceZoneUc := ctn.Get(ServiceZoneUsecaseDIName).(websystem_usecase.IServiceZoneUsecase)
			txManager := ctn.Get(TransactionManagerDIName).(database.TransactionManager)
			return websystem_usecase.NewServiceUsecase(serviceRepo, serviceZoneUc, txManager), nil
		},
	}

	serviceZoneDef := di.Def{
		Name:  ServiceZoneUsecaseDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get(ServiceZoneRepoDIName).(settingRepository.IServiceZoneRepository)
			txManager := ctn.Get(TransactionManagerDIName).(database.TransactionManager)
			return websystem_usecase.NewServiceZoneUsecase(repo, txManager), nil
		},
	}

	distancePricingRuleDef := di.Def{
		Name:  DistancePricingRuleUsecaseDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get(DistancePricingRuleRepoDIName).(settingRepository.IDistancePricingRuleRepository)
			serviceRepo := ctn.Get(ServiceRepoDIName).(settingRepository.IServiceRepository)
			txManager := ctn.Get(TransactionManagerDIName).(database.TransactionManager)
			return websystem_usecase.NewDistancePricingRuleUsecase(repo, serviceRepo, txManager), nil
		},
	}

	surchargeRuleDef := di.Def{
		Name:  SurchargeRuleUsecaseDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get(SurchargeRuleRepoDIName).(settingRepository.ISurchargeRuleRepository)
			serviceRepo := ctn.Get(ServiceRepoDIName).(settingRepository.IServiceRepository)
			zoneRepo := ctn.Get(ZoneRepoDIName).(account_repo.IZoneRepository)
			transactionManager := ctn.Get(TransactionManagerDIName).(database.TransactionManager)
			return websystem_usecase.NewSurchargeRuleUsecase(repo, serviceRepo, zoneRepo, transactionManager), nil
		},
	}

	packageSizePricingDef := di.Def{
		Name:  PackageSizePricingUsecaseDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get(PackageSizePricingRepoDIName).(settingRepository.IPackageSizePricingRepository)
			serviceRepo := ctn.Get(ServiceRepoDIName).(settingRepository.IServiceRepository)
			txManager := ctn.Get(TransactionManagerDIName).(database.TransactionManager)
			return websystem_usecase.NewPackageSizePricingUsecase(repo, serviceRepo, txManager), nil
		},
	}

	authAdminDef := di.Def{
		Name:  AuthAdminUsecaseDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			adminRepo := ctn.Get(SystemAdminRepoDIName).(settingRepository.ISystemAdminRepository)
			loginHistoryRepo := ctn.Get(SystemLoginHistoryRepoDIName).(settingRepository.ISystemLoginHistoryRepository)
			refreshTokenRepo := ctn.Get(SystemAdminRefreshTokenRepoDIName).(settingRepository.ISystemAdminRefreshTokenRepository)
			txManager := ctn.Get(TransactionManagerDIName).(database.TransactionManager)
			return websystem_usecase.NewAuthAdminUsecase(adminRepo, loginHistoryRepo, refreshTokenRepo, txManager), nil
		},
	}

	return builder.Add(
		userProfileDef,
		notifyDef,
		settingDef,
		otpDef,
		zoneDef,
		sidebarDef,
		serviceZoneDef,
		serviceDef,
		distancePricingRuleDef,
		surchargeRuleDef,
		packageSizePricingDef,
		authAdminDef,
	)
}
