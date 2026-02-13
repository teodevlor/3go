package web_system

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"go-structure/global"
	common "go-structure/internal/common"
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
	if u.zoneRepository == nil {
		return nil, nil
	}
	polygonBytes, err := json.Marshal(req.Polygon)
	if err != nil {
		return nil, err
	}
	zoneInput := &model.Zone{
		Code:            req.Code,
		Name:            req.Name,
		Polygon:         string(polygonBytes),
		PriceMultiplier: req.PriceMultiplier,
		IsActive:        req.IsActive,
	}
	var zone *model.Zone
	err = u.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if _, err := u.zoneRepository.GetZoneByCode(txCtx, req.Code); err == nil {
			return ErrZoneCodeExists
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		created, err := u.zoneRepository.CreateZone(txCtx, zoneInput)
		if err != nil {
			return err
		}
		zone = created
		return nil
	})
	if err != nil {
		return nil, err
	}
	global.GetChannelLogger("zone").Info("create zone", zap.Any("zone", zone))
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
	if u.zoneRepository == nil {
		return nil, ErrZoneNotFound
	}
	zone, err := u.zoneRepository.GetZoneByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrZoneNotFound
		}
		return nil, err
	}
	item := zoneTransformer.ToZoneItemDto(zone)
	return &item, nil
}

func (u *zoneUsecase) ListZones(ctx context.Context, page int, limit int, search string) (*dto.ListZonesResponseDto, error) {
	if u.zoneRepository == nil {
		return &dto.ListZonesResponseDto{Items: nil, Pagination: dto_common.PaginationMeta{Page: page, Limit: limit, Total: 0}}, nil
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

	total, err := u.zoneRepository.CountZones(ctx, search)
	if err != nil {
		return nil, err
	}
	zones, err := u.zoneRepository.ListZones(ctx, search, limit32, offset)
	if err != nil {
		return nil, err
	}
	items := make([]dto.ZoneItemDto, 0, len(zones))
	for _, z := range zones {
		items = append(items, zoneTransformer.ToZoneItemDto(z))
	}
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
	if u.zoneRepository == nil {
		return nil, ErrZoneNotFound
	}
	polygonBytes, err := json.Marshal(req.Polygon)
	if err != nil {
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
	var zone *model.Zone
	err = u.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		updated, err := u.zoneRepository.UpdateZone(txCtx, zoneInput)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrZoneNotFound
			}
			return err
		}
		zone = updated
		return nil
	})
	if err != nil {
		return nil, err
	}
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
	if u.zoneRepository == nil {
		return ErrZoneNotFound
	}
	return u.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		return u.zoneRepository.DeleteZone(txCtx, id)
	})
}
