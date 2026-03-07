package registry

import (
	"go-structure/config"
	app_driver_controller "go-structure/internal/controller/app_driver"
	otpcontroller "go-structure/internal/controller"
	controller "go-structure/internal/controller/app_user"
	websystem_controller "go-structure/internal/controller/web_system"
	usecase_pkg "go-structure/internal/usecase"
	app_driver_usecase "go-structure/internal/usecase/app_driver"
	usecase "go-structure/internal/usecase/app_user"
	websystem_usecase "go-structure/internal/usecase/web_system"

	"github.com/sarulabs/di"
)

const (
	UserProfileControllerDIName         = "user_profile_controller_di"
	OTPControllerDIName                 = "otp_controller_di"
	StorageControllerDIName             = "storage_controller_di"
	ZoneControllerDIName                = "zone_controller_di"
	SidebarControllerDIName             = "sidebar_controller_di"
	ServiceControllerDIName             = "service_controller_di"
	DistancePricingRuleControllerDIName = "distance_pricing_rule_controller_di"
	SurchargeConditionControllerDIName  = "surcharge_condition_controller_di"
	SurchargeRuleControllerDIName       = "surcharge_rule_controller_di"
	PackageSizePricingControllerDIName  = "package_size_pricing_controller_di"
	AuthAdminControllerDIName           = "auth_admin_controller_di"
	RoleControllerDIName                = "role_controller_di"
	AdminControllerDIName               = "admin_controller_di"
	PermissionControllerDIName          = "permission_controller_di"
	DriverDocumentTypeControllerDIName  = "driver_document_type_controller_di"
	DriverProfileControllerDIName       = "driver_profile_controller_di"
	DriverDocumentControllerDIName      = "driver_document_controller_di"
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

	storageDef := di.Def{
		Name:  StorageControllerDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			uc := ctn.Get(StorageUsecaseDIName).(usecase_pkg.IStorageUsecase)
			driverDocumentUc := ctn.Get(DriverDocumentUsecaseDIName).(app_driver_usecase.IDriverDocumentUsecase)
			cfg := ctn.Get(ConfigDIName).(*config.Config)
			return otpcontroller.NewStorageController(uc, driverDocumentUc, cfg.Storage), nil
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

	surchargeConditionDef := di.Def{
		Name:  SurchargeConditionControllerDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			uc := ctn.Get(SurchargeConditionUsecaseDIName).(websystem_usecase.ISurchargeConditionUsecase)
			return websystem_controller.NewSurchargeConditionController(uc), nil
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

	roleDef := di.Def{
		Name:  RoleControllerDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			uc := ctn.Get(RoleUsecaseDIName).(websystem_usecase.IRoleUsecase)
			return websystem_controller.NewRoleController(uc), nil
		},
	}

	adminDef := di.Def{
		Name:  AdminControllerDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			uc := ctn.Get(AdminUsecaseDIName).(websystem_usecase.IAdminUsecase)
			return websystem_controller.NewAdminController(uc), nil
		},
	}

	permissionDef := di.Def{
		Name:  PermissionControllerDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			uc := ctn.Get(PermissionUsecaseDIName).(websystem_usecase.IPermissionUsecase)
			return websystem_controller.NewPermissionController(uc), nil
		},
	}

	driverDocumentTypeDef := di.Def{
		Name:  DriverDocumentTypeControllerDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			uc := ctn.Get(DriverDocumentTypeUsecaseDIName).(app_driver_usecase.IDriverDocumentTypeUsecase)
			return app_driver_controller.NewDriverDocumentTypeController(uc), nil
		},
	}

	driverProfileDef := di.Def{
		Name:  DriverProfileControllerDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			uc := ctn.Get(DriverProfileUsecaseDIName).(app_driver_usecase.IDriverProfileUsecase)
			return app_driver_controller.NewDriverProfileController(uc), nil
		},
	}

	driverDocumentDef := di.Def{
		Name:  DriverDocumentControllerDIName,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			uc := ctn.Get(DriverDocumentUsecaseDIName).(app_driver_usecase.IDriverDocumentUsecase)
			return app_driver_controller.NewDriverDocumentController(uc), nil
		},
	}

	return builder.Add(
		userProfileDef,
		otpDef,
		storageDef,
		zoneDef,
		sidebarDef,
		serviceDef,
		distancePricingRuleDef,
		surchargeConditionDef,
		surchargeRuleDef,
		packageSizePricingDef,
		authAdminDef,
		roleDef,
		adminDef,
		permissionDef,
		driverDocumentTypeDef,
		driverProfileDef,
		driverDocumentDef,
	)
}
