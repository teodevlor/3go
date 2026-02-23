package v1

import (
	ctl "go-structure/internal/controller"
	controller "go-structure/internal/controller/app_user"
	websystemctl "go-structure/internal/controller/web_system"
	"go-structure/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewApiV1(
	router *gin.Engine,
	userProfileController controller.UserProfileController,
	otpController ctl.OTPController,
	permissionChecker middleware.AdminPermissionChecker,
	authAdminController websystemctl.AuthAdminController,
	zoneController websystemctl.ZoneController,
	sidebarController websystemctl.SidebarController,
	serviceController websystemctl.ServiceController,
	distancePricingRuleController websystemctl.DistancePricingRuleController,
	surchargeRuleController websystemctl.SurchargeRuleController,
	packageSizePricingController websystemctl.PackageSizePricingController,
	roleController websystemctl.RoleController,
	adminController websystemctl.AdminController,
	permissionController websystemctl.PermissionController,
	storageController ctl.StorageController,
) {
	apiV1 := router.Group("api/v1")
	{
		apiV1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "OK bay giờ test commit"})
		})

		// modules
		NewUserProfileApi(userProfileController).InitUserProfileApi(apiV1, userProfileController)
		NewOTPApi(otpController).InitOTPApI(apiV1, otpController)
		NewStorageApi(storageController).InitStorageApi(apiV1, storageController)
		NewWebSystemApi(
			authAdminController,
			zoneController,
			sidebarController,
			serviceController,
			distancePricingRuleController,
			surchargeRuleController,
			packageSizePricingController,
			roleController,
			adminController,
			permissionController,
		).InitWebSystemApi(apiV1,
			middleware.AdminAuthMiddleware(),
			permissionChecker,
			authAdminController,
			zoneController,
			sidebarController,
			serviceController,
			distancePricingRuleController,
			surchargeRuleController,
			packageSizePricingController,
			roleController,
			adminController,
			permissionController,
		)
	}
}
