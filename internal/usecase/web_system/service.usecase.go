package web_system

import (
	"context"
	"errors"
	"strings"

	"go-structure/global"
	common "go-structure/internal/common"
	dto_common "go-structure/internal/dto/common"
	dto "go-structure/internal/dto/web_system"
	"go-structure/internal/helper/database"
	pgdb "go-structure/internal/orm/db/postgres"
	websystem "go-structure/internal/repository/model/web_system"
	websystem_repo "go-structure/internal/repository/web_system"
	serviceTransformer "go-structure/internal/transformer/web_system"

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
		serviceRepository   websystem_repo.IServiceRepository
		serviceZoneUsecase  IServiceZoneUsecase
		transactionManager  database.TransactionManager
	}
)

func NewServiceUsecase(serviceRepository websystem_repo.IServiceRepository, serviceZoneUsecase IServiceZoneUsecase, transactionManager database.TransactionManager) IServiceUsecase {
	return &serviceUsecase{
		serviceRepository:   serviceRepository,
		serviceZoneUsecase:  serviceZoneUsecase,
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
	if u.serviceRepository == nil || u.serviceZoneUsecase == nil {
		return nil, nil
	}
	zoneUUIDs, err := parseZoneIDs(req.ZoneIDs)
	if err != nil {
		return nil, err
	}
	params := pgdb.CreateServiceParams{
		Code:      req.Code,
		Name:      req.Name,
		BasePrice: common.Float64ToNumeric(req.BasePrice),
		MinPrice:  common.Float64ToNumeric(req.MinPrice),
		IsActive:  req.IsActive,
	}
	var svc *websystem.Service
	err = u.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if _, err := u.serviceRepository.GetServiceByCode(txCtx, req.Code); err == nil {
			return ErrServiceCodeExists
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		created, err := u.serviceRepository.CreateService(txCtx, params)
		if err != nil {
			return err
		}
		svc = created
		if err := u.serviceZoneUsecase.SetZonesForService(txCtx, created.ID, zoneUUIDs); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	global.GetChannelLogger("service").Info("create service", zap.Any("service", svc))
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
	if u.serviceRepository == nil || u.serviceZoneUsecase == nil {
		return nil, ErrServiceNotFound
	}
	svc, err := u.serviceRepository.GetServiceByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrServiceNotFound
		}
		return nil, err
	}
	zoneIDs, _ := u.serviceZoneUsecase.GetZoneIDsByServiceID(ctx, id)
	zoneIDStrs := make([]string, 0, len(zoneIDs))
	for _, zid := range zoneIDs {
		zoneIDStrs = append(zoneIDStrs, zid.String())
	}
	item := serviceTransformer.ToServiceItemDto(svc, zoneIDStrs)
	return &item, nil
}

func (u *serviceUsecase) ListServices(ctx context.Context, page int, limit int, search string) (*dto.ListServicesResponseDto, error) {
	if u.serviceRepository == nil || u.serviceZoneUsecase == nil {
		return &dto.ListServicesResponseDto{Items: nil, Pagination: dto_common.PaginationMeta{Page: page, Limit: limit, Total: 0}}, nil
	}
	if page < 1 {
		page = common.DefaultPage
	}
	if limit < 1 || limit > common.MaxLimit {
		limit = common.DefaultLimit
	}
	search = strings.TrimSpace(search)
	if len(search) > common.MaxSearchLen {
		search = search[:common.MaxSearchLen]
	}
	offset := int32((page - 1) * limit)
	limit32 := int32(limit)

	total, err := u.serviceRepository.CountServices(ctx, search)
	if err != nil {
		return nil, err
	}
	services, err := u.serviceRepository.ListServices(ctx, search, limit32, offset)
	if err != nil {
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
	if u.serviceRepository == nil || u.serviceZoneUsecase == nil {
		return nil, ErrServiceNotFound
	}
	zoneUUIDs, err := parseZoneIDs(req.ZoneIDs)
	if err != nil {
		return nil, err
	}
	params := pgdb.UpdateServiceParams{
		ID:        id,
		Code:      req.Code,
		Name:      req.Name,
		BasePrice: common.Float64ToNumeric(req.BasePrice),
		MinPrice:  common.Float64ToNumeric(req.MinPrice),
		IsActive:  req.IsActive,
	}
	var svc *websystem.Service
	err = u.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		updated, err := u.serviceRepository.UpdateService(txCtx, params)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrServiceNotFound
			}
			return err
		}
		svc = updated
		if err := u.serviceZoneUsecase.SetZonesForService(txCtx, id, zoneUUIDs); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	zoneIDStrs := make([]string, 0, len(zoneUUIDs))
	for _, zid := range zoneUUIDs {
		zoneIDStrs = append(zoneIDStrs, zid.String())
	}
	item := serviceTransformer.ToServiceItemDto(svc, zoneIDStrs)
	return &item, nil
}

func (u *serviceUsecase) DeleteService(ctx context.Context, id uuid.UUID) error {
	if u.serviceRepository == nil {
		return ErrServiceNotFound
	}
	return u.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		return u.serviceRepository.DeleteService(txCtx, id)
	})
}
