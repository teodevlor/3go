package appuser

import (
	"go-structure/internal/common"
	dto "go-structure/internal/dto/app_user"
	usecase "go-structure/internal/usecase/app_user"
	"go-structure/pkg/validator"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
)

type (
	UserProfileController interface {
		RegisterUserProfile(c *gin.Context) *common.ResponseData
	}

	userProfileController struct {
		userProfileUsecase usecase.IUserProfileUsecase
	}
)

func NewUserProfileController(userProfileUsecase usecase.IUserProfileUsecase) UserProfileController {
	return &userProfileController{userProfileUsecase: userProfileUsecase}
}

func (ctl *userProfileController) RegisterUserProfile(c *gin.Context) *common.ResponseData {
	var req dto.UserRegisterRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msg := validator.TranslateMessage(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, []string{msg}, "")
	}
	logger.Info("RegisterUserProfile", "phone", req.Phone)
	if err := validator.ValidatePassword(req.Password); err != nil {
		return common.ErrorResponse(common.StatusUnprocessableEntity, []string{err.Error()}, "")
	}
	result, err := ctl.userProfileUsecase.RegisterUserProfile(c.Request.Context(), &req)
	if err != nil {
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()}, "")
	}
	return common.SuccessResponse(common.StatusOK, result, "User profile registered successfully")
}
