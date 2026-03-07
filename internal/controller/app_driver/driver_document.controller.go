package app_driver

import (
	"errors"

	"go-structure/internal/common"
	"go-structure/internal/controller"
	dto "go-structure/internal/dto/app_driver"
	usecase "go-structure/internal/usecase/app_driver"
	"go-structure/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	DriverDocumentController interface {
		Create(c *gin.Context) *common.ResponseData
		BulkCreate(c *gin.Context) *common.ResponseData
		GetByID(c *gin.Context) *common.ResponseData
		ListByDriverID(c *gin.Context) *common.ResponseData
		Update(c *gin.Context) *common.ResponseData
		UpdateStatus(c *gin.Context) *common.ResponseData
		BulkUpdate(c *gin.Context) *common.ResponseData
		Delete(c *gin.Context) *common.ResponseData
	}

	driverDocumentController struct {
		*controller.BaseController
		uc usecase.IDriverDocumentUsecase
	}
)

func NewDriverDocumentController(uc usecase.IDriverDocumentUsecase) DriverDocumentController {
	return &driverDocumentController{
		BaseController: controller.NewBaseController(),
		uc:             uc,
	}
}

func (ctrl *driverDocumentController) Create(c *gin.Context) *common.ResponseData {
	var req dto.CreateDriverDocumentRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}
	result, err := ctrl.uc.Create(c.Request.Context(), &req)
	if err != nil {
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (ctrl *driverDocumentController) BulkCreate(c *gin.Context) *common.ResponseData {
	var req dto.BulkCreateDriverDocumentsRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}
	result, err := ctrl.uc.BulkCreate(c.Request.Context(), &req)
	if err != nil {
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (ctrl *driverDocumentController) GetByID(c *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}
	result, err := ctrl.uc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrDriverDocumentNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (ctrl *driverDocumentController) ListByDriverID(c *gin.Context) *common.ResponseData {
	driverIDStr := c.Param("driver_id")
	if driverIDStr == "" {
		driverIDStr = c.Param("id")
	}
	if driverIDStr == "" {
		return common.ErrorResponse(common.StatusBadRequest, []string{"driver_id là bắt buộc"})
	}
	driverID, err := uuid.Parse(driverIDStr)
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"driver_id không hợp lệ"})
	}
	result, err := ctrl.uc.ListByDriverID(c.Request.Context(), driverID)
	if err != nil {
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (ctrl *driverDocumentController) Update(c *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}
	var req dto.UpdateDriverDocumentRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}
	result, err := ctrl.uc.Update(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrDriverDocumentNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (ctrl *driverDocumentController) UpdateStatus(c *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}

	var req dto.UpdateDriverDocumentStatusRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}

	status := req.Status
	updateReq := dto.BulkUpdateDriverDocumentsRequestDto{
		Items: []dto.BulkUpdateDriverDocumentItemDto{
			{
				ID:     id.String(),
				Status: &status,
			},
		},
	}

	result, err := ctrl.uc.BulkUpdate(c.Request.Context(), &updateReq)
	if err != nil {
		if errors.Is(err, usecase.ErrDriverDocumentNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	if len(result.Items) == 0 {
		return common.ErrorResponse(common.StatusInternalServerError, []string{"cập nhật trạng thái thất bại"})
	}
	return common.SuccessResponse(common.StatusOK, result.Items[0])
}

func (ctrl *driverDocumentController) BulkUpdate(c *gin.Context) *common.ResponseData {
	var req dto.BulkUpdateDriverDocumentsRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := validator.Translate(err)
		return common.ErrorResponse(common.StatusUnprocessableEntity, msgs)
	}
	result, err := ctrl.uc.BulkUpdate(c.Request.Context(), &req)
	if err != nil {
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, result)
}

func (ctrl *driverDocumentController) Delete(c *gin.Context) *common.ResponseData {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return common.ErrorResponse(common.StatusBadRequest, []string{"id không hợp lệ"})
	}
	if err := ctrl.uc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, usecase.ErrDriverDocumentNotFound) {
			return common.ErrorResponse(common.StatusNotFound, []string{err.Error()})
		}
		return common.ErrorResponse(common.StatusInternalServerError, []string{err.Error()})
	}
	return common.SuccessResponse(common.StatusOK, gin.H{"message": "Đã xóa giấy tờ"})
}
