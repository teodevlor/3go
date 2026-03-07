package web_system

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"go-structure/global"
	"go-structure/internal/constants"
	"go-structure/internal/middleware"
	dto_common "go-structure/internal/dto/common"
	dto "go-structure/internal/dto/web_system"
	"go-structure/internal/helper/database"
	"go-structure/internal/repository"
	"go-structure/internal/repository/model"
	zoneTransformer "go-structure/internal/transformer/web_system"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

var (
	ErrZoneNotFound   = errors.New("không tìm thấy zone")
	ErrZoneCodeExists = errors.New("mã zone đã tồn tại")
)

type (
	IZoneUsecase interface {
		CreateZone(ctx context.Context, req *dto.CreateZoneRequestDto) (*dto.CreateZoneResponseDto, error)
		GetZone(ctx context.Context, id uuid.UUID) (*dto.ZoneItemDto, error)
		ListZones(ctx context.Context, page int, limit int, search string) (*dto.ListZonesResponseDto, error)
		UpdateZone(ctx context.Context, id uuid.UUID, req *dto.UpdateZoneRequestDto) (*dto.ZoneItemDto, error)
		DeleteZone(ctx context.Context, id uuid.UUID) error
	}

	zoneUsecase struct {
		zoneRepository     repository.IZoneRepository
		transactionManager database.TransactionManager
	}
)

func NewZoneUsecase(zoneRepository repository.IZoneRepository, transactionManager database.TransactionManager) IZoneUsecase {
	return &zoneUsecase{
		zoneRepository:     zoneRepository,
		transactionManager: transactionManager,
	}
}

func (u *zoneUsecase) CreateZone(ctx context.Context, req *dto.CreateZoneRequestDto) (*dto.CreateZoneResponseDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("CreateZone: start", zap.String(global.KeyCorrelationID, cid), zap.String("code", req.Code))

	if u.zoneRepository == nil {
		return nil, nil
	}
	polygonBytes, err := json.Marshal(req.Polygon)
	if err != nil {
		global.Logger.Error("CreateZone: failed to marshal polygon", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	zoneInput := &model.Zone{
		Code:            req.Code,
		Name:            req.Name,
		Polygon:         string(polygonBytes),
		PriceMultiplier: req.PriceMultiplier,
		IsActive:        req.IsActive,
	}

	zone, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (*model.Zone, error) {
			if _, err := u.zoneRepository.GetZoneByCode(txCtx, req.Code); err == nil {
				return nil, ErrZoneCodeExists
			} else if !errors.Is(err, pgx.ErrNoRows) {
				return nil, err
			}
			return u.zoneRepository.CreateZone(txCtx, zoneInput)
		},
	)
	if err != nil {
		global.Logger.Error("CreateZone: transaction failed", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("CreateZone: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("zone_id", zone.ID.String()))
	return &dto.CreateZoneResponseDto{
		ID:              zone.ID.String(),
		Code:            zone.Code,
		Name:            zone.Name,
		PriceMultiplier: zone.PriceMultiplier,
		Polygon:         string(polygonBytes),
		IsActive:        zone.IsActive,
	}, nil
}

func (u *zoneUsecase) GetZone(ctx context.Context, id uuid.UUID) (*dto.ZoneItemDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("GetZone: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

	if u.zoneRepository == nil {
		global.Logger.Error("GetZone: zone repository nil", zap.String(global.KeyCorrelationID, cid))
		return nil, ErrZoneNotFound
	}
	zone, err := u.zoneRepository.GetZoneByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			global.Logger.Error("GetZone: zone not found", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
			return nil, ErrZoneNotFound
		}
		global.Logger.Error("GetZone: failed to get zone", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("GetZone: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	item := zoneTransformer.ToZoneItemDto(zone)
	return &item, nil
}

func (u *zoneUsecase) ListZones(ctx context.Context, page int, limit int, search string) (*dto.ListZonesResponseDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("ListZones: start", zap.String(global.KeyCorrelationID, cid), zap.Int("page", page), zap.Int("limit", limit))

	if u.zoneRepository == nil {
		return &dto.ListZonesResponseDto{Items: nil, Pagination: dto_common.PaginationMeta{Page: page, Limit: limit, Total: 0}}, nil
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

	total, err := u.zoneRepository.CountZones(ctx, search)
	if err != nil {
		global.Logger.Error("ListZones: failed to count zones", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	zones, err := u.zoneRepository.ListZones(ctx, search, limit32, offset)
	if err != nil {
		global.Logger.Error("ListZones: failed to list zones", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	items := make([]dto.ZoneItemDto, 0, len(zones))
	for _, z := range zones {
		items = append(items, zoneTransformer.ToZoneItemDto(z))
	}
	global.Logger.Info("ListZones: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.Int64("total", total))
	return &dto.ListZonesResponseDto{
		Items: items,
		Pagination: dto_common.PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil
}

func (u *zoneUsecase) UpdateZone(ctx context.Context, id uuid.UUID, req *dto.UpdateZoneRequestDto) (*dto.ZoneItemDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("UpdateZone: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

	if u.zoneRepository == nil {
		global.Logger.Error("UpdateZone: zone repository nil", zap.String(global.KeyCorrelationID, cid))
		return nil, ErrZoneNotFound
	}
	polygonBytes, err := json.Marshal(req.Polygon)
	if err != nil {
		global.Logger.Error("UpdateZone: failed to marshal polygon", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	zoneInput := &model.Zone{
		ID:              id,
		Code:            req.Code,
		Name:            req.Name,
		Polygon:         string(polygonBytes),
		PriceMultiplier: req.PriceMultiplier,
		IsActive:        req.IsActive,
	}

	zone, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (*model.Zone, error) {
			updated, err := u.zoneRepository.UpdateZone(txCtx, zoneInput)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return nil, ErrZoneNotFound
				}
				return nil, err
			}
			return updated, nil
		},
	)
	if err != nil {
		global.Logger.Error("UpdateZone: transaction failed", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("UpdateZone: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	return &dto.ZoneItemDto{
		ID:              zone.ID.String(),
		Code:            zone.Code,
		Name:            zone.Name,
		PriceMultiplier: zone.PriceMultiplier,
		Polygon:         string(polygonBytes),
		IsActive:        zone.IsActive,
	}, nil
}

func (u *zoneUsecase) DeleteZone(ctx context.Context, id uuid.UUID) error {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("DeleteZone: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

	if u.zoneRepository == nil {
		global.Logger.Error("DeleteZone: zone repository nil", zap.String(global.KeyCorrelationID, cid))
		return ErrZoneNotFound
	}
	_, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (struct{}, error) {
			return struct{}{}, u.zoneRepository.DeleteZone(txCtx, id)
		},
	)
	if err != nil {
		global.Logger.Error("DeleteZone: failed to delete zone", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return err
	}
	global.Logger.Info("DeleteZone: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	return nil
}
