package registry

import (
	otpcontroller "go-structure/internal/controller"
	controller "go-structure/internal/controller/app_user"
	websystem_controller "go-structure/internal/controller/web_system"
	usecase_pkg "go-structure/internal/usecase"
	usecase "go-structure/internal/usecase/app_user"
	websystem_usecase "go-structure/internal/usecase/web_system"

	"github.com/sarulabs/di"
)

const (
	UserProfileControllerDIName         = "user_profile_controller_di"
	OTPControllerDIName                 = "otp_controller_di"
	ZoneControllerDIName                = "zone_controller_di"
	SidebarControllerDIName             = "sidebar_controller_di"
	ServiceControllerDIName             = "service_controller_di"
	DistancePricingRuleControllerDIName = "distance_pricing_rule_controller_di"
	SurchargeRuleControllerDIName       = "surcharge_rule_controller_di"
	PackageSizePricingControllerDIName  = "package_size_pricing_controller_di"
	AuthAdminControllerDIName           = "auth_admin_controller_di"
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

	zoneDef := di.Def{
		Name:  ZoneControllerDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			uc := ctn.Get(ZoneUsecaseDIName).(websystem_usecase.IZoneUsecase)
			return websystem_controller.NewZoneController(uc), nil
		},
	}

	sidebarDef := di.Def{
		Name:  SidebarControllerDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			uc := ctn.Get(SidebarUsecaseDIName).(websystem_usecase.ISidebarUsecase)
			return websystem_controller.NewSidebarController(uc), nil
		},
	}

	serviceDef := di.Def{
		Name:  ServiceControllerDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			uc := ctn.Get(ServiceUsecaseDIName).(websystem_usecase.IServiceUsecase)
			return websystem_controller.NewServiceController(uc), nil
		},
	}

	distancePricingRuleDef := di.Def{
		Name:  DistancePricingRuleControllerDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			uc := ctn.Get(DistancePricingRuleUsecaseDIName).(websystem_usecase.IDistancePricingRuleUsecase)
			return websystem_controller.NewDistancePricingRuleController(uc), nil
		},
	}

	surchargeRuleDef := di.Def{
		Name:  SurchargeRuleControllerDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			uc := ctn.Get(SurchargeRuleUsecaseDIName).(websystem_usecase.ISurchargeRuleUsecase)
			return websystem_controller.NewSurchargeRuleController(uc), nil
		},
	}

	packageSizePricingDef := di.Def{
		Name:  PackageSizePricingControllerDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			uc := ctn.Get(PackageSizePricingUsecaseDIName).(websystem_usecase.IPackageSizePricingUsecase)
			return websystem_controller.NewPackageSizePricingController(uc), nil
		},
	}

	authAdminDef := di.Def{
		Name:  AuthAdminControllerDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			uc := ctn.Get(AuthAdminUsecaseDIName).(websystem_usecase.IAuthAdminUsecase)
			return websystem_controller.NewAuthAdminController(uc), nil
		},
	}

	return builder.Add(
		userProfileDef,
		otpDef,
		zoneDef,
		sidebarDef,
		serviceDef,
		distancePricingRuleDef,
		surchargeRuleDef,
		packageSizePricingDef,
		authAdminDef,
	)
}
