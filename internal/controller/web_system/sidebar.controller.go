package web_system

import (
	"errors"
	"strconv"

	"go-structure/internal/common"
	"go-structure/internal/controller"
	dto "go-structure/internal/dto/web_system"
	usecase "go-structure/internal/usecase/web_system"
	"go-structure/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	SidebarController interface {
		CreateSidebar(c *gin.Context) *common.ResponseData
		GetSidebar(c *gin.Context) *common.ResponseData
		ListSidebars(c *gin.Context) *common.ResponseData
		UpdateSidebar(c *gin.Context) *common.ResponseData
		DeleteSidebar(c *gin.Context) *common.ResponseData
	}

	sidebarController struct {
		*controller.BaseController
		sidebarUsecase usecase.ISidebarUsecase
	}
)

func NewSidebarController(sidebarUsecase usecase.ISidebarUsecase) SidebarController {
	return &sidebarController{
		BaseController: controller.NewBaseController(),
		sidebarUsecase: sidebarUsecase,
	}
}

func (sc *sidebarController) CreateSidebar(c *gin.Context) *common.ResponseData {
	var req dto.CreateSidebarRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}
	result, err := sc.sidebarUsecase.CreateSidebar(c.Request.Context(), &req)
	if err != nil {
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (sc *sidebarController) GetSidebar(c *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}
	result, err := sc.sidebarUsecase.GetSidebar(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrSidebarNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (sc *sidebarController) ListSidebars(c *gin.Context) *common.ResponseData {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	contextFilter := c.DefaultQuery("context", "")
	result, err := sc.sidebarUsecase.ListSidebars(c.Request.Context(), contextFilter, page, limit)
	if err != nil {
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (sc *sidebarController) UpdateSidebar(c *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}
	var req dto.UpdateSidebarRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}
	result, err := sc.sidebarUsecase.UpdateSidebar(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrSidebarNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (sc *sidebarController) DeleteSidebar(c *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}
	if err := sc.sidebarUsecase.DeleteSidebar(c.Request.Context(), id); err != nil {
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, gin.H{"message": "Đã xóa sidebar"})
}

