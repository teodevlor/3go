package app_driver

import (
	"errors"
	"strconv"

	"go-structure/internal/common"
	"go-structure/internal/controller"
	dto "go-structure/internal/dto/app_driver"
	usecase "go-structure/internal/usecase/app_driver"
	"go-structure/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	DriverDocumentTypeController interface {
		Create(c *gin.Context) *common.ResponseData
		GetByID(c *gin.Context) *common.ResponseData
		List(c *gin.Context) *common.ResponseData
		GetRequiredDocuments(c *gin.Context) *common.ResponseData
		Update(c *gin.Context) *common.ResponseData
		Delete(c *gin.Context) *common.ResponseData
	}

	driverDocumentTypeController struct {
		*controller.BaseController
		uc usecase.IDriverDocumentTypeUsecase
	}
)

func NewDriverDocumentTypeController(uc usecase.IDriverDocumentTypeUsecase) DriverDocumentTypeController {
	return &driverDocumentTypeController{
		BaseController: controller.NewBaseController(),
		uc:              uc,
	}
}

func (ctrl *driverDocumentTypeController) Create(c *gin.Context) *common.ResponseData {
	var req dto.CreateDriverDocumentTypeRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}
	result, err := ctrl.uc.Create(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, usecase.ErrDriverDocumentTypeCodeExists) {
			return common.ErrorResponse(common.StatusUnprocessableEntity, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (ctrl *driverDocumentTypeController) GetByID(c *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}
	result, err := ctrl.uc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrDriverDocumentTypeNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (ctrl *driverDocumentTypeController) List(c *gin.Context) *common.ResponseData {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.DefaultQuery("search", "")
	serviceID := c.Query("service_id")
	var serviceIDPtr *string
	if serviceID != "" {
		serviceIDPtr = &serviceID
	}
	result, err := ctrl.uc.List(c.Request.Context(), page, limit, search, serviceIDPtr)
	if err != nil {
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (ctrl *driverDocumentTypeController) GetRequiredDocuments(c *gin.Context) *common.ResponseData {
	serviceIDStr := c.Query("service_id")
	if serviceIDStr == "" {
		return common.ErrorResponse(common.StatusBadRequest, []string{"service_id là bắt buộc"})
	}
	serviceID, err := uuid.Parse(serviceIDStr)
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"service_id không hợp lệ"})
	}
	result, err := ctrl.uc.GetRequiredByServiceID(c.Request.Context(), serviceID)
	if err != nil {
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (ctrl *driverDocumentTypeController) Update(c *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}
	var req dto.UpdateDriverDocumentTypeRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}
	result, err := ctrl.uc.Update(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrDriverDocumentTypeNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		if errors.Is(err, usecase.ErrDriverDocumentTypeCodeExists) {
			return common.ErrorResponse(common.StatusUnprocessableEntity, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (ctrl *driverDocumentTypeController) Delete(c *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}
	if err := ctrl.uc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, usecase.ErrDriverDocumentTypeNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, gin.H{"message": "Đã xóa loại giấy tờ"})
}
