package appuser

import (
	"errors"
	"go-structure/internal/common"
	"go-structure/internal/controller"
	dto "go-structure/internal/dto/app_user"
	usecase "go-structure/internal/usecase/app_user"
	"go-structure/pkg/validator"

	"github.com/gin-gonic/gin"
)

type (
	UserProfileController interface {
		RegisterUserProfile(c *gin.Context) *common.ResponseData
		ActiveUserProfile(c *gin.Context) *common.ResponseData
		LoginUserProfile(c *gin.Context) *common.ResponseData
		GetUserProfile(c *gin.Context) *common.ResponseData
		RefreshToken(c *gin.Context) *common.ResponseData
		Logout(c *gin.Context) *common.ResponseData
		UpdateUserProfile(c *gin.Context) *common.ResponseData
		ForgotPassword(c *gin.Context) *common.ResponseData
		ResetPassword(c *gin.Context) *common.ResponseData
	}

	userProfileController struct {
		*controller.BaseController
		userProfileUsecase usecase.IUserProfileUsecase
	}
)

func NewUserProfileController(userProfileUsecase usecase.IUserProfileUsecase) UserProfileController {
	return &userProfileController{
		BaseController:     controller.NewBaseController(),
		userProfileUsecase: userProfileUsecase,
	}
}

func (ctl *userProfileController) RegisterUserProfile(c *gin.Context) *common.ResponseData {
	var req dto.UserRegisterRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}

	if err := validator.ValidatePassword(req.Password); err != nil {
		return common.ErrorResponse(common.StatusUnprocessableEntity, []string{err.Error()})
	}
	result, err := ctl.userProfileUsecase.RegisterUserProfile(c.Request.Context(), &req)
	if err != nil {
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (ctl *userProfileController) ActiveUserProfile(c *gin.Context) *common.ResponseData {
	var req dto.UserActiveRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}
	clientInfo := ctl.GetClientInfo(c)
	verified, err := ctl.userProfileUsecase.ActiveUserProfile(c.Request.Context(), req.Phone, req.Code, clientInfo.IP, clientInfo.UserAgent)
	if err != nil {
		if errors.Is(err, usecase.ErrUserAlreadyActive) {
			return common.ErrorResponse(common.StatusBadRequest, []string{err.Error()})
		}
		if errors.Is(err, usecase.ErrInvalidOTP) {
			return common.ErrorResponse(common.StatusBadRequest, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, verified)
}

func (ctl *userProfileController) LoginUserProfile(c *gin.Context) *common.ResponseData {
	var req dto.UserLoginRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}

	clientInfo := ctl.GetClientInfo(c)
	result, err := ctl.userProfileUsecase.LoginUserProfile(c.Request.Context(), &req, clientInfo.IP, clientInfo.UserAgent)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		if errors.Is(err, usecase.ErrInvalidPassword) {
			return common.ErrorResponse(common.StatusUnauthorized, []string{err.Error()})
		}
		if errors.Is(err, usecase.ErrUserNotActive) {
			return common.ErrorResponse(common.StatusForbidden, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (ctl *userProfileController) GetUserProfile(c *gin.Context) *common.ResponseData {
	accountID, errResp := ctl.GetAccountIDFromContext(c)
	if errResp != nil {
		return errResp
	}

	profile, err := ctl.userProfileUsecase.GetUserProfile(c.Request.Context(), accountID)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}

	return common.SuccessResponse(common.StatusOK, profile)
}

func (ctl *userProfileController) RefreshToken(c *gin.Context) *common.ResponseData {
	var req dto.RefreshTokenRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}

	result, err := ctl.userProfileUsecase.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidRefreshToken) {
			return common.ErrorResponse(common.StatusUnauthorized, []string{err.Error()})
		}
		if errors.Is(err, usecase.ErrUserNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		if errors.Is(err, usecase.ErrUserNotActive) {
			return common.ErrorResponse(common.StatusForbidden, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}

	return common.SuccessResponse(common.StatusOK, result)
}

func (ctl *userProfileController) Logout(c *gin.Context) *common.ResponseData {
	accountID, errResp := ctl.GetAccountIDFromContext(c)
	if errResp != nil {
		return errResp
	}

	err := ctl.userProfileUsecase.Logout(c.Request.Context(), accountID)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}

	return common.SuccessResponse(common.StatusOK, dto.LogoutResponseDto{
		UserMessage: common.BaseMessageLogoutSuccess,
	})
}

func (ctl *userProfileController) UpdateUserProfile(c *gin.Context) *common.ResponseData {
	accountID, errResp := ctl.GetAccountIDFromContext(c)
	if errResp != nil {
		return errResp
	}

	var req dto.UpdateUserProfile
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}

	result, err := ctl.userProfileUsecase.UpdateUserProfile(c.Request.Context(), accountID, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}

	return common.SuccessResponse(common.StatusOK, result)
}

func (ctl *userProfileController) ForgotPassword(c *gin.Context) *common.ResponseData {
	var req dto.ForgotPasswordRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}

	result, err := ctl.userProfileUsecase.ForgotPassword(c.Request.Context(), req.Phone)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		var retryErr *common.ErrorWithRetryAfter
		if errors.As(err, &retryErr) {
			return common.ErrorResponse(common.StatusTooManyRequests, []string{retryErr.Error()})
		}
		if errors.Is(err, common.ErrResendTooSoon) {
			return common.ErrorResponse(common.StatusTooManyRequests, []string{err.Error()})
		}
		if errors.Is(err, common.ErrResendMaxExceeded) {
			return common.ErrorResponse(common.StatusTooManyRequests, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}

	return common.SuccessResponse(common.StatusOK, result)
}

func (ctl *userProfileController) ResetPassword(c *gin.Context) *common.ResponseData {
	var req dto.ResetPasswordRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}

	if err := validator.ValidatePassword(req.NewPassword); err != nil {
		return common.ErrorResponse(common.StatusUnprocessableEntity, []string{err.Error()})
	}

	if req.NewPassword != req.ConfirmPassword {
		return common.ErrorResponse(common.StatusUnprocessableEntity, []string{"confirm_password: Mật khẩu xác nhận không khớp"})
	}

	clientInfo := ctl.GetClientInfo(c)
	result, err := ctl.userProfileUsecase.ResetPassword(c.Request.Context(), req.Phone, req.Code, req.NewPassword, clientInfo.IP, clientInfo.UserAgent)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		if errors.Is(err, usecase.ErrInvalidOTP) {
			return common.ErrorResponse(common.StatusBadRequest, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}

	return common.SuccessResponse(common.StatusOK, result)
}
