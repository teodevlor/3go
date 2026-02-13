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
	ServiceController interface {
		CreateService(c *gin.Context) *common.ResponseData
		GetService(c *gin.Context) *common.ResponseData
		ListServices(c *gin.Context) *common.ResponseData
		UpdateService(c *gin.Context) *common.ResponseData
		DeleteService(c *gin.Context) *common.ResponseData
	}

	serviceController struct {
		*controller.BaseController
		serviceUsecase usecase.IServiceUsecase
	}
)

func NewServiceController(serviceUsecase usecase.IServiceUsecase) ServiceController {
	return &serviceController{
		BaseController: controller.NewBaseController(),
		serviceUsecase: serviceUsecase,
	}
}

func (su *serviceController) CreateService(c *gin.Context) *common.ResponseData {
	var req dto.CreateServiceRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}

	result, err := su.serviceUsecase.CreateService(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, usecase.ErrServiceCodeExists) {
			return common.ErrorResponse(common.StatusUnprocessableEntity, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (su *serviceController) GetService(c *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}
	result, err := su.serviceUsecase.GetService(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrServiceNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (su *serviceController) ListServices(c *gin.Context) *common.ResponseData {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.DefaultQuery("search", "")
	result, err := su.serviceUsecase.ListServices(c.Request.Context(), page, limit, search)
	if err != nil {
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (su *serviceController) UpdateService(c *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}
	var req dto.UpdateServiceRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}
	result, err := su.serviceUsecase.UpdateService(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrServiceNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (su *serviceController) DeleteService(c *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}
	if err := su.serviceUsecase.DeleteService(c.Request.Context(), id); err != nil {
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, gin.H{"message": "Đã xóa dịch vụ"})
}
