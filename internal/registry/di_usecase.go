package registry

import (
	"go-structure/internal/adapter"
	"go-structure/internal/helper/database"
	account_repo "go-structure/internal/repository"
	user_profile_repo "go-structure/internal/repository/app_user"
	usecase_pkg "go-structure/internal/usecase"
	user_profile_usecase "go-structure/internal/usecase/app_user"

	"github.com/sarulabs/di"
)

const (
	UserProfileUsecaseDIName = "user_profile_usecase_di"
	NotifyUsecaseDIName      = "notify_usecase_di"
	SettingUsecaseDIName     = "setting_usecase_di"
	OTPUsecaseDIName         = "otp_usecase_di"
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
			settingRepo := ctn.Get(SettingRepoDIName).(account_repo.ISettingRepository)
			return usecase_pkg.NewSettingUsecase(settingRepo), nil
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

	return builder.Add(userProfileDef, notifyDef, settingDef, otpDef)
}
