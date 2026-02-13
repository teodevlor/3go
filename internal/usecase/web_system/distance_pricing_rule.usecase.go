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

var ErrDistancePricingRuleNotFound = errors.New("không tìm thấy quy tắc giá theo khoảng cách")

type (
	IDistancePricingRuleUsecase interface {
		Create(ctx context.Context, req *dto.CreateDistancePricingRuleRequestDto) (*dto.CreateDistancePricingRuleResponseDto, error)
		GetByID(ctx context.Context, id uuid.UUID) (*dto.DistancePricingRuleItemDto, error)
		List(ctx context.Context, serviceID *uuid.UUID) ([]dto.DistancePricingRuleItemDto, error)
		Update(ctx context.Context, id uuid.UUID, req *dto.UpdateDistancePricingRuleRequestDto) (*dto.DistancePricingRuleItemDto, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	distancePricingRuleUsecase struct {
		repo                  websystem_repo.IDistancePricingRuleRepository
		serviceRepo           websystem_repo.IServiceRepository
		transactionManager    database.TransactionManager
	}
)

func NewDistancePricingRuleUsecase(repo websystem_repo.IDistancePricingRuleRepository, serviceRepo websystem_repo.IServiceRepository, transactionManager database.TransactionManager) IDistancePricingRuleUsecase {
	return &distancePricingRuleUsecase{
		repo:               repo,
		serviceRepo:        serviceRepo,
		transactionManager: transactionManager,
	}
}

func (u *distancePricingRuleUsecase) Create(ctx context.Context, req *dto.CreateDistancePricingRuleRequestDto) (*dto.CreateDistancePricingRuleResponseDto, error) {
	if u.repo == nil {
		return nil, nil
	}
	serviceID, err := uuid.Parse(req.ServiceID)
	if err != nil {
		return nil, err
	}
	params := pgdb.CreateDistancePricingRuleParams{
		ServiceID:  serviceID,
		FromKm:     common.Float64ToNumeric(req.FromKm),
		ToKm:       common.Float64ToNumeric(req.ToKm),
		PricePerKm: common.Float64ToNumeric(req.PricePerKm),
		IsActive:   req.IsActive,
	}

	var rule *websystem.DistancePricingRule
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
		rule = created
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &dto.CreateDistancePricingRuleResponseDto{
		ID:         rule.ID.String(),
		ServiceID:  rule.ServiceID.String(),
		FromKm:     rule.FromKm,
		ToKm:       rule.ToKm,
		PricePerKm: rule.PricePerKm,
		IsActive:   rule.IsActive,
	}, nil
}

func (u *distancePricingRuleUsecase) GetByID(ctx context.Context, id uuid.UUID) (*dto.DistancePricingRuleItemDto, error) {
	if u.repo == nil {
		return nil, ErrDistancePricingRuleNotFound
	}
	rule, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDistancePricingRuleNotFound
		}
		return nil, err
	}
	item := serviceTransformer.ToDistancePricingRuleItemDto(rule)
	return &item, nil
}

func (u *distancePricingRuleUsecase) List(ctx context.Context, serviceID *uuid.UUID) ([]dto.DistancePricingRuleItemDto, error) {
	if u.repo == nil {
		return nil, nil
	}
	rules, err := u.repo.List(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	items := make([]dto.DistancePricingRuleItemDto, 0, len(rules))
	for _, r := range rules {
		items = append(items, serviceTransformer.ToDistancePricingRuleItemDto(r))
	}
	return items, nil
}

func (u *distancePricingRuleUsecase) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateDistancePricingRuleRequestDto) (*dto.DistancePricingRuleItemDto, error) {
	if u.repo == nil {
		return nil, ErrDistancePricingRuleNotFound
	}
	serviceID, err := uuid.Parse(req.ServiceID)
	if err != nil {
		return nil, err
	}
	params := pgdb.UpdateDistancePricingRuleParams{
		ID:         id,
		ServiceID:  serviceID,
		FromKm:     common.Float64ToNumeric(req.FromKm),
		ToKm:       common.Float64ToNumeric(req.ToKm),
		PricePerKm: common.Float64ToNumeric(req.PricePerKm),
		IsActive:   req.IsActive,
	}
	var rule *websystem.DistancePricingRule
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
				return ErrDistancePricingRuleNotFound
			}
			return err
		}
		rule = updated
		return nil
	})
	if err != nil {
		return nil, err
	}
	item := serviceTransformer.ToDistancePricingRuleItemDto(rule)
	return &item, nil
}

func (u *distancePricingRuleUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if u.repo == nil {
		return ErrDistancePricingRuleNotFound
	}
	return u.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		return u.repo.Delete(txCtx, id)
	})
}
