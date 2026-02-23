package v1

import (
	websystemctl "go-structure/internal/controller/web_system"
	"go-structure/internal/constants"
	"go-structure/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	WebSystemApi interface {
		InitWebSystemApi(router *gin.RouterGroup,
			authMiddleware gin.HandlerFunc,
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
		)
	}

	webSystemApi struct {
		authAdminController           websystemctl.AuthAdminController
		zoneController                websystemctl.ZoneController
		sidebarController             websystemctl.SidebarController
		serviceController             websystemctl.ServiceController
		distancePricingRuleController websystemctl.DistancePricingRuleController
		surchargeRuleController       websystemctl.SurchargeRuleController
		packageSizePricingController  websystemctl.PackageSizePricingController
		roleController                websystemctl.RoleController
		adminController               websystemctl.AdminController
		permissionController          websystemctl.PermissionController
	}
)

const (
	API_MODULE_WEB_SYSTEM = "system"
)

func NewWebSystemApi(
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
) WebSystemApi {
	return &webSystemApi{
		authAdminController:           authAdminController,
		zoneController:                zoneController,
		sidebarController:             sidebarController,
		serviceController:             serviceController,
		distancePricingRuleController: distancePricingRuleController,
		surchargeRuleController:       surchargeRuleController,
		packageSizePricingController:  packageSizePricingController,
		roleController:                roleController,
		adminController:               adminController,
		permissionController:          permissionController,
	}
}

func (a *webSystemApi) InitWebSystemApi(
	router *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
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
) {
	systemRoutes := router.Group(API_MODULE_WEB_SYSTEM)
	{
		// Auth Admin (không cần đăng nhập)
		systemRoutes.POST("auth/login", func(c *gin.Context) {
			resp := authAdminController.LoginAdmin(c)
			c.JSON(http.StatusOK, resp)
		})
		systemRoutes.POST("auth/refresh", func(c *gin.Context) {
			resp := authAdminController.RefreshToken(c)
			c.JSON(http.StatusOK, resp)
		})

		protected := systemRoutes.Group("", authMiddleware)
		{
			// Zones
			protected.POST("zones", middleware.RequirePermission(permissionChecker, constants.PermissionZoneCreate), func(c *gin.Context) {
				resp := zoneController.CreateZone(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.GET("zones", middleware.RequirePermission(permissionChecker, constants.PermissionZoneList), func(c *gin.Context) {
				resp := zoneController.ListZones(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.GET("zones/:id", middleware.RequirePermission(permissionChecker, constants.PermissionZoneRead), func(c *gin.Context) {
				resp := zoneController.GetZone(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.PUT("zones/:id", middleware.RequirePermission(permissionChecker, constants.PermissionZoneUpdate), func(c *gin.Context) {
				resp := zoneController.UpdateZone(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.DELETE("zones/:id", middleware.RequirePermission(permissionChecker, constants.PermissionZoneDelete), func(c *gin.Context) {
				resp := zoneController.DeleteZone(c)
				c.JSON(http.StatusOK, resp)
			})

			// Services
			protected.POST("services", middleware.RequirePermission(permissionChecker, constants.PermissionServiceCreate), func(c *gin.Context) {
				resp := serviceController.CreateService(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.GET("services", middleware.RequirePermission(permissionChecker, constants.PermissionServiceList), func(c *gin.Context) {
				resp := serviceController.ListServices(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.GET("services/:id", middleware.RequirePermission(permissionChecker, constants.PermissionServiceRead), func(c *gin.Context) {
				resp := serviceController.GetService(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.PUT("services/:id", middleware.RequirePermission(permissionChecker, constants.PermissionServiceUpdate), func(c *gin.Context) {
				resp := serviceController.UpdateService(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.DELETE("services/:id", middleware.RequirePermission(permissionChecker, constants.PermissionServiceDelete), func(c *gin.Context) {
				resp := serviceController.DeleteService(c)
				c.JSON(http.StatusOK, resp)
			})

			// Distance pricing rules (quy tắc giá theo km)
			protected.GET("distance-pricing-rules", middleware.RequirePermission(permissionChecker, constants.PermissionDistancePricingRuleList), func(c *gin.Context) {
				resp := distancePricingRuleController.List(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.POST("distance-pricing-rules", middleware.RequirePermission(permissionChecker, constants.PermissionDistancePricingRuleCreate), func(c *gin.Context) {
				resp := distancePricingRuleController.Create(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.GET("distance-pricing-rules/:id", middleware.RequirePermission(permissionChecker, constants.PermissionDistancePricingRuleRead), func(c *gin.Context) {
				resp := distancePricingRuleController.GetByID(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.PUT("distance-pricing-rules/:id", middleware.RequirePermission(permissionChecker, constants.PermissionDistancePricingRuleUpdate), func(c *gin.Context) {
				resp := distancePricingRuleController.Update(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.DELETE("distance-pricing-rules/:id", middleware.RequirePermission(permissionChecker, constants.PermissionDistancePricingRuleDelete), func(c *gin.Context) {
				resp := distancePricingRuleController.Delete(c)
				c.JSON(http.StatusOK, resp)
			})

			// Surcharge rules (quy tắc phụ thu)
			protected.GET("surcharge-rules", middleware.RequirePermission(permissionChecker, constants.PermissionSurchargeRuleList), func(c *gin.Context) {
				resp := surchargeRuleController.List(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.POST("surcharge-rules", middleware.RequirePermission(permissionChecker, constants.PermissionSurchargeRuleCreate), func(c *gin.Context) {
				resp := surchargeRuleController.Create(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.GET("surcharge-rules/:id", middleware.RequirePermission(permissionChecker, constants.PermissionSurchargeRuleRead), func(c *gin.Context) {
				resp := surchargeRuleController.GetByID(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.PUT("surcharge-rules/:id", middleware.RequirePermission(permissionChecker, constants.PermissionSurchargeRuleUpdate), func(c *gin.Context) {
				resp := surchargeRuleController.Update(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.DELETE("surcharge-rules/:id", middleware.RequirePermission(permissionChecker, constants.PermissionSurchargeRuleDelete), func(c *gin.Context) {
				resp := surchargeRuleController.Delete(c)
				c.JSON(http.StatusOK, resp)
			})

			// Package size pricing (quy tắc giá theo kích thước gói)
			protected.GET("package-size-pricings", middleware.RequirePermission(permissionChecker, constants.PermissionPackageSizePricingList), func(c *gin.Context) {
				resp := packageSizePricingController.List(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.POST("package-size-pricings", middleware.RequirePermission(permissionChecker, constants.PermissionPackageSizePricingCreate), func(c *gin.Context) {
				resp := packageSizePricingController.Create(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.GET("package-size-pricings/:id", middleware.RequirePermission(permissionChecker, constants.PermissionPackageSizePricingRead), func(c *gin.Context) {
				resp := packageSizePricingController.GetByID(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.PUT("package-size-pricings/:id", middleware.RequirePermission(permissionChecker, constants.PermissionPackageSizePricingUpdate), func(c *gin.Context) {
				resp := packageSizePricingController.Update(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.DELETE("package-size-pricings/:id", middleware.RequirePermission(permissionChecker, constants.PermissionPackageSizePricingDelete), func(c *gin.Context) {
				resp := packageSizePricingController.Delete(c)
				c.JSON(http.StatusOK, resp)
			})

			// Sidebars
			protected.POST("sidebars", middleware.RequirePermission(permissionChecker, constants.PermissionSidebarCreate), func(c *gin.Context) {
				resp := sidebarController.CreateSidebar(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.GET("sidebars", middleware.RequirePermission(permissionChecker, constants.PermissionSidebarList), func(c *gin.Context) {
				resp := sidebarController.ListSidebars(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.GET("sidebars/:id", middleware.RequirePermission(permissionChecker, constants.PermissionSidebarRead), func(c *gin.Context) {
				resp := sidebarController.GetSidebar(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.PUT("sidebars/:id", middleware.RequirePermission(permissionChecker, constants.PermissionSidebarUpdate), func(c *gin.Context) {
				resp := sidebarController.UpdateSidebar(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.DELETE("sidebars/:id", middleware.RequirePermission(permissionChecker, constants.PermissionSidebarDelete), func(c *gin.Context) {
				resp := sidebarController.DeleteSidebar(c)
				c.JSON(http.StatusOK, resp)
			})

			// Roles (vai trò)
			protected.POST("roles", middleware.RequirePermission(permissionChecker, constants.PermissionRoleCreate), func(c *gin.Context) {
				resp := roleController.Create(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.GET("roles", middleware.RequirePermission(permissionChecker, constants.PermissionRoleList), func(c *gin.Context) {
				resp := roleController.List(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.GET("roles/:id", middleware.RequirePermission(permissionChecker, constants.PermissionRoleRead), func(c *gin.Context) {
				resp := roleController.GetByID(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.PUT("roles/:id", middleware.RequirePermission(permissionChecker, constants.PermissionRoleUpdate), func(c *gin.Context) {
				resp := roleController.Update(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.DELETE("roles/:id", middleware.RequirePermission(permissionChecker, constants.PermissionRoleDelete), func(c *gin.Context) {
				resp := roleController.Delete(c)
				c.JSON(http.StatusOK, resp)
			})

			// Admins (quản trị viên, kèm role_ids)
			protected.POST("admins", middleware.RequirePermission(permissionChecker, constants.PermissionAdminCreate), func(c *gin.Context) {
				resp := adminController.Create(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.GET("admins", middleware.RequirePermission(permissionChecker, constants.PermissionAdminList), func(c *gin.Context) {
				resp := adminController.List(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.GET("admins/:id", middleware.RequirePermission(permissionChecker, constants.PermissionAdminRead), func(c *gin.Context) {
				resp := adminController.GetByID(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.PUT("admins/:id", middleware.RequirePermission(permissionChecker, constants.PermissionAdminUpdate), func(c *gin.Context) {
				resp := adminController.Update(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.DELETE("admins/:id", middleware.RequirePermission(permissionChecker, constants.PermissionAdminDelete), func(c *gin.Context) {
				resp := adminController.Delete(c)
				c.JSON(http.StatusOK, resp)
			})

			// Permissions (quyền)
			protected.POST("permissions", middleware.RequirePermission(permissionChecker, constants.PermissionPermissionCreate), func(c *gin.Context) {
				resp := permissionController.Create(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.GET("permissions", middleware.RequirePermission(permissionChecker, constants.PermissionPermissionList), func(c *gin.Context) {
				resp := permissionController.List(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.GET("permissions/:id", middleware.RequirePermission(permissionChecker, constants.PermissionPermissionRead), func(c *gin.Context) {
				resp := permissionController.GetByID(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.PUT("permissions/:id", middleware.RequirePermission(permissionChecker, constants.PermissionPermissionUpdate), func(c *gin.Context) {
				resp := permissionController.Update(c)
				c.JSON(http.StatusOK, resp)
			})
			protected.DELETE("permissions/:id", middleware.RequirePermission(permissionChecker, constants.PermissionPermissionDelete), func(c *gin.Context) {
				resp := permissionController.Delete(c)
				c.JSON(http.StatusOK, resp)
			})
		}
	}
}
