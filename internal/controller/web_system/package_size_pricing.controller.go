package web_system

import (
	"errors"

	"go-structure/internal/common"
	"go-structure/internal/controller"
	dto "go-structure/internal/dto/web_system"
	usecase "go-structure/internal/usecase/web_system"
	"go-structure/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	PackageSizePricingController interface {
		Create(c *gin.Context) *common.ResponseData
		GetByID(c *gin.Context) *common.ResponseData
		List(c *gin.Context) *common.ResponseData
		Update(c *gin.Context) *common.ResponseData
		Delete(c *gin.Context) *common.ResponseData
	}

	packageSizePricingController struct {
		*controller.BaseController
		uc usecase.IPackageSizePricingUsecase
	}
)

func NewPackageSizePricingController(uc usecase.IPackageSizePricingUsecase) PackageSizePricingController {
	return &packageSizePricingController{
		BaseController: controller.NewBaseController(),
		uc:             uc,
	}
}

func (c *packageSizePricingController) Create(ctx *gin.Context) *common.ResponseData {
	var req dto.CreatePackageSizePricingRequestDto
	if err := ctx.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}
	result, err := c.uc.Create(ctx.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, usecase.ErrServiceNotFound) {
			return common.ErrorResponse(common.StatusBadRequest, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (c *packageSizePricingController) GetByID(ctx *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}
	result, err := c.uc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrPackageSizePricingNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (c *packageSizePricingController) List(ctx *gin.Context) *common.ResponseData {
	var serviceID *uuid.UUID
	if s := ctx.Query("service_id"); s != "" {
		parsed, err := uuid.Parse(s)
		if err != nil {
			return common.ErrorResponse(common.StatusBadRequest, []string{"service_id không hợp lệ"})
		}
		serviceID = &parsed
	}
	result, err := c.uc.List(ctx.Request.Context(), serviceID)
	if err != nil {
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (c *packageSizePricingController) Update(ctx *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}
	var req dto.UpdatePackageSizePricingRequestDto
	if err := ctx.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}
	result, err := c.uc.Update(ctx.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrPackageSizePricingNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		if errors.Is(err, usecase.ErrServiceNotFound) {
			return common.ErrorResponse(common.StatusBadRequest, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (c *packageSizePricingController) Delete(ctx *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}
	if err := c.uc.Delete(ctx.Request.Context(), id); err != nil {
		if errors.Is(err, usecase.ErrPackageSizePricingNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, gin.H{"message": "Đã xóa quy tắc giá theo kích thước gói"})
}
