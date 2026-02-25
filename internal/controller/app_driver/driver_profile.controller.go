package app_driver

import (
	"errors"

	"go-structure/internal/common"
	"go-structure/internal/controller"
	dto "go-structure/internal/dto/app_driver"
	"go-structure/internal/middleware"
	usecase "go-structure/internal/usecase/app_driver"
	"go-structure/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	DriverProfileController interface {
		RegisterDriver(c *gin.Context) *common.ResponseData
		VerifyDriverOtp(c *gin.Context) *common.ResponseData
		LoginDriver(c *gin.Context) *common.ResponseData

		GoOnline(c *gin.Context) *common.ResponseData
		GoOffline(c *gin.Context) *common.ResponseData
		PingOnline(c *gin.Context) *common.ResponseData
	}

	driverProfileController struct {
		*controller.BaseController
		uc usecase.IDriverProfileUsecase
	}
)

func NewDriverProfileController(uc usecase.IDriverProfileUsecase) DriverProfileController {
	return &driverProfileController{
		BaseController: controller.NewBaseController(),
		uc:             uc,
	}
}

func (ctrl *driverProfileController) RegisterDriver(c *gin.Context) *common.ResponseData {
	var req dto.DriverRegisterRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}
	if err := validator.ValidatePassword(req.Password); err != nil {
		return common.ErrorResponse(common.StatusUnprocessableEntity, []string{err.Error()})
	}
	result, err := ctrl.uc.RegisterDriver(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, usecase.ErrDriverAlreadyRegistered) {
			return common.ErrorResponse(common.StatusBadRequest, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (ctrl *driverProfileController) VerifyDriverOtp(c *gin.Context) *common.ResponseData {
	var req dto.DriverVerifyOtpRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}
	clientInfo := ctrl.GetClientInfo(c)
	result, err := ctrl.uc.VerifyDriverOtp(c.Request.Context(), req.Phone, req.Code, clientInfo.IP, clientInfo.UserAgent)
	if err != nil {
		if errors.Is(err, usecase.ErrDriverNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		if errors.Is(err, usecase.ErrDriverInvalidOTP) {
			return common.ErrorResponse(common.StatusBadRequest, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (ctrl *driverProfileController) LoginDriver(c *gin.Context) *common.ResponseData {
	var req dto.DriverLoginRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}
	clientInfo := ctrl.GetClientInfo(c)
	result, err := ctrl.uc.LoginDriver(c.Request.Context(), &req, clientInfo.IP, clientInfo.UserAgent)
	if err != nil {
		if errors.Is(err, usecase.ErrDriverNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		if errors.Is(err, usecase.ErrDriverInvalidPassword) {
			return common.ErrorResponse(common.StatusBadRequest, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (ctrl *driverProfileController) GoOnline(c *gin.Context) *common.ResponseData {
	var req dto.DriverLocationStatusRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}

	accountIDVal, exists := c.Get(middleware.ContextAccountIDKey)
	if !exists {
		return common.ErrorResponse(common.StatusUnauthorized, []string{"Unauthorized"})
	}
	accountID, ok := accountIDVal.(uuid.UUID)
	if !ok {
		return common.ErrorResponse(common.StatusUnauthorized, []string{"Invalid account in context"})
	}

	if err := ctrl.uc.GoOnline(c.Request.Context(), accountID, &req); err != nil {
		if errors.Is(err, usecase.ErrDriverNotFound) || errors.Is(err, usecase.ErrDriverNotActive) {
			return common.ErrorResponse(common.StatusBadRequest, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}

	return common.SuccessResponse(common.StatusOK, map[string]string{
		"message": "Driver is online",
	})
}

func (ctrl *driverProfileController) GoOffline(c *gin.Context) *common.ResponseData {
	accountIDVal, exists := c.Get(middleware.ContextAccountIDKey)
	if !exists {
		return common.ErrorResponse(common.StatusUnauthorized, []string{"Unauthorized"})
	}
	accountID, ok := accountIDVal.(uuid.UUID)
	if !ok {
		return common.ErrorResponse(common.StatusUnauthorized, []string{"Invalid account in context"})
	}

	if err := ctrl.uc.GoOffline(c.Request.Context(), accountID); err != nil {
		if errors.Is(err, usecase.ErrDriverNotFound) {
			return common.ErrorResponse(common.StatusBadRequest, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}

	return common.SuccessResponse(common.StatusOK, map[string]string{
		"message": "Driver is offline",
	})
}

func (ctrl *driverProfileController) PingOnline(c *gin.Context) *common.ResponseData {
	var req dto.DriverLocationStatusRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}

	accountIDVal, exists := c.Get(middleware.ContextAccountIDKey)
	if !exists {
		return common.ErrorResponse(common.StatusUnauthorized, []string{"Unauthorized"})
	}
	accountID, ok := accountIDVal.(uuid.UUID)
	if !ok {
		return common.ErrorResponse(common.StatusUnauthorized, []string{"Invalid account in context"})
	}

	if err := ctrl.uc.PingOnline(c.Request.Context(), accountID, &req); err != nil {
		if errors.Is(err, usecase.ErrDriverNotFound) || errors.Is(err, usecase.ErrDriverNotActive) {
			return common.ErrorResponse(common.StatusBadRequest, []string{err.Error()})
		}
		if errors.Is(err, usecase.ErrDriverPingTooFrequent) {
			return common.ErrorResponse(common.StatusTooManyRequests, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}

	return common.SuccessResponse(common.StatusOK, map[string]string{
		"message": "Driver ping updated",
	})
}
