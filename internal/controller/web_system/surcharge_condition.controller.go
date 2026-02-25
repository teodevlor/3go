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
	SurchargeConditionController interface {
		Create(c *gin.Context) *common.ResponseData
		GetByID(c *gin.Context) *common.ResponseData
		List(c *gin.Context) *common.ResponseData
		Update(c *gin.Context) *common.ResponseData
		Delete(c *gin.Context) *common.ResponseData
	}

	surchargeConditionController struct {
		*controller.BaseController
		uc usecase.ISurchargeConditionUsecase
	}
)

func NewSurchargeConditionController(uc usecase.ISurchargeConditionUsecase) SurchargeConditionController {
	return &surchargeConditionController{
		BaseController: controller.NewBaseController(),
		uc:             uc,
	}
}

func (s *surchargeConditionController) Create(c *gin.Context) *common.ResponseData {
	var req dto.CreateSurchargeConditionRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}
	result, err := s.uc.Create(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, usecase.ErrSurchargeConditionCodeExists) {
			return common.ErrorResponse(common.StatusUnprocessableEntity, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (s *surchargeConditionController) GetByID(c *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}
	result, err := s.uc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrSurchargeConditionNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (s *surchargeConditionController) List(c *gin.Context) *common.ResponseData {
	result, err := s.uc.List(c.Request.Context())
	if err != nil {
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (s *surchargeConditionController) Update(c *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}
	var req dto.UpdateSurchargeConditionRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}
	result, err := s.uc.Update(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrSurchargeConditionNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		if errors.Is(err, usecase.ErrSurchargeConditionCodeExists) {
			return common.ErrorResponse(common.StatusUnprocessableEntity, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (s *surchargeConditionController) Delete(c *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}
	if err := s.uc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, usecase.ErrSurchargeConditionNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, gin.H{"message": "Đã xóa điều kiện phụ thu"})
}

