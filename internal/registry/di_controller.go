package registry

import (
	otpcontroller "go-structure/internal/controller"
	controller "go-structure/internal/controller/app_user"
	usecase_pkg "go-structure/internal/usecase"
	usecase "go-structure/internal/usecase/app_user"

	"github.com/sarulabs/di"
)

const (
	UserProfileControllerDIName = "user_profile_controller_di"
	OTPControllerDIName         = "otp_controller_di"
)

func buildControllers() error {
	userProfileDef := di.Def{
		Name:  UserProfileControllerDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			uc := ctn.Get(UserProfileUsecaseDIName).(usecase.IUserProfileUsecase)
			return controller.NewUserProfileController(uc), nil
		},
	}

	otpDef := di.Def{
		Name:  OTPControllerDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			uc := ctn.Get(OTPUsecaseDIName).(usecase_pkg.IOTPUsecase)
			return otpcontroller.NewOTPController(uc), nil
		},
	}

	return builder.Add(userProfileDef, otpDef)
}
