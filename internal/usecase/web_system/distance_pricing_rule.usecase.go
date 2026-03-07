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
		repo               websystem_repo.IDistancePricingRuleRepository
		serviceRepo        websystem_repo.IServiceRepository
		transactionManager database.TransactionManager
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
	params := pgdb.CreateDistancePricingRuleParams{
		ServiceID:  serviceID,
		FromKm:     common.Float64ToNumeric(req.FromKm),
		ToKm:       common.Float64ToNumeric(req.ToKm),
		PricePerKm: common.Float64ToNumeric(req.PricePerKm),
		IsActive:   req.IsActive,
	}

	rule, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (*websystem.DistancePricingRule, error) {
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
	global.Logger.Info("Create: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", rule.ID.String()))
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
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("GetByID: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

	if u.repo == nil {
		global.Logger.Error("GetByID: repository nil", zap.String(global.KeyCorrelationID, cid))
		return nil, ErrDistancePricingRuleNotFound
	}
	rule, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			global.Logger.Error("GetByID: rule not found", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
			return nil, ErrDistancePricingRuleNotFound
		}
		global.Logger.Error("GetByID: failed to get rule", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("GetByID: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	var svc *websystem.Service
	if u.serviceRepo != nil {
		if s, err := u.serviceRepo.GetServiceByID(ctx, rule.ServiceID); err == nil {
			svc = s
		}
	}
	item := serviceTransformer.ToDistancePricingRuleItemDtoWithService(rule, svc)
	return &item, nil
}

func (u *distancePricingRuleUsecase) List(ctx context.Context, serviceID *uuid.UUID) ([]dto.DistancePricingRuleItemDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("List: start", zap.String(global.KeyCorrelationID, cid))

	if u.repo == nil {
		return nil, nil
	}
	rules, err := u.repo.List(ctx, serviceID)
	if err != nil {
		global.Logger.Error("List: failed to list rules", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("List: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.Int("count", len(rules)))
	serviceByID := make(map[uuid.UUID]*websystem.Service)
	if u.serviceRepo != nil {
		seen := make(map[uuid.UUID]struct{})
		for _, r := range rules {
			if _, ok := seen[r.ServiceID]; ok {
				continue
			}
			seen[r.ServiceID] = struct{}{}
			if s, err := u.serviceRepo.GetServiceByID(ctx, r.ServiceID); err == nil {
				serviceByID[r.ServiceID] = s
			}
		}
	}
	items := make([]dto.DistancePricingRuleItemDto, 0, len(rules))
	for _, r := range rules {
		svc := serviceByID[r.ServiceID]
		items = append(items, serviceTransformer.ToDistancePricingRuleItemDtoWithService(r, svc))
	}
	return items, nil
}

func (u *distancePricingRuleUsecase) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateDistancePricingRuleRequestDto) (*dto.DistancePricingRuleItemDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("Update: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

	if u.repo == nil {
		global.Logger.Error("Update: repository nil", zap.String(global.KeyCorrelationID, cid))
		return nil, ErrDistancePricingRuleNotFound
	}
	serviceID, err := uuid.Parse(req.ServiceID)
	if err != nil {
		global.Logger.Error("Update: failed to parse service_id", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
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

	rule, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (*websystem.DistancePricingRule, error) {
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
					return nil, ErrDistancePricingRuleNotFound
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
	var svc *websystem.Service
	if u.serviceRepo != nil {
		if s, err := u.serviceRepo.GetServiceByID(ctx, rule.ServiceID); err == nil {
			svc = s
		}
	}
	item := serviceTransformer.ToDistancePricingRuleItemDtoWithService(rule, svc)
	return &item, nil
}

func (u *distancePricingRuleUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("Delete: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

	if u.repo == nil {
		global.Logger.Error("Delete: repository nil", zap.String(global.KeyCorrelationID, cid))
		return ErrDistancePricingRuleNotFound
	}
	_, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (struct{}, error) {
			return struct{}{}, u.repo.Delete(txCtx, id)
		},
	)
	if err != nil {
		global.Logger.Error("Delete: failed to delete rule", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return err
	}
	global.Logger.Info("Delete: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	return nil
}
