package web_system

import (
	"errors"
	"go-structure/internal/common"
	"go-structure/internal/controller"
	websystemdto "go-structure/internal/dto/web_system"
	websystemusecase "go-structure/internal/usecase/web_system"
	"go-structure/pkg/validator"

	"github.com/gin-gonic/gin"
)

type (
	AuthAdminController interface {
		LoginAdmin(c *gin.Context) *common.ResponseData
		RefreshToken(c *gin.Context) *common.ResponseData
	}

	authAdminController struct {
		*controller.BaseController
		authAdminUsecase websystemusecase.IAuthAdminUsecase
	}
)

func NewAuthAdminController(authAdminUsecase websystemusecase.IAuthAdminUsecase) AuthAdminController {
	return &authAdminController{
		BaseController:   controller.NewBaseController(),
		authAdminUsecase: authAdminUsecase,
	}
}

func (ctl *authAdminController) LoginAdmin(c *gin.Context) *common.ResponseData {
	var req websystemdto.AdminLoginRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}

	clientInfo := ctl.GetClientInfo(c)
	result, err := ctl.authAdminUsecase.LoginAdmin(c.Request.Context(), &req, clientInfo.IP, clientInfo.UserAgent)
	if err != nil {
		if errors.Is(err, websystemusecase.ErrAdminNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		if errors.Is(err, websystemusecase.ErrAdminInvalidPassword) {
			return common.ErrorResponse(common.StatusUnauthorized, []string{err.Error()})
		}
		if errors.Is(err, websystemusecase.ErrAdminNotActive) {
			return common.ErrorResponse(common.StatusForbidden, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (ctl *authAdminController) RefreshToken(c *gin.Context) *common.ResponseData {
	var req websystemdto.AdminRefreshTokenRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}

	clientInfo := ctl.GetClientInfo(c)
	result, err := ctl.authAdminUsecase.RefreshToken(c.Request.Context(), req.RefreshToken, clientInfo.IP, clientInfo.UserAgent)
	if err != nil {
		if errors.Is(err, websystemusecase.ErrInvalidRefreshToken) {
			return common.ErrorResponse(common.StatusUnauthorized, []string{err.Error()})
		}
		if errors.Is(err, websystemusecase.ErrAdminNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		if errors.Is(err, websystemusecase.ErrAdminNotActive) {
			return common.ErrorResponse(common.StatusForbidden, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}

	return common.SuccessResponse(common.StatusOK, result)
}
