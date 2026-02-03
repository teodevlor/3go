package registry

import (
	"go-structure/internal/adapter"
	account_repo "go-structure/internal/repository"
	user_profile_repo "go-structure/internal/repository/app_user"
	notify_usecase "go-structure/internal/usecase"
	user_profile_usecase "go-structure/internal/usecase/app_user"

	"github.com/sarulabs/di"
)

const (
	UserProfileUsecaseDIName = "user_profile_usecase_di"
	NotifyUsecaseDIName      = "notify_usecase_di"
)

func buildUsecases() error {
	userProfileDef := di.Def{
		Name:  UserProfileUsecaseDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			usrProfileRepo := ctn.Get(UserProfileRepoDIName).(user_profile_repo.IUserProfileRepository)
			accountRepo := ctn.Get(AccountRepoDIName).(account_repo.IAccountRepository)
			notifyUcAny := ctn.Get(NotifyUsecaseDIName)

			var notifyUc notify_usecase.INotifyUsecase
			if notifyUcAny != nil {
				notifyUc = notifyUcAny.(notify_usecase.INotifyUsecase)
			}
			return user_profile_usecase.NewUserProfileUsecase(usrProfileRepo, accountRepo, notifyUc), nil
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
			return notify_usecase.NewNotifyUsecase(ad.(adapter.INotifyAdapter)), nil
		},
	}

	return builder.Add(userProfileDef, notifyDef)
}
