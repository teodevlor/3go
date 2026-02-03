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
	}

	userRoutesPrivate := router.Group(API_MODULE_USER)
	userRoutesPrivate.Use(middlewares.AuthMiddleware())
	{
		userRoutesPrivate.GET("me", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "user info (authenticated)",
				"user":    c.MustGet("userId"),
			})
		})
	}
}
