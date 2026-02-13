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
	DistancePricingRuleController interface {
		Create(c *gin.Context) *common.ResponseData
		GetByID(c *gin.Context) *common.ResponseData
		List(c *gin.Context) *common.ResponseData
		Update(c *gin.Context) *common.ResponseData
		Delete(c *gin.Context) *common.ResponseData
	}

	distancePricingRuleController struct {
		*controller.BaseController
		uc usecase.IDistancePricingRuleUsecase
	}
)

func NewDistancePricingRuleController(uc usecase.IDistancePricingRuleUsecase) DistancePricingRuleController {
	return &distancePricingRuleController{
		BaseController: controller.NewBaseController(),
		uc:             uc,
	}
}

func (dpr *distancePricingRuleController) Create(c *gin.Context) *common.ResponseData {
	var req dto.CreateDistancePricingRuleRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}

	result, err := dpr.uc.Create(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, usecase.ErrServiceNotFound) {
			return common.ErrorResponse(common.StatusBadRequest, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (dpr *distancePricingRuleController) GetByID(c *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}

	result, err := dpr.uc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrDistancePricingRuleNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (dpr *distancePricingRuleController) List(c *gin.Context) *common.ResponseData {
	var serviceID *uuid.UUID
	if serviceIDStr := c.Query("service_id"); serviceIDStr != "" {
		parsedID, err := uuid.Parse(serviceIDStr)
		if err != nil {
			return common.ErrorResponse(common.StatusBadRequest, []string{"service_id không hợp lệ"})
		}
		serviceID = &parsedID
	}

	result, err := dpr.uc.List(c.Request.Context(), serviceID)
	if err != nil {
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (dpr *distancePricingRuleController) Update(c *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}

	var req dto.UpdateDistancePricingRuleRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}

	result, err := dpr.uc.Update(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrDistancePricingRuleNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		if errors.Is(err, usecase.ErrServiceNotFound) {
			return common.ErrorResponse(common.StatusBadRequest, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (dpr *distancePricingRuleController) Delete(c *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}

	if err := dpr.uc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, usecase.ErrDistancePricingRuleNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, gin.H{"message": "Đã xóa quy tắc giá theo khoảng cách"})
}
