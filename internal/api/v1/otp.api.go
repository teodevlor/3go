package v1

import (
	otpcontroller "go-structure/internal/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	OTPApi interface {
		InitOTPApI(router *gin.RouterGroup, otpController otpcontroller.OTPController)
	}

	otpApi struct {
		otpController otpcontroller.OTPController
	}
)

const (
	API_MODULE_OTP = "otp"
)

func NewOTPApi(otpController otpcontroller.OTPController) OTPApi {
	return &otpApi{otpController: otpController}
}

func (a *otpApi) InitOTPApI(router *gin.RouterGroup, otpController otpcontroller.OTPController) {
	otpRoutes := router.Group(API_MODULE_OTP)
	{
		otpRoutes.POST("resend", func(c *gin.Context) {
			resp := otpController.ResendOTP(c)
			c.JSON(http.StatusOK, resp)
		})
	}
}
