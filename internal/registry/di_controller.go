package registry

import (
	controller "go-structure/internal/controller/app_user"
	usecase "go-structure/internal/usecase/app_user"

	"github.com/sarulabs/di"
)

const (
	UserProfileControllerDIName = "user_profile_controller_di"
)

func buildControllers() error {
	def := di.Def{
		Name:  UserProfileControllerDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			uc := ctn.Get(UserProfileUsecaseDIName).(usecase.IUserProfileUsecase)
			return controller.NewUserProfileController(uc), nil
		},
	}
	return builder.Add(def)
}
