package v1

import (
	"go-structure/internal/constants"
	app_driver_controller "go-structure/internal/controller/app_driver"
	websystemctl "go-structure/internal/controller/web_system"
	"go-structure/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	API_MODULE_APP_DRIVER = "driver"
)

type (
	AppDriverApi interface {
		InitAppDriverApi(
			router *gin.RouterGroup,
			authMiddleware gin.HandlerFunc,
			permissionChecker middleware.AdminPermissionChecker,
		)
	}

	appDriverApi struct {
		driverDocumentTypeController app_driver_controller.DriverDocumentTypeController
		driverProfileController      app_driver_controller.DriverProfileController
		driverDocumentController     app_driver_controller.DriverDocumentController
		zoneController               websystemctl.ZoneController
	}
)

func NewAppDriverApi(
	driverDocumentTypeController app_driver_controller.DriverDocumentTypeController,
	driverProfileController app_driver_controller.DriverProfileController,
	driverDocumentController app_driver_controller.DriverDocumentController,
	zoneController websystemctl.ZoneController,
) AppDriverApi {
	return &appDriverApi{
		driverDocumentTypeController: driverDocumentTypeController,
		driverProfileController:      driverProfileController,
		driverDocumentController:     driverDocumentController,
		zoneController:               zoneController,
	}
}

func (a *appDriverApi) InitAppDriverApi(
	router *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	permissionChecker middleware.AdminPermissionChecker,
) {
	publicDriver := router.Group("public/" + API_MODULE_APP_DRIVER)
	{
		publicDriver.GET("document-types/required", func(c *gin.Context) {
			resp := a.driverDocumentTypeController.GetRequiredDocuments(c)
			c.JSON(http.StatusOK, resp)
		})
		publicDriver.GET("services", func(c *gin.Context) {
			resp := a.zoneController.ListZones(c)
			c.JSON(http.StatusOK, resp)
		})
		publicDriver.POST("documents/bulk", func(c *gin.Context) {
			resp := a.driverDocumentController.BulkCreate(c)
			c.JSON(http.StatusOK, resp)
		})
	}

	authDriver := router.Group("auth/" + API_MODULE_APP_DRIVER)
	{
		authDriver.POST("register", func(c *gin.Context) {
			resp := a.driverProfileController.RegisterDriver(c)
			c.JSON(http.StatusOK, resp)
		})
		authDriver.POST("verify-otp", func(c *gin.Context) {
			resp := a.driverProfileController.VerifyDriverOtp(c)
			c.JSON(http.StatusOK, resp)
		})
		authDriver.POST("login", func(c *gin.Context) {
			resp := a.driverProfileController.LoginDriver(c)
			c.JSON(http.StatusOK, resp)
		})
	}

	driverRoute := router.Group(API_MODULE_APP_DRIVER, middleware.AuthMiddleware())
	{
		driverRoute.POST("online", func(c *gin.Context) {
			resp := a.driverProfileController.GoOnline(c)
			c.JSON(http.StatusOK, resp)
		})
		driverRoute.POST("offline", func(c *gin.Context) {
			resp := a.driverProfileController.GoOffline(c)
			c.JSON(http.StatusOK, resp)
		})
		driverRoute.POST("ping", func(c *gin.Context) {
			resp := a.driverProfileController.PingOnline(c)
			c.JSON(http.StatusOK, resp)
		})
	}

	driverRoutes := router.Group(API_MODULE_APP_DRIVER)
	protected := driverRoutes.Group("", authMiddleware)
	{
		// Driver document types
		protected.POST("document-types", middleware.RequirePermission(permissionChecker, constants.PermissionDriverDocumentTypeCreate), func(c *gin.Context) {
			resp := a.driverDocumentTypeController.Create(c)
			c.JSON(http.StatusOK, resp)
		})
		protected.GET("document-types", middleware.RequirePermission(permissionChecker, constants.PermissionDriverDocumentTypeList), func(c *gin.Context) {
			resp := a.driverDocumentTypeController.List(c)
			c.JSON(http.StatusOK, resp)
		})
		protected.GET("document-types/:id", middleware.RequirePermission(permissionChecker, constants.PermissionDriverDocumentTypeRead), func(c *gin.Context) {
			resp := a.driverDocumentTypeController.GetByID(c)
			c.JSON(http.StatusOK, resp)
		})
		protected.PUT("document-types/:id", middleware.RequirePermission(permissionChecker, constants.PermissionDriverDocumentTypeUpdate), func(c *gin.Context) {
			resp := a.driverDocumentTypeController.Update(c)
			c.JSON(http.StatusOK, resp)
		})
		protected.DELETE("document-types/:id", middleware.RequirePermission(permissionChecker, constants.PermissionDriverDocumentTypeDelete), func(c *gin.Context) {
			resp := a.driverDocumentTypeController.Delete(c)
			c.JSON(http.StatusOK, resp)
		})

		// Driver documents: CRUD + bulk create + bulk update (PATCH)
		protected.POST("documents", middleware.RequirePermission(permissionChecker, constants.PermissionDriverDocumentCreate), func(c *gin.Context) {
			resp := a.driverDocumentController.Create(c)
			c.JSON(http.StatusOK, resp)
		})
		protected.GET("documents", middleware.RequirePermission(permissionChecker, constants.PermissionDriverDocumentList), func(c *gin.Context) {
			resp := a.driverDocumentController.ListByDriverID(c)
			c.JSON(http.StatusOK, resp)
		})
		protected.GET("documents/:id", middleware.RequirePermission(permissionChecker, constants.PermissionDriverDocumentRead), func(c *gin.Context) {
			resp := a.driverDocumentController.GetByID(c)
			c.JSON(http.StatusOK, resp)
		})
		protected.PUT("documents/:id", middleware.RequirePermission(permissionChecker, constants.PermissionDriverDocumentUpdate), func(c *gin.Context) {
			resp := a.driverDocumentController.Update(c)
			c.JSON(http.StatusOK, resp)
		})
		protected.PATCH("documents/bulk", func(c *gin.Context) {
			resp := a.driverDocumentController.BulkUpdate(c)
			c.JSON(http.StatusOK, resp)
		})
		protected.DELETE("documents/:id", middleware.RequirePermission(permissionChecker, constants.PermissionDriverDocumentDelete), func(c *gin.Context) {
			resp := a.driverDocumentController.Delete(c)
			c.JSON(http.StatusOK, resp)
		})
	}
}
