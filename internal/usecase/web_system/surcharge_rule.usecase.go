package web_system

import (
	"context"
	"errors"

	common "go-structure/internal/common"
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
		Create(ctx context.Context, req *dto.CreateSurchargeRuleRequestDto) (*dto.SurchargeRuleItemDto, error)
		GetByID(ctx context.Context, id uuid.UUID) (*dto.SurchargeRuleItemDto, error)
		List(ctx context.Context, serviceID, zoneID *uuid.UUID) ([]dto.SurchargeRuleItemDto, error)
		Update(ctx context.Context, id uuid.UUID, req *dto.UpdateSurchargeRuleRequestDto) (*dto.SurchargeRuleItemDto, error)
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

func (u *surchargeRuleUsecase) Create(ctx context.Context, req *dto.CreateSurchargeRuleRequestDto) (*dto.SurchargeRuleItemDto, error) {
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
	condition := req.Condition
	if condition == nil {
		condition = []byte("{}")
	}
	params := pgdb.CreateSurchargeRuleParams{
		ServiceID:     serviceID,
		ZoneID:        zoneID,
		SurchargeType: req.SurchargeType,
		Amount:        common.Float64ToNumeric(req.Amount),
		Unit:          req.Unit,
		Condition:     condition,
		IsActive:      req.IsActive,
	}

	var rule *websystem_model.SurchargeRule
	err = u.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if u.serviceRepo != nil {
			if _, err := u.serviceRepo.GetServiceByID(txCtx, serviceID); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ErrServiceNotFound
				}
				return err
			}
		}
		if u.zoneRepo != nil {
			if _, err := u.zoneRepo.GetZoneByID(txCtx, zoneID); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ErrZoneNotFound
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

func (u *surchargeRuleUsecase) List(ctx context.Context, serviceID, zoneID *uuid.UUID) ([]dto.SurchargeRuleItemDto, error) {
	if u.repo == nil {
		return nil, nil
	}
	rules, err := u.repo.List(ctx, serviceID, zoneID)
	if err != nil {
		return nil, err
	}
	items := make([]dto.SurchargeRuleItemDto, 0, len(rules))
	for _, r := range rules {
		items = append(items, serviceTransformer.ToSurchargeRuleItemDto(r))
	}
	return items, nil
}

func (u *surchargeRuleUsecase) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateSurchargeRuleRequestDto) (*dto.SurchargeRuleItemDto, error) {
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
	condition := req.Condition
	if condition == nil {
		condition = []byte("{}")
	}
	params := pgdb.UpdateSurchargeRuleParams{
		ID:            id,
		ServiceID:     serviceID,
		ZoneID:        zoneID,
		SurchargeType: req.SurchargeType,
		Amount:        common.Float64ToNumeric(req.Amount),
		Unit:          req.Unit,
		Condition:     condition,
		IsActive:      req.IsActive,
	}

	var rule *websystem_model.SurchargeRule
	err = u.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if u.serviceRepo != nil {
			if _, err := u.serviceRepo.GetServiceByID(txCtx, serviceID); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ErrServiceNotFound
				}
				return err
			}
		}
		if u.zoneRepo != nil {
			if _, err := u.zoneRepo.GetZoneByID(txCtx, zoneID); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return ErrZoneNotFound
				}
				return err
			}
		}
		updated, err := u.repo.Update(txCtx, params)
		if err != nil {
			return err
		}
		rule = updated
		return nil
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSurchargeRuleNotFound
		}
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
