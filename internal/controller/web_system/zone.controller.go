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
	ZoneController interface {
		CreateZone(c *gin.Context) *common.ResponseData
		GetZone(c *gin.Context) *common.ResponseData
		ListZones(c *gin.Context) *common.ResponseData
		UpdateZone(c *gin.Context) *common.ResponseData
		DeleteZone(c *gin.Context) *common.ResponseData
	}

	zoneController struct {
		*controller.BaseController
		zoneUsecase usecase.IZoneUsecase
	}
)

func NewZoneController(zoneUsecase usecase.IZoneUsecase) ZoneController {
	return &zoneController{
		BaseController: controller.NewBaseController(),
		zoneUsecase:    zoneUsecase,
	}
}

func (zc *zoneController) CreateZone(c *gin.Context) *common.ResponseData {
	var req dto.CreateZoneRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}

	if err := req.Polygon.Validate(); err != nil {
		return common.ErrorResponse(common.StatusUnprocessableEntity, []string{err.Error()})
	}

	result, err := zc.zoneUsecase.CreateZone(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, usecase.ErrZoneCodeExists) {
			return common.ErrorResponse(common.StatusUnprocessableEntity, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}

	return common.SuccessResponse(common.StatusOK, result)
}

func (zc *zoneController) GetZone(c *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}
	result, err := zc.zoneUsecase.GetZone(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrZoneNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (zc *zoneController) ListZones(c *gin.Context) *common.ResponseData {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.DefaultQuery("search", "")
	result, err := zc.zoneUsecase.ListZones(c.Request.Context(), page, limit, search)
	if err != nil {
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (zc *zoneController) UpdateZone(c *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}
	var req dto.UpdateZoneRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}
	if err := req.Polygon.Validate(); err != nil {
		return common.ErrorResponse(common.StatusUnprocessableEntity, []string{err.Error()})
	}
	result, err := zc.zoneUsecase.UpdateZone(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrZoneNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (zc *zoneController) DeleteZone(c *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}
	if err := zc.zoneUsecase.DeleteZone(c.Request.Context(), id); err != nil {
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, gin.H{"message": "Đã xóa zone"})
}
