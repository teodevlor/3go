package web_system

import (
	"context"
	"errors"

	common "go-structure/internal/common"
	dto "go-structure/internal/dto/web_system"
	"go-structure/internal/helper/database"
	pgdb "go-structure/internal/orm/db/postgres"
	websystem "go-structure/internal/repository/model/web_system"
	websystem_repo "go-structure/internal/repository/web_system"
	serviceTransformer "go-structure/internal/transformer/web_system"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	if u.repo == nil {
		return nil, nil
	}
	serviceID, err := uuid.Parse(req.ServiceID)
	if err != nil {
		return nil, err
	}
	params := pgdb.CreatePackageSizePricingParams{
		ServiceID:   serviceID,
		PackageSize: req.PackageSize,
		ExtraPrice:  common.Float64ToNumeric(req.ExtraPrice),
		IsActive:    req.IsActive,
	}

	var item *websystem.PackageSizePricing
	err = u.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if u.serviceRepo != nil {
			if _, err := u.serviceRepo.GetServiceByID(txCtx, serviceID); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ErrServiceNotFound
				}
				return err
			}
		}
		created, err := u.repo.Create(txCtx, params)
		if err != nil {
			return err
		}
		item = created
		return nil
	})
	if err != nil {
		return nil, err
	}
	d := serviceTransformer.ToPackageSizePricingItemDto(item)
	return &d, nil
}

func (u *packageSizePricingUsecase) GetByID(ctx context.Context, id uuid.UUID) (*dto.PackageSizePricingItemDto, error) {
	if u.repo == nil {
		return nil, ErrPackageSizePricingNotFound
	}
	item, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPackageSizePricingNotFound
		}
		return nil, err
	}
	d := serviceTransformer.ToPackageSizePricingItemDto(item)
	return &d, nil
}

func (u *packageSizePricingUsecase) List(ctx context.Context, serviceID *uuid.UUID) ([]dto.PackageSizePricingItemDto, error) {
	if u.repo == nil {
		return nil, nil
	}
	items, err := u.repo.List(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.PackageSizePricingItemDto, 0, len(items))
	for _, r := range items {
		out = append(out, serviceTransformer.ToPackageSizePricingItemDto(r))
	}
	return out, nil
}

func (u *packageSizePricingUsecase) Update(ctx context.Context, id uuid.UUID, req *dto.UpdatePackageSizePricingRequestDto) (*dto.PackageSizePricingItemDto, error) {
	if u.repo == nil {
		return nil, ErrPackageSizePricingNotFound
	}
	serviceID, err := uuid.Parse(req.ServiceID)
	if err != nil {
		return nil, err
	}
	params := pgdb.UpdatePackageSizePricingParams{
		ID:          id,
		ServiceID:   serviceID,
		PackageSize: req.PackageSize,
		ExtraPrice:  common.Float64ToNumeric(req.ExtraPrice),
		IsActive:    req.IsActive,
	}

	var item *websystem.PackageSizePricing
	err = u.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if u.serviceRepo != nil {
			if _, err := u.serviceRepo.GetServiceByID(txCtx, serviceID); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ErrServiceNotFound
				}
				return err
			}
		}
		updated, err := u.repo.Update(txCtx, params)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrPackageSizePricingNotFound
			}
			return err
		}
		item = updated
		return nil
	})
	if err != nil {
		return nil, err
	}
	d := serviceTransformer.ToPackageSizePricingItemDto(item)
	return &d, nil
}

func (u *packageSizePricingUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if u.repo == nil {
		return ErrPackageSizePricingNotFound
	}
	return u.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		return u.repo.Delete(txCtx, id)
	})
}
