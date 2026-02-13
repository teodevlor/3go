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
	SurchargeRuleController interface {
		Create(c *gin.Context) *common.ResponseData
		GetByID(c *gin.Context) *common.ResponseData
		List(c *gin.Context) *common.ResponseData
		Update(c *gin.Context) *common.ResponseData
		Delete(c *gin.Context) *common.ResponseData
	}

	surchargeRuleController struct {
		*controller.BaseController
		uc usecase.ISurchargeRuleUsecase
	}
)

func NewSurchargeRuleController(uc usecase.ISurchargeRuleUsecase) SurchargeRuleController {
	return &surchargeRuleController{
		BaseController: controller.NewBaseController(),
		uc:             uc,
	}
}

func (s *surchargeRuleController) Create(c *gin.Context) *common.ResponseData {
	var req dto.CreateSurchargeRuleRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}
	result, err := s.uc.Create(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, usecase.ErrServiceNotFound) {
			return common.ErrorResponse(common.StatusBadRequest, []string{err.Error()})
		}
		if errors.Is(err, usecase.ErrZoneNotFound) {
			return common.ErrorResponse(common.StatusBadRequest, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (s *surchargeRuleController) GetByID(c *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}
	result, err := s.uc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrSurchargeRuleNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (s *surchargeRuleController) List(c *gin.Context) *common.ResponseData {
	var serviceID, zoneID *uuid.UUID
	if s := c.Query("service_id"); s != "" {
		parsed, err := uuid.Parse(s)
		if err != nil {
			return common.ErrorResponse(common.StatusBadRequest, []string{"service_id không hợp lệ"})
		}
		serviceID = &parsed
	}
	if z := c.Query("zone_id"); z != "" {
		parsed, err := uuid.Parse(z)
		if err != nil {
			return common.ErrorResponse(common.StatusBadRequest, []string{"zone_id không hợp lệ"})
		}
		zoneID = &parsed
	}
	result, err := s.uc.List(c.Request.Context(), serviceID, zoneID)
	if err != nil {
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (s *surchargeRuleController) Update(c *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}
	var req dto.UpdateSurchargeRuleRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}
	result, err := s.uc.Update(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrSurchargeRuleNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		if errors.Is(err, usecase.ErrServiceNotFound) || errors.Is(err, usecase.ErrZoneNotFound) {
			return common.ErrorResponse(common.StatusBadRequest, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (s *surchargeRuleController) Delete(c *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}
	if err := s.uc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, usecase.ErrSurchargeRuleNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, gin.H{"message": "Đã xóa quy tắc phụ thu"})
}
