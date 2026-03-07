package web_system

import (
	"context"
	"errors"
	"strings"

	"go-structure/global"
	common "go-structure/internal/common"
	"go-structure/internal/constants"
	"go-structure/internal/middleware"
	dto_common "go-structure/internal/dto/common"
	dto "go-structure/internal/dto/web_system"
	"go-structure/internal/helper/database"
	pgdb "go-structure/orm/db/postgres"
	account_repo "go-structure/internal/repository"
	websystem "go-structure/internal/repository/model/web_system"
	websystem_repo "go-structure/internal/repository/web_system"
	serviceTransformer "go-structure/internal/transformer/web_system"
	"go-structure/pkg/validator"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

var (
	ErrServiceNotFound   = errors.New("không tìm thấy dịch vụ")
	ErrServiceCodeExists = errors.New("mã dịch vụ đã tồn tại")
)

type (
	IServiceUsecase interface {
		CreateService(ctx context.Context, req *dto.CreateServiceRequestDto) (*dto.CreateServiceResponseDto, error)
		GetService(ctx context.Context, id uuid.UUID) (*dto.ServiceItemDto, error)
		ListServices(ctx context.Context, page int, limit int, search string) (*dto.ListServicesResponseDto, error)
		UpdateService(ctx context.Context, id uuid.UUID, req *dto.UpdateServiceRequestDto) (*dto.ServiceItemDto, error)
		DeleteService(ctx context.Context, id uuid.UUID) error
	}

	serviceUsecase struct {
		serviceRepository  websystem_repo.IServiceRepository
		serviceZoneUsecase IServiceZoneUsecase
		zoneRepo           account_repo.IZoneRepository
		transactionManager database.TransactionManager
	}
)

func NewServiceUsecase(
	serviceRepository websystem_repo.IServiceRepository,
	serviceZoneUsecase IServiceZoneUsecase,
	zoneRepo account_repo.IZoneRepository,
	transactionManager database.TransactionManager,
) IServiceUsecase {
	return &serviceUsecase{
		serviceRepository:  serviceRepository,
		serviceZoneUsecase: serviceZoneUsecase,
		zoneRepo:           zoneRepo,
		transactionManager: transactionManager,
	}
}

func parseZoneIDs(zoneIDs []string) ([]uuid.UUID, error) {
	if len(zoneIDs) == 0 {
		return nil, nil
	}
	ids := make([]uuid.UUID, 0, len(zoneIDs))
	for _, s := range zoneIDs {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		id, err := uuid.Parse(s)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (u *serviceUsecase) CreateService(ctx context.Context, req *dto.CreateServiceRequestDto) (*dto.CreateServiceResponseDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("CreateService: start", zap.String(global.KeyCorrelationID, cid), zap.String("code", req.Code))

	if u.serviceRepository == nil || u.serviceZoneUsecase == nil {
		return nil, nil
	}
	zoneUUIDs, err := parseZoneIDs(req.ZoneIDs)
	if err != nil {
		global.Logger.Error("CreateService: failed to parse zone IDs", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	if u.zoneRepo != nil {
		if err := validator.CheckExistsMany(ctx, zoneUUIDs, func(ctx context.Context, id uuid.UUID) error {
			_, err := u.zoneRepo.GetZoneByID(ctx, id)
			return err
		}, ErrZoneNotFound); err != nil {
			global.Logger.Error("CreateService: failed to validate zones", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
			return nil, err
		}
	}
	params := pgdb.CreateServiceParams{
		Code:      req.Code,
		Name:      req.Name,
		BasePrice: common.Float64ToNumeric(req.BasePrice),
		MinPrice:  common.Float64ToNumeric(req.MinPrice),
		IsActive:  req.IsActive,
	}

	svc, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (*websystem.Service, error) {
			if _, err := u.serviceRepository.GetServiceByCode(txCtx, req.Code); err == nil {
				return nil, ErrServiceCodeExists
			} else if !errors.Is(err, pgx.ErrNoRows) {
				return nil, err
			}
			created, err := u.serviceRepository.CreateService(txCtx, params)
			if err != nil {
				return nil, err
			}
			if err := u.serviceZoneUsecase.SetZonesForService(txCtx, created.ID, zoneUUIDs); err != nil {
				return nil, err
			}
			return created, nil
		},
	)
	if err != nil {
		global.Logger.Error("CreateService: transaction failed", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("CreateService: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("service_id", svc.ID.String()))
	zoneIDStrs := make([]string, 0, len(zoneUUIDs))
	for _, id := range zoneUUIDs {
		zoneIDStrs = append(zoneIDStrs, id.String())
	}
	return &dto.CreateServiceResponseDto{
		ID:        svc.ID.String(),
		Code:      svc.Code,
		Name:      svc.Name,
		BasePrice: svc.BasePrice,
		MinPrice:  svc.MinPrice,
		IsActive:  svc.IsActive,
		ZoneIDs:   zoneIDStrs,
	}, nil
}

func (u *serviceUsecase) GetService(ctx context.Context, id uuid.UUID) (*dto.ServiceItemDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("GetService: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

	if u.serviceRepository == nil || u.serviceZoneUsecase == nil {
		global.Logger.Error("GetService: repository nil", zap.String(global.KeyCorrelationID, cid))
		return nil, ErrServiceNotFound
	}
	svc, err := u.serviceRepository.GetServiceByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			global.Logger.Error("GetService: service not found", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
			return nil, ErrServiceNotFound
		}
		global.Logger.Error("GetService: failed to get service", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	zoneIDs, _ := u.serviceZoneUsecase.GetZoneIDsByServiceID(ctx, id)
	zoneIDStrs := make([]string, 0, len(zoneIDs))
	for _, zid := range zoneIDs {
		zoneIDStrs = append(zoneIDStrs, zid.String())
	}
	global.Logger.Info("GetService: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	item := serviceTransformer.ToServiceItemDto(svc, zoneIDStrs)
	return &item, nil
}

func (u *serviceUsecase) ListServices(ctx context.Context, page int, limit int, search string) (*dto.ListServicesResponseDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("ListServices: start", zap.String(global.KeyCorrelationID, cid), zap.Int("page", page), zap.Int("limit", limit))

	if u.serviceRepository == nil || u.serviceZoneUsecase == nil {
		return &dto.ListServicesResponseDto{Items: nil, Pagination: dto_common.PaginationMeta{Page: page, Limit: limit, Total: 0}}, nil
	}
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

	total, err := u.serviceRepository.CountServices(ctx, search)
	if err != nil {
		global.Logger.Error("ListServices: failed to count", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	services, err := u.serviceRepository.ListServices(ctx, search, limit32, offset)
	if err != nil {
		global.Logger.Error("ListServices: failed to list", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	items := make([]dto.ServiceItemDto, 0, len(services))
	for _, s := range services {
		zoneIDs, _ := u.serviceZoneUsecase.GetZoneIDsByServiceID(ctx, s.ID)
		zoneIDStrs := make([]string, 0, len(zoneIDs))
		for _, zid := range zoneIDs {
			zoneIDStrs = append(zoneIDStrs, zid.String())
		}
		items = append(items, serviceTransformer.ToServiceItemDto(s, zoneIDStrs))
	}
	global.Logger.Info("ListServices: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.Int64("total", total))
	return &dto.ListServicesResponseDto{
		Items: items,
		Pagination: dto_common.PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil
}

func (u *serviceUsecase) UpdateService(ctx context.Context, id uuid.UUID, req *dto.UpdateServiceRequestDto) (*dto.ServiceItemDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("UpdateService: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

	if u.serviceRepository == nil || u.serviceZoneUsecase == nil {
		return nil, ErrServiceNotFound
	}
	zoneUUIDs, err := parseZoneIDs(req.ZoneIDs)
	if err != nil {
		global.Logger.Error("UpdateService: failed to parse zone IDs", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	if u.zoneRepo != nil {
		if err := validator.CheckExistsMany(ctx, zoneUUIDs, func(ctx context.Context, id uuid.UUID) error {
			_, err := u.zoneRepo.GetZoneByID(ctx, id)
			return err
		}, ErrZoneNotFound); err != nil {
			global.Logger.Error("UpdateService: failed to validate zones", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
			return nil, err
		}
	}
	params := pgdb.UpdateServiceParams{
		ID:        id,
		Code:      req.Code,
		Name:      req.Name,
		BasePrice: common.Float64ToNumeric(req.BasePrice),
		MinPrice:  common.Float64ToNumeric(req.MinPrice),
		IsActive:  req.IsActive,
	}

	svc, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (*websystem.Service, error) {
			updated, err := u.serviceRepository.UpdateService(txCtx, params)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return nil, ErrServiceNotFound
				}
				return nil, err
			}
			if err := u.serviceZoneUsecase.SetZonesForService(txCtx, id, zoneUUIDs); err != nil {
				return nil, err
			}
			return updated, nil
		},
	)
	if err != nil {
		global.Logger.Error("UpdateService: transaction failed", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	zoneIDStrs := make([]string, 0, len(zoneUUIDs))
	for _, zid := range zoneUUIDs {
		zoneIDStrs = append(zoneIDStrs, zid.String())
	}
	global.Logger.Info("UpdateService: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	item := serviceTransformer.ToServiceItemDto(svc, zoneIDStrs)
	return &item, nil
}

func (u *serviceUsecase) DeleteService(ctx context.Context, id uuid.UUID) error {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("DeleteService: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

	if u.serviceRepository == nil {
		global.Logger.Error("DeleteService: repository nil", zap.String(global.KeyCorrelationID, cid))
		return ErrServiceNotFound
	}
	_, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (struct{}, error) {
			return struct{}{}, u.serviceRepository.DeleteService(txCtx, id)
		},
	)
	if err != nil {
		global.Logger.Error("DeleteService: failed to delete service", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return err
	}
	global.Logger.Info("DeleteService: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	return nil
}
