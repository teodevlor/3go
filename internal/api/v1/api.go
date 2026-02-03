package v1

import (
	controller "go-structure/internal/controller/app_user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewApiV1(
	router *gin.Engine,
	userProfileController controller.UserProfileController,
) {
	apiV1 := router.Group("api/v1")
	{
		apiV1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "OK"})
		})

		// modules
		NewUserProfileApi(userProfileController).InitUserProfileApi(apiV1, userProfileController)
	}
}
