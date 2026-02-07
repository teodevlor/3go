package v1

import (
	controller "go-structure/internal/controller/app_user"
	middlewares "go-structure/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	UserProfileApi interface {
		InitUserProfileApi(router *gin.RouterGroup, userProfileController controller.UserProfileController)
	}

	userProfileApi struct {
		userProfileController controller.UserProfileController
	}
)

const (
	API_MODULE_USER = "auth/user"
)

func NewUserProfileApi(userProfileController controller.UserProfileController) UserProfileApi {
	return &userProfileApi{userProfileController: userProfileController}
}

func (u *userProfileApi) InitUserProfileApi(router *gin.RouterGroup, userProfileController controller.UserProfileController) {
	userRoutesPublic := router.Group(API_MODULE_USER)
	{
		userRoutesPublic.POST("register", func(c *gin.Context) {
			resp := userProfileController.RegisterUserProfile(c)
			c.JSON(http.StatusOK, resp)
		})
		userRoutesPublic.POST("login", func(c *gin.Context) {
			resp := userProfileController.LoginUserProfile(c)
			c.JSON(http.StatusOK, resp)
		})
		userRoutesPublic.POST("active", func(c *gin.Context) {
			resp := userProfileController.ActiveUserProfile(c)
			c.JSON(http.StatusOK, resp)
		})
		userRoutesPublic.POST("refresh-token", func(c *gin.Context) {
			resp := userProfileController.RefreshToken(c)
			c.JSON(http.StatusOK, resp)
		})
		userRoutesPublic.POST("forgot-password", func(c *gin.Context) {
			resp := userProfileController.ForgotPassword(c)
			c.JSON(http.StatusOK, resp)
		})
		userRoutesPublic.POST("reset-password", func(c *gin.Context) {
			resp := userProfileController.ResetPassword(c)
			c.JSON(http.StatusOK, resp)
		})
	}

	userRoutesPrivate := router.Group(API_MODULE_USER)
	userRoutesPrivate.Use(middlewares.AuthMiddleware())
	{
		userRoutesPrivate.GET("profile", func(c *gin.Context) {
			resp := userProfileController.GetUserProfile(c)
			c.JSON(http.StatusOK, resp)
		})
		userRoutesPrivate.PUT("update-user-profile", func(c *gin.Context) {
			resp := userProfileController.UpdateUserProfile(c)
			c.JSON(http.StatusOK, resp)
		})
		userRoutesPrivate.POST("logout", func(c *gin.Context) {
			resp := userProfileController.Logout(c)
			c.JSON(http.StatusOK, resp)
		})
	}
}
