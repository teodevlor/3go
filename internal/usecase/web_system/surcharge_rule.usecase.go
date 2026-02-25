package web_system

import (
	"context"
	"errors"

	common "go-structure/internal/common"
	dto_common "go-structure/internal/dto/common"
	dto "go-structure/internal/dto/web_system"
	"go-structure/internal/helper/database"
	pgdb "go-structure/internal/orm/db/postgres"
	account_repo "go-structure/internal/repository"
	websystem_model "go-structure/internal/repository/model/web_system"
	websystem_repo "go-structure/internal/repository/web_system"
	serviceTransformer "go-structure/internal/transformer/web_system"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrSurchargeRuleNotFound = errors.New("không tìm thấy quy tắc phụ thu")

type (
	ISurchargeRuleUsecase interface {
		Create(ctx context.Context, adminID uuid.UUID, req *dto.CreateSurchargeRuleRequestDto) (*dto.SurchargeRuleItemDto, error)
		GetByID(ctx context.Context, id uuid.UUID) (*dto.SurchargeRuleItemDto, error)
		List(ctx context.Context, serviceID, zoneID *uuid.UUID) (*dto.ListSurchargeRulesResponseDto, error)
		Update(ctx context.Context, adminID uuid.UUID, id uuid.UUID, req *dto.UpdateSurchargeRuleRequestDto) (*dto.SurchargeRuleItemDto, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	surchargeRuleUsecase struct {
		repo               websystem_repo.ISurchargeRuleRepository
		serviceRepo        websystem_repo.IServiceRepository
		zoneRepo           account_repo.IZoneRepository
		transactionManager database.TransactionManager
	}
)

func NewSurchargeRuleUsecase(
	repo websystem_repo.ISurchargeRuleRepository,
	serviceRepo websystem_repo.IServiceRepository,
	zoneRepo account_repo.IZoneRepository,
	transactionManager database.TransactionManager,
) ISurchargeRuleUsecase {
	return &surchargeRuleUsecase{
		repo:               repo,
		serviceRepo:        serviceRepo,
		zoneRepo:           zoneRepo,
		transactionManager: transactionManager,
	}
}

func (u *surchargeRuleUsecase) Create(ctx context.Context, adminID uuid.UUID, req *dto.CreateSurchargeRuleRequestDto) (*dto.SurchargeRuleItemDto, error) {
	if u.repo == nil {
		return nil, nil
	}
	serviceID, err := uuid.Parse(req.ServiceID)
	if err != nil {
		return nil, err
	}
	zoneID, err := uuid.Parse(req.ZoneID)
	if err != nil {
		return nil, err
	}
	params := pgdb.CreateSurchargeRuleParams{
		ServiceID: serviceID,
		ZoneID:    zoneID,
		Amount:    common.Float64ToNumeric(req.Amount),
		Unit:      req.Unit,
		IsActive:  req.IsActive,
		Priority:  int32(req.Priority),
		CreatedBy: adminID,
		UpdatedBy: adminID,
	}

	rule, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (*websystem_model.SurchargeRule, error) {
			if u.serviceRepo != nil {
				if _, err := u.serviceRepo.GetServiceByID(txCtx, serviceID); err != nil {
					if errors.Is(err, pgx.ErrNoRows) {
						return nil, ErrServiceNotFound
					}
					return nil, err
				}
			}
			if u.zoneRepo != nil {
				if _, err := u.zoneRepo.GetZoneByID(txCtx, zoneID); err != nil {
					if errors.Is(err, pgx.ErrNoRows) {
						return nil, ErrZoneNotFound
					}
					return nil, err
				}
			}
			return u.repo.Create(txCtx, params)
		},
	)
	if err != nil {
		return nil, err
	}
	item := serviceTransformer.ToSurchargeRuleItemDto(rule)
	return &item, nil
}

func (u *surchargeRuleUsecase) GetByID(ctx context.Context, id uuid.UUID) (*dto.SurchargeRuleItemDto, error) {
	if u.repo == nil {
		return nil, ErrSurchargeRuleNotFound
	}
	rule, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSurchargeRuleNotFound
		}
		return nil, err
	}
	item := serviceTransformer.ToSurchargeRuleItemDto(rule)
	return &item, nil
}

func (u *surchargeRuleUsecase) List(ctx context.Context, serviceID, zoneID *uuid.UUID) (*dto.ListSurchargeRulesResponseDto, error) {
	if u.repo == nil {
		return &dto.ListSurchargeRulesResponseDto{
			Items:      nil,
			Pagination: dto_common.PaginationMeta{Page: 1, Limit: 0, Total: 0},
		}, nil
	}
	rules, err := u.repo.List(ctx, serviceID, zoneID)
	if err != nil {
		return nil, err
	}
	items := make([]dto.SurchargeRuleItemDto, 0, len(rules))
	for _, r := range rules {
		items = append(items, serviceTransformer.ToSurchargeRuleItemDto(r))
	}
	return &dto.ListSurchargeRulesResponseDto{
		Items: items,
		Pagination: dto_common.PaginationMeta{
			Page:  1,
			Limit: len(items),
			Total: int64(len(items)),
		},
	}, nil
}

func (u *surchargeRuleUsecase) Update(ctx context.Context, adminID uuid.UUID, id uuid.UUID, req *dto.UpdateSurchargeRuleRequestDto) (*dto.SurchargeRuleItemDto, error) {
	if u.repo == nil {
		return nil, ErrSurchargeRuleNotFound
	}
	serviceID, err := uuid.Parse(req.ServiceID)
	if err != nil {
		return nil, err
	}
	zoneID, err := uuid.Parse(req.ZoneID)
	if err != nil {
		return nil, err
	}
	params := pgdb.UpdateSurchargeRuleParams{
		ID:        id,
		ServiceID: serviceID,
		ZoneID:    zoneID,
		Amount:    common.Float64ToNumeric(req.Amount),
		Unit:      req.Unit,
		IsActive:  req.IsActive,
		Priority:  int32(req.Priority),
		UpdatedBy: adminID,
	}

	rule, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (*websystem_model.SurchargeRule, error) {
			if u.serviceRepo != nil {
				if _, err := u.serviceRepo.GetServiceByID(txCtx, serviceID); err != nil {
					if errors.Is(err, pgx.ErrNoRows) {
						return nil, ErrServiceNotFound
					}
					return nil, err
				}
			}
			if u.zoneRepo != nil {
				if _, err := u.zoneRepo.GetZoneByID(txCtx, zoneID); err != nil {
					if errors.Is(err, pgx.ErrNoRows) {
						return nil, ErrZoneNotFound
					}
					return nil, err
				}
			}
			updated, err := u.repo.Update(txCtx, params)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return nil, ErrSurchargeRuleNotFound
				}
				return nil, err
			}
			return updated, nil
		},
	)
	if err != nil {
		return nil, err
	}
	item := serviceTransformer.ToSurchargeRuleItemDto(rule)
	return &item, nil
}

func (u *surchargeRuleUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if u.repo == nil {
		return ErrSurchargeRuleNotFound
	}
	return u.repo.Delete(ctx, id)
}
