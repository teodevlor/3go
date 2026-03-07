package web_system

import (
	"context"
	"errors"

	"go-structure/global"
	common "go-structure/internal/common"
	dto "go-structure/internal/dto/web_system"
	"go-structure/internal/helper/database"
	"go-structure/internal/middleware"
	pgdb "go-structure/orm/db/postgres"
	websystem "go-structure/internal/repository/model/web_system"
	websystem_repo "go-structure/internal/repository/web_system"
	serviceTransformer "go-structure/internal/transformer/web_system"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

var ErrPackageSizePricingNotFound = errors.New("không tìm thấy quy tắc giá theo kích thước gói")

type (
	IPackageSizePricingUsecase interface {
		Create(ctx context.Context, req *dto.CreatePackageSizePricingRequestDto) (*dto.PackageSizePricingItemDto, error)
		GetByID(ctx context.Context, id uuid.UUID) (*dto.PackageSizePricingItemDto, error)
		List(ctx context.Context, serviceID *uuid.UUID) ([]dto.PackageSizePricingItemDto, error)
		Update(ctx context.Context, id uuid.UUID, req *dto.UpdatePackageSizePricingRequestDto) (*dto.PackageSizePricingItemDto, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	packageSizePricingUsecase struct {
		repo               websystem_repo.IPackageSizePricingRepository
		serviceRepo        websystem_repo.IServiceRepository
		transactionManager database.TransactionManager
	}
)

func NewPackageSizePricingUsecase(
	repo websystem_repo.IPackageSizePricingRepository,
	serviceRepo websystem_repo.IServiceRepository,
	transactionManager database.TransactionManager,
) IPackageSizePricingUsecase {
	return &packageSizePricingUsecase{
		repo:               repo,
		serviceRepo:        serviceRepo,
		transactionManager: transactionManager,
	}
}

func (u *packageSizePricingUsecase) Create(ctx context.Context, req *dto.CreatePackageSizePricingRequestDto) (*dto.PackageSizePricingItemDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("Create: start", zap.String(global.KeyCorrelationID, cid), zap.String("service_id", req.ServiceID))

	if u.repo == nil {
		return nil, nil
	}
	serviceID, err := uuid.Parse(req.ServiceID)
	if err != nil {
		global.Logger.Error("Create: failed to parse service_id", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	params := pgdb.CreatePackageSizePricingParams{
		ServiceID:   serviceID,
		PackageSize: req.PackageSize,
		ExtraPrice:  common.Float64ToNumeric(req.ExtraPrice),
		IsActive:    req.IsActive,
	}

	item, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (*websystem.PackageSizePricing, error) {
			if u.serviceRepo != nil {
				if _, err := u.serviceRepo.GetServiceByID(txCtx, serviceID); err != nil {
					if errors.Is(err, pgx.ErrNoRows) {
						return nil, ErrServiceNotFound
					}
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
	global.Logger.Info("Create: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", item.ID.String()))
	d := serviceTransformer.ToPackageSizePricingItemDto(item)
	return &d, nil
}

func (u *packageSizePricingUsecase) GetByID(ctx context.Context, id uuid.UUID) (*dto.PackageSizePricingItemDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("GetByID: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

	if u.repo == nil {
		global.Logger.Error("GetByID: repository nil", zap.String(global.KeyCorrelationID, cid))
		return nil, ErrPackageSizePricingNotFound
	}
	item, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			global.Logger.Error("GetByID: item not found", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
			return nil, ErrPackageSizePricingNotFound
		}
		global.Logger.Error("GetByID: failed to get item", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("GetByID: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	d := serviceTransformer.ToPackageSizePricingItemDto(item)
	return &d, nil
}

func (u *packageSizePricingUsecase) List(ctx context.Context, serviceID *uuid.UUID) ([]dto.PackageSizePricingItemDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("List: start", zap.String(global.KeyCorrelationID, cid))

	if u.repo == nil {
		return nil, nil
	}
	items, err := u.repo.List(ctx, serviceID)
	if err != nil {
		global.Logger.Error("List: failed to list", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	out := make([]dto.PackageSizePricingItemDto, 0, len(items))
	for _, r := range items {
		out = append(out, serviceTransformer.ToPackageSizePricingItemDto(r))
	}
	global.Logger.Info("List: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.Int("count", len(out)))
	return out, nil
}

func (u *packageSizePricingUsecase) Update(ctx context.Context, id uuid.UUID, req *dto.UpdatePackageSizePricingRequestDto) (*dto.PackageSizePricingItemDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("Update: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

	if u.repo == nil {
		global.Logger.Error("Update: repository nil", zap.String(global.KeyCorrelationID, cid))
		return nil, ErrPackageSizePricingNotFound
	}
	serviceID, err := uuid.Parse(req.ServiceID)
	if err != nil {
		global.Logger.Error("Update: failed to parse service_id", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	params := pgdb.UpdatePackageSizePricingParams{
		ID:          id,
		ServiceID:   serviceID,
		PackageSize: req.PackageSize,
		ExtraPrice:  common.Float64ToNumeric(req.ExtraPrice),
		IsActive:    req.IsActive,
	}

	item, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (*websystem.PackageSizePricing, error) {
			if u.serviceRepo != nil {
				if _, err := u.serviceRepo.GetServiceByID(txCtx, serviceID); err != nil {
					if errors.Is(err, pgx.ErrNoRows) {
						return nil, ErrServiceNotFound
					}
					return nil, err
				}
			}
			updated, err := u.repo.Update(txCtx, params)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return nil, ErrPackageSizePricingNotFound
				}
				return nil, err
			}
			return updated, nil
		},
	)
	if err != nil {
		global.Logger.Error("Update: transaction failed", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("Update: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	d := serviceTransformer.ToPackageSizePricingItemDto(item)
	return &d, nil
}

func (u *packageSizePricingUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("Delete: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

	if u.repo == nil {
		global.Logger.Error("Delete: repository nil", zap.String(global.KeyCorrelationID, cid))
		return ErrPackageSizePricingNotFound
	}
	_, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (struct{}, error) {
			return struct{}{}, u.repo.Delete(txCtx, id)
		},
	)
	if err != nil {
		global.Logger.Error("Delete: failed to delete", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return err
	}
	global.Logger.Info("Delete: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	return nil
}
