package v1

import (
	websystemctl "go-structure/internal/controller/web_system"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	WebSystemApi interface {
		InitWebSystemApi(router *gin.RouterGroup, authAdminController websystemctl.AuthAdminController, zoneController websystemctl.ZoneController, sidebarController websystemctl.SidebarController, serviceController websystemctl.ServiceController, distancePricingRuleController websystemctl.DistancePricingRuleController, surchargeRuleController websystemctl.SurchargeRuleController, packageSizePricingController websystemctl.PackageSizePricingController)
	}

	webSystemApi struct {
		authAdminController           websystemctl.AuthAdminController
		zoneController                websystemctl.ZoneController
		sidebarController             websystemctl.SidebarController
		serviceController             websystemctl.ServiceController
		distancePricingRuleController websystemctl.DistancePricingRuleController
		surchargeRuleController       websystemctl.SurchargeRuleController
		packageSizePricingController  websystemctl.PackageSizePricingController
	}
)

const (
	API_MODULE_WEB_SYSTEM = "system"
)

func NewWebSystemApi(authAdminController websystemctl.AuthAdminController, zoneController websystemctl.ZoneController, sidebarController websystemctl.SidebarController, serviceController websystemctl.ServiceController, distancePricingRuleController websystemctl.DistancePricingRuleController, surchargeRuleController websystemctl.SurchargeRuleController, packageSizePricingController websystemctl.PackageSizePricingController) WebSystemApi {
	return &webSystemApi{
		authAdminController:           authAdminController,
		zoneController:                zoneController,
		sidebarController:             sidebarController,
		serviceController:             serviceController,
		distancePricingRuleController: distancePricingRuleController,
		surchargeRuleController:       surchargeRuleController,
		packageSizePricingController:  packageSizePricingController,
	}
}

func (a *webSystemApi) InitWebSystemApi(
	router *gin.RouterGroup,
	authAdminController websystemctl.AuthAdminController,
	zoneController websystemctl.ZoneController,
	sidebarController websystemctl.SidebarController,
	serviceController websystemctl.ServiceController,
	distancePricingRuleController websystemctl.DistancePricingRuleController,
	surchargeRuleController websystemctl.SurchargeRuleController,
	packageSizePricingController websystemctl.PackageSizePricingController,
) {
	systemRoutes := router.Group(API_MODULE_WEB_SYSTEM)
	{
		// Auth Admin
		systemRoutes.POST("auth/login", func(c *gin.Context) {
			resp := authAdminController.LoginAdmin(c)
			c.JSON(http.StatusOK, resp)
		})
		systemRoutes.POST("auth/refresh", func(c *gin.Context) {
			resp := authAdminController.RefreshToken(c)
			c.JSON(http.StatusOK, resp)
		})

		// Zones
		systemRoutes.POST("zones", func(c *gin.Context) {
			resp := zoneController.CreateZone(c)
			c.JSON(http.StatusOK, resp)
		})
		systemRoutes.GET("zones", func(c *gin.Context) {
			resp := zoneController.ListZones(c)
			c.JSON(http.StatusOK, resp)
		})
		systemRoutes.GET("zones/:id", func(c *gin.Context) {
			resp := zoneController.GetZone(c)
			c.JSON(http.StatusOK, resp)
		})
		systemRoutes.PUT("zones/:id", func(c *gin.Context) {
			resp := zoneController.UpdateZone(c)
			c.JSON(http.StatusOK, resp)
		})
		systemRoutes.DELETE("zones/:id", func(c *gin.Context) {
			resp := zoneController.DeleteZone(c)
			c.JSON(http.StatusOK, resp)
		})

		// Services
		systemRoutes.POST("services", func(c *gin.Context) {
			resp := serviceController.CreateService(c)
			c.JSON(http.StatusOK, resp)
		})
		systemRoutes.GET("services", func(c *gin.Context) {
			resp := serviceController.ListServices(c)
			c.JSON(http.StatusOK, resp)
		})
		systemRoutes.GET("services/:id", func(c *gin.Context) {
			resp := serviceController.GetService(c)
			c.JSON(http.StatusOK, resp)
		})
		systemRoutes.PUT("services/:id", func(c *gin.Context) {
			resp := serviceController.UpdateService(c)
			c.JSON(http.StatusOK, resp)
		})
		systemRoutes.DELETE("services/:id", func(c *gin.Context) {
			resp := serviceController.DeleteService(c)
			c.JSON(http.StatusOK, resp)
		})

		// Distance pricing rules (quy tắc giá theo km)
		systemRoutes.GET("distance-pricing-rules", func(c *gin.Context) {
			resp := distancePricingRuleController.List(c)
			c.JSON(http.StatusOK, resp)
		})
		systemRoutes.POST("distance-pricing-rules", func(c *gin.Context) {
			resp := distancePricingRuleController.Create(c)
			c.JSON(http.StatusOK, resp)
		})
		systemRoutes.GET("distance-pricing-rules/:id", func(c *gin.Context) {
			resp := distancePricingRuleController.GetByID(c)
			c.JSON(http.StatusOK, resp)
		})
		systemRoutes.PUT("distance-pricing-rules/:id", func(c *gin.Context) {
			resp := distancePricingRuleController.Update(c)
			c.JSON(http.StatusOK, resp)
		})
		systemRoutes.DELETE("distance-pricing-rules/:id", func(c *gin.Context) {
			resp := distancePricingRuleController.Delete(c)
			c.JSON(http.StatusOK, resp)
		})

		// Surcharge rules (quy tắc phụ thu)
		systemRoutes.GET("surcharge-rules", func(c *gin.Context) {
			resp := surchargeRuleController.List(c)
			c.JSON(http.StatusOK, resp)
		})
		systemRoutes.POST("surcharge-rules", func(c *gin.Context) {
			resp := surchargeRuleController.Create(c)
			c.JSON(http.StatusOK, resp)
		})
		systemRoutes.GET("surcharge-rules/:id", func(c *gin.Context) {
			resp := surchargeRuleController.GetByID(c)
			c.JSON(http.StatusOK, resp)
		})
		systemRoutes.PUT("surcharge-rules/:id", func(c *gin.Context) {
			resp := surchargeRuleController.Update(c)
			c.JSON(http.StatusOK, resp)
		})
		systemRoutes.DELETE("surcharge-rules/:id", func(c *gin.Context) {
			resp := surchargeRuleController.Delete(c)
			c.JSON(http.StatusOK, resp)
		})

		// Package size pricing (quy tắc giá theo kích thước gói)
		systemRoutes.GET("package-size-pricings", func(c *gin.Context) {
			resp := packageSizePricingController.List(c)
			c.JSON(http.StatusOK, resp)
		})
		systemRoutes.POST("package-size-pricings", func(c *gin.Context) {
			resp := packageSizePricingController.Create(c)
			c.JSON(http.StatusOK, resp)
		})
		systemRoutes.GET("package-size-pricings/:id", func(c *gin.Context) {
			resp := packageSizePricingController.GetByID(c)
			c.JSON(http.StatusOK, resp)
		})
		systemRoutes.PUT("package-size-pricings/:id", func(c *gin.Context) {
			resp := packageSizePricingController.Update(c)
			c.JSON(http.StatusOK, resp)
		})
		systemRoutes.DELETE("package-size-pricings/:id", func(c *gin.Context) {
			resp := packageSizePricingController.Delete(c)
			c.JSON(http.StatusOK, resp)
		})

		// Sidebars
		systemRoutes.POST("sidebars", func(c *gin.Context) {
			resp := sidebarController.CreateSidebar(c)
			c.JSON(http.StatusOK, resp)
		})
		systemRoutes.GET("sidebars", func(c *gin.Context) {
			resp := sidebarController.ListSidebars(c)
			c.JSON(http.StatusOK, resp)
		})
		systemRoutes.GET("sidebars/:id", func(c *gin.Context) {
			resp := sidebarController.GetSidebar(c)
			c.JSON(http.StatusOK, resp)
		})
		systemRoutes.PUT("sidebars/:id", func(c *gin.Context) {
			resp := sidebarController.UpdateSidebar(c)
			c.JSON(http.StatusOK, resp)
		})
		systemRoutes.DELETE("sidebars/:id", func(c *gin.Context) {
			resp := sidebarController.DeleteSidebar(c)
			c.JSON(http.StatusOK, resp)
		})
	}
}
