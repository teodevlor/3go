package web_system

import (
	"context"
	"errors"

	"go-structure/global"
	dto_common "go-structure/internal/dto/common"
	dto "go-structure/internal/dto/web_system"
	"go-structure/internal/helper/database"
	"go-structure/internal/middleware"
	pgdb "go-structure/orm/db/postgres"
	websystem_model "go-structure/internal/repository/model/web_system"
	websystem_repo "go-structure/internal/repository/web_system"
	serviceTransformer "go-structure/internal/transformer/web_system"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

var (
	ErrSurchargeConditionNotFound   = errors.New("không tìm thấy điều kiện phụ thu")
	ErrSurchargeConditionCodeExists = errors.New("mã điều kiện phụ thu đã tồn tại")
)

type (
	ISurchargeConditionUsecase interface {
		Create(ctx context.Context, req *dto.CreateSurchargeConditionRequestDto) (*dto.SurchargeConditionItemDto, error)
		GetByID(ctx context.Context, id uuid.UUID) (*dto.SurchargeConditionItemDto, error)
		List(ctx context.Context) (*dto.ListSurchargeConditionsResponseDto, error)
		Update(ctx context.Context, id uuid.UUID, req *dto.UpdateSurchargeConditionRequestDto) (*dto.SurchargeConditionItemDto, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	surchargeConditionUsecase struct {
		repo               websystem_repo.ISurchargeConditionRepository
		transactionManager database.TransactionManager
	}
)

func NewSurchargeConditionUsecase(
	repo websystem_repo.ISurchargeConditionRepository,
	transactionManager database.TransactionManager,
) ISurchargeConditionUsecase {
	return &surchargeConditionUsecase{
		repo:               repo,
		transactionManager: transactionManager,
	}
}

func (u *surchargeConditionUsecase) Create(ctx context.Context, req *dto.CreateSurchargeConditionRequestDto) (*dto.SurchargeConditionItemDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("Create: start", zap.String(global.KeyCorrelationID, cid), zap.String("code", req.Code))

	if u.repo == nil {
		return nil, nil
	}

	if err := websystem_model.ValidateConditionConfig(req.ConditionType, req.Config); err != nil {
		global.Logger.Error("Create: failed to validate condition config", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}

	params := pgdb.CreateSurchargeConditionParams{
		Code:          req.Code,
		Name:          req.Name,
		ConditionType: req.ConditionType,
		Config:        req.Config,
		IsActive:      pgtype.Bool{Bool: req.IsActive, Valid: true},
	}

	condition, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (*websystem_model.SurchargeCondition, error) {
			if _, err := u.repo.GetByCode(txCtx, req.Code); err == nil {
				return nil, ErrSurchargeConditionCodeExists
			} else if !errors.Is(err, pgx.ErrNoRows) {
				return nil, err
			}
			return u.repo.Create(txCtx, params)
		},
	)
	if err != nil {
		global.Logger.Error("Create: transaction failed", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("Create: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", condition.ID.String()))
	item := serviceTransformer.ToSurchargeConditionItemDto(condition)
	return &item, nil
}

func (u *surchargeConditionUsecase) GetByID(ctx context.Context, id uuid.UUID) (*dto.SurchargeConditionItemDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("GetByID: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

	if u.repo == nil {
		global.Logger.Error("GetByID: repository nil", zap.String(global.KeyCorrelationID, cid))
		return nil, ErrSurchargeConditionNotFound
	}
	cond, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			global.Logger.Error("GetByID: condition not found", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
			return nil, ErrSurchargeConditionNotFound
		}
		global.Logger.Error("GetByID: failed to get condition", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("GetByID: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	item := serviceTransformer.ToSurchargeConditionItemDto(cond)
	return &item, nil
}

func (u *surchargeConditionUsecase) List(ctx context.Context) (*dto.ListSurchargeConditionsResponseDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("List: start", zap.String(global.KeyCorrelationID, cid))

	if u.repo == nil {
		return &dto.ListSurchargeConditionsResponseDto{
			Items:      nil,
			Pagination: dto_common.PaginationMeta{Page: 1, Limit: 0, Total: 0},
		}, nil
	}
	conds, err := u.repo.List(ctx)
	if err != nil {
		global.Logger.Error("List: failed to list conditions", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("List: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.Int("count", len(conds)))
	items := make([]dto.SurchargeConditionItemDto, 0, len(conds))
	for _, c := range conds {
		items = append(items, serviceTransformer.ToSurchargeConditionItemDto(c))
	}
	return &dto.ListSurchargeConditionsResponseDto{
		Items: items,
		Pagination: dto_common.PaginationMeta{
			Page:  1,
			Limit: len(items),
			Total: int64(len(items)),
		},
	}, nil
}

func (u *surchargeConditionUsecase) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateSurchargeConditionRequestDto) (*dto.SurchargeConditionItemDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("Update: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

	if u.repo == nil {
		global.Logger.Error("Update: repository nil", zap.String(global.KeyCorrelationID, cid))
		return nil, ErrSurchargeConditionNotFound
	}

	if err := websystem_model.ValidateConditionConfig(req.ConditionType, req.Config); err != nil {
		global.Logger.Error("Update: failed to validate condition config", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}

	params := pgdb.UpdateSurchargeConditionParams{
		ID:            id,
		Code:          req.Code,
		Name:          req.Name,
		ConditionType: req.ConditionType,
		Config:        req.Config,
		IsActive:      pgtype.Bool{Bool: req.IsActive, Valid: true},
	}

	cond, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (*websystem_model.SurchargeCondition, error) {
			existing, err := u.repo.GetByCode(txCtx, req.Code)
			if err == nil && existing.ID != id {
				return nil, ErrSurchargeConditionCodeExists
			}
			if err != nil && !errors.Is(err, pgx.ErrNoRows) {
				return nil, err
			}

			updated, err := u.repo.Update(txCtx, params)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return nil, ErrSurchargeConditionNotFound
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
	item := serviceTransformer.ToSurchargeConditionItemDto(cond)
	return &item, nil
}

func (u *surchargeConditionUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("Delete: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

	if u.repo == nil {
		global.Logger.Error("Delete: repository nil", zap.String(global.KeyCorrelationID, cid))
		return ErrSurchargeConditionNotFound
	}
	_, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (struct{}, error) {
			if err := u.repo.Delete(txCtx, id); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return struct{}{}, ErrSurchargeConditionNotFound
				}
				return struct{}{}, err
			}
			return struct{}{}, nil
		},
	)
	if err != nil {
		global.Logger.Error("Delete: failed to delete condition", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return err
	}
	global.Logger.Info("Delete: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	return nil
}
