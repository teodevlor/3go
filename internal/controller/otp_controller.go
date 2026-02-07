package controller

import (
	"errors"
	"fmt"

	"go-structure/internal/common"
	"go-structure/internal/dto"
	"go-structure/internal/usecase"
	"go-structure/pkg/validator"

	"github.com/gin-gonic/gin"
)

type (
	OTPController interface {
		ResendOTP(c *gin.Context) *common.ResponseData
	}

	otpController struct {
		*BaseController
		otpUsecase usecase.IOTPUsecase
	}
)

func NewOTPController(otpUsecase usecase.IOTPUsecase) OTPController {
	return &otpController{
		BaseController: NewBaseController(),
		otpUsecase:     otpUsecase,
	}
}

func (ctl *otpController) ResendOTP(c *gin.Context) *common.ResponseData {
	var req dto.OTPResendRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}

	_, err := ctl.otpUsecase.ResendOTP(c.Request.Context(), req.Target, req.Purpose)
	if err != nil {
		if errors.Is(err, common.ErrResendTooSoon) || errors.Is(err, common.ErrResendMaxExceeded) {
			errMsg := err.Error()
			var withRetry *common.ErrorWithRetryAfter
			if errors.As(err, &withRetry) && withRetry.RetryAfterSeconds > 0 {
				errMsg = fmt.Sprintf("%s%s", errMsg, fmt.Sprintf(common.BaseMessageResendOTPMaxExceededRetryAfter, withRetry.RetryAfterSeconds))
			}
			return common.ErrorResponse(common.StatusTooManyRequests, []string{errMsg})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}

	return common.SuccessResponse(common.StatusOK, dto.OTPResendResponseDto{
		UserMessage: common.AppUserMessageResendOTPSuccess,
	})
}
