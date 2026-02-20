package v1

import (
	otpcontroller "go-structure/internal/controller"
	controller "go-structure/internal/controller/app_user"
	websystemctl "go-structure/internal/controller/web_system"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewApiV1(
	router *gin.Engine,
	userProfileController controller.UserProfileController,
	otpController otpcontroller.OTPController,
	authAdminController websystemctl.AuthAdminController,
	zoneController websystemctl.ZoneController,
	sidebarController websystemctl.SidebarController,
	serviceController websystemctl.ServiceController,
	distancePricingRuleController websystemctl.DistancePricingRuleController,
	surchargeRuleController websystemctl.SurchargeRuleController,
	packageSizePricingController websystemctl.PackageSizePricingController,
) {
	apiV1 := router.Group("api/v1")
	{
		apiV1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "OK bay giờ test commit"})
		})

		// modules
		NewUserProfileApi(userProfileController).InitUserProfileApi(apiV1, userProfileController)
		NewOTPApi(otpController).InitOTPApI(apiV1, otpController)
		NewWebSystemApi(
			authAdminController,
			zoneController,
			sidebarController,
			serviceController,
			distancePricingRuleController,
			surchargeRuleController,
			packageSizePricingController,
		).InitWebSystemApi(apiV1,
			authAdminController,
			zoneController,
			sidebarController,
			serviceController,
			distancePricingRuleController,
			surchargeRuleController,
			packageSizePricingController,
		)
	}
}
