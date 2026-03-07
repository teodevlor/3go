package app_driver

import (
	"context"
	"errors"
	"strings"

	"go-structure/global"
	"go-structure/internal/constants"
	dto "go-structure/internal/dto/app_driver"
	dto_common "go-structure/internal/dto/common"
	"go-structure/internal/helper/database"
	"go-structure/internal/middleware"
	appdriverrepo "go-structure/internal/repository/app_driver"
	appdrivermodel "go-structure/internal/repository/model/app_driver"
	appdrivertransformer "go-structure/internal/transformer/app_driver"
	websystem_repo "go-structure/internal/repository/web_system"
	pgdb "go-structure/orm/db/postgres"
	"go-structure/pkg/validator"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

var (
	ErrDriverDocumentTypeNotFound   = errors.New("không tìm thấy loại giấy tờ")
	ErrDriverDocumentTypeCodeExists = errors.New("mã loại giấy tờ đã tồn tại cho phạm vi này (chung hoặc theo dịch vụ)")
	ErrServiceNotFoundForDocType    = errors.New("không tìm thấy dịch vụ")
)

type (
	IDriverDocumentTypeUsecase interface {
		Create(ctx context.Context, req *dto.CreateDriverDocumentTypeRequestDto) (*dto.CreateDriverDocumentTypeResponseDto, error)
		GetByID(ctx context.Context, id uuid.UUID) (*dto.DriverDocumentTypeItemDto, error)
		List(ctx context.Context, page, limit int, search string, serviceID *string) (*dto.ListDriverDocumentTypesResponseDto, error)
		GetRequiredByServiceID(ctx context.Context, serviceID uuid.UUID) (*dto.RequiredDriverDocumentTypesResponseDto, error)
		Update(ctx context.Context, id uuid.UUID, req *dto.UpdateDriverDocumentTypeRequestDto) (*dto.DriverDocumentTypeItemDto, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	driverDocumentTypeUsecase struct {
		repo        appdriverrepo.IDriverDocumentTypeRepository
		serviceRepo websystem_repo.IServiceRepository
		txManager   database.TransactionManager
	}
)

func NewDriverDocumentTypeUsecase(
	repo appdriverrepo.IDriverDocumentTypeRepository,
	serviceRepo websystem_repo.IServiceRepository,
	txManager database.TransactionManager,
) IDriverDocumentTypeUsecase {
	return &driverDocumentTypeUsecase{
		repo:        repo,
		serviceRepo: serviceRepo,
		txManager:   txManager,
	}
}

func serviceIDFromString(s *string) *uuid.UUID {
	if s == nil || *s == "" {
		return nil
	}
	id, err := uuid.Parse(*s)
	if err != nil {
		return nil
	}
	return &id
}

func (u *driverDocumentTypeUsecase) Create(ctx context.Context, req *dto.CreateDriverDocumentTypeRequestDto) (*dto.CreateDriverDocumentTypeResponseDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("Create: start", zap.String(global.KeyCorrelationID, cid), zap.String("code", req.Code))

	if u.serviceRepo != nil {
		if err := validator.CheckExistsOptionalString(ctx, req.ServiceID, func(ctx context.Context, id uuid.UUID) error {
			_, err := u.serviceRepo.GetServiceByID(ctx, id)
			return err
		}, ErrServiceNotFoundForDocType); err != nil {
			global.Logger.Error("Create: failed to validate service", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
			return nil, err
		}
	}

	serviceID := serviceIDFromString(req.ServiceID)
	params := pgdb.CreateDriverDocumentTypeParams{
		Code:              req.Code,
		Name:              req.Name,
		Description:       req.Description,
		IsRequired:        req.IsRequired,
		RequireExpireDate: req.RequireExpireDate,
		ServiceID:         serviceID,
		IsActive:          req.IsActive,
	}

	driver_document_type, err := database.WithTransaction(
		u.txManager,
		ctx,
		func(txCtx context.Context) (*appdrivermodel.DriverDocumentType, error) {
			if serviceID != nil {
				if _, err := u.repo.GetByCodeAndServiceID(txCtx, req.Code, *serviceID); err == nil {
					return nil, ErrDriverDocumentTypeCodeExists
				} else if !errors.Is(err, pgx.ErrNoRows) {
					return nil, err
				}
			} else {
				if _, err := u.repo.GetByCodeGlobal(txCtx, req.Code); err == nil {
					return nil, ErrDriverDocumentTypeCodeExists
				} else if !errors.Is(err, pgx.ErrNoRows) {
					return nil, err
				}
			}
			return u.repo.Create(txCtx, params)
		},
	)
	if err != nil {
		global.Logger.Error("Create: transaction failed", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("Create: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", driver_document_type.ID.String()))
	res := appdrivertransformer.ToCreateDriverDocumentTypeResponseDto(driver_document_type)
	return &res, nil
}

func (u *driverDocumentTypeUsecase) GetByID(ctx context.Context, id uuid.UUID) (*dto.DriverDocumentTypeItemDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("GetByID: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

	m, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			global.Logger.Error("GetByID: document type not found", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
			return nil, ErrDriverDocumentTypeNotFound
		}
		global.Logger.Error("GetByID: failed to get document type", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("GetByID: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	item := appdrivertransformer.ToDriverDocumentTypeItemDto(m)
	return &item, nil
}

func (u *driverDocumentTypeUsecase) List(ctx context.Context, page, limit int, search string, serviceID *string) (*dto.ListDriverDocumentTypesResponseDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("List: start", zap.String(global.KeyCorrelationID, cid), zap.Int("page", page), zap.Int("limit", limit))

	if page < 1 {
		page = constants.DefaultPage
	}
	if limit < 1 || limit > constants.MaxLimit {
		limit = constants.DefaultLimit
	}
	search = strings.TrimSpace(search)
	if len(search) > constants.MaxSearchLen {
		search = search[:constants.MaxSearchLen]
	}
	offset := int32((page - 1) * limit)
	limit32 := int32(limit)

	var total int64
	var items []*appdrivermodel.DriverDocumentType
	if serviceID != nil && *serviceID != "" {
		sid, err := uuid.Parse(*serviceID)
		if err != nil {
			global.Logger.Error("List: failed to parse service_id", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
			return &dto.ListDriverDocumentTypesResponseDto{
				Items:      nil,
				Pagination: dto_common.PaginationMeta{Page: page, Limit: limit, Total: 0},
			}, nil
		}
		total, err = u.repo.CountByServiceID(ctx, sid, search)
		if err != nil {
			global.Logger.Error("List: failed to count by service", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
			return nil, err
		}
		items, err = u.repo.ListByServiceID(ctx, sid, search, limit32, offset)
		if err != nil {
			global.Logger.Error("List: failed to list by service", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
			return nil, err
		}
	} else {
		var err error
		total, err = u.repo.Count(ctx, search)
		if err != nil {
			global.Logger.Error("List: failed to count", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
			return nil, err
		}
		items, err = u.repo.List(ctx, search, limit32, offset)
		if err != nil {
			global.Logger.Error("List: failed to list", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
			return nil, err
		}
	}

	dtos := make([]dto.DriverDocumentTypeItemDto, 0, len(items))
	for _, m := range items {
		dtos = append(dtos, appdrivertransformer.ToDriverDocumentTypeItemDto(m))
	}
	global.Logger.Info("List: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.Int64("total", total))
	return &dto.ListDriverDocumentTypesResponseDto{
		Items: dtos,
		Pagination: dto_common.PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil
}

func (u *driverDocumentTypeUsecase) GetRequiredByServiceID(ctx context.Context, serviceID uuid.UUID) (*dto.RequiredDriverDocumentTypesResponseDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("GetRequiredByServiceID: start", zap.String(global.KeyCorrelationID, cid), zap.String("service_id", serviceID.String()))

	items, err := u.repo.GetRequiredByServiceID(ctx, serviceID)
	if err != nil {
		global.Logger.Error("GetRequiredByServiceID: failed to get required types", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	dtos := make([]dto.DriverDocumentTypeItemDto, 0, len(items))
	for _, m := range items {
		dtos = append(dtos, appdrivertransformer.ToDriverDocumentTypeItemDto(m))
	}
	global.Logger.Info("GetRequiredByServiceID: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.Int("count", len(dtos)))
	return &dto.RequiredDriverDocumentTypesResponseDto{Items: dtos}, nil
}

func (u *driverDocumentTypeUsecase) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateDriverDocumentTypeRequestDto) (*dto.DriverDocumentTypeItemDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("Update: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

	serviceID := serviceIDFromString(req.ServiceID)
	params := pgdb.UpdateDriverDocumentTypeParams{
		ID:                id,
		Code:              req.Code,
		Name:              req.Name,
		Description:       req.Description,
		IsRequired:        req.IsRequired,
		RequireExpireDate: req.RequireExpireDate,
		ServiceID:         serviceID,
		IsActive:          req.IsActive,
	}

	updated, err := database.WithTransaction(
		u.txManager,
		ctx,
		func(txCtx context.Context) (*appdrivermodel.DriverDocumentType, error) {
			existing, err := u.repo.GetByID(txCtx, id)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return nil, ErrDriverDocumentTypeNotFound
				}
				return nil, err
			}
			if existing.Code != req.Code || !serviceIDEqual(existing.ServiceID, serviceID) {
				if serviceID != nil {
					if _, err := u.repo.GetByCodeAndServiceID(txCtx, req.Code, *serviceID); err == nil {
						return nil, ErrDriverDocumentTypeCodeExists
					} else if !errors.Is(err, pgx.ErrNoRows) {
						return nil, err
					}
				} else {
					if _, err := u.repo.GetByCodeGlobal(txCtx, req.Code); err == nil {
						return nil, ErrDriverDocumentTypeCodeExists
					} else if !errors.Is(err, pgx.ErrNoRows) {
						return nil, err
					}
				}
			}
			return u.repo.Update(txCtx, params)
		},
	)
	if err != nil {
		global.Logger.Error("Update: transaction failed", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("Update: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	item := appdrivertransformer.ToDriverDocumentTypeItemDto(updated)
	return &item, nil
}

func serviceIDEqual(a, b *uuid.UUID) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func (u *driverDocumentTypeUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("Delete: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

	_, err := database.WithTransaction(
		u.txManager,
		ctx,
		func(txCtx context.Context) (struct{}, error) {
			_, err := u.repo.GetByID(txCtx, id)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return struct{}{}, ErrDriverDocumentTypeNotFound
				}
				return struct{}{}, err
			}
			return struct{}{}, u.repo.Delete(txCtx, id)
		},
	)
	if err != nil {
		if errors.Is(err, ErrDriverDocumentTypeNotFound) {
			global.Logger.Error("Delete: document type not found", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
		} else {
			global.Logger.Error("Delete: failed to delete document type", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		}
		return err
	}
	global.Logger.Info("Delete: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	return nil
}
