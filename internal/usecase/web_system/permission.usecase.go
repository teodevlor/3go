package web_system

import (
	"context"
	"errors"
	"strings"

	"go-structure/global"
	"go-structure/internal/constants"
	dto_common "go-structure/internal/dto/common"
	dto "go-structure/internal/dto/web_system"
	"go-structure/internal/helper/database"
	"go-structure/internal/middleware"
	pgdb "go-structure/orm/db/postgres"
	websystem_model "go-structure/internal/repository/model/web_system"
	websystem_repo "go-structure/internal/repository/web_system"
	permissionTransformer "go-structure/internal/transformer/web_system"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

var (
	ErrPermissionNotFound   = errors.New("không tìm thấy quyền")
	ErrPermissionCodeExists = errors.New("quyền với resource và action này đã tồn tại")
)

func permissionCode(resource, action string) string {
	return resource + "." + action
}

type (
	IPermissionUsecase interface {
		Create(ctx context.Context, req *dto.CreatePermissionRequestDto) (*dto.PermissionItemDto, error)
		GetByID(ctx context.Context, id uuid.UUID) (*dto.PermissionItemDto, error)
		List(ctx context.Context, page, limit int, search string) (*dto.ListPermissionsResponseDto, error)
		Update(ctx context.Context, id uuid.UUID, req *dto.UpdatePermissionRequestDto) (*dto.PermissionItemDto, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	permissionUsecase struct {
		repo               websystem_repo.IPermissionRepository
		transactionManager database.TransactionManager
	}
)

func NewPermissionUsecase(repo websystem_repo.IPermissionRepository, transactionManager database.TransactionManager) IPermissionUsecase {
	return &permissionUsecase{
		repo:               repo,
		transactionManager: transactionManager,
	}
}

func (u *permissionUsecase) Create(ctx context.Context, req *dto.CreatePermissionRequestDto) (*dto.PermissionItemDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("Create: start", zap.String(global.KeyCorrelationID, cid), zap.String("resource", req.Resource), zap.String("action", req.Action))

	if u.repo == nil {
		return nil, nil
	}
	params := pgdb.CreatePermissionParams{
		Resource: req.Resource,
		Action:   req.Action,
		Name:     req.Name,
		Description: pgtype.Text{
			String: req.Description,
			Valid:  req.Description != "",
		},
	}

	perm, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (*websystem_model.Permission, error) {
			code := permissionCode(req.Resource, req.Action)
			if _, err := u.repo.GetByCode(txCtx, code); err == nil {
				return nil, ErrPermissionCodeExists
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
	global.Logger.Info("Create: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", perm.ID.String()))
	item := permissionTransformer.ToPermissionItemDto(perm)
	return &item, nil
}

func (u *permissionUsecase) GetByID(ctx context.Context, id uuid.UUID) (*dto.PermissionItemDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("GetByID: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

	if u.repo == nil {
		global.Logger.Error("GetByID: repository nil", zap.String(global.KeyCorrelationID, cid))
		return nil, ErrPermissionNotFound
	}
	perm, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			global.Logger.Error("GetByID: permission not found", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
			return nil, ErrPermissionNotFound
		}
		global.Logger.Error("GetByID: failed to get permission", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("GetByID: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	item := permissionTransformer.ToPermissionItemDto(perm)
	return &item, nil
}

func (u *permissionUsecase) List(ctx context.Context, page, limit int, search string) (*dto.ListPermissionsResponseDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("List: start", zap.String(global.KeyCorrelationID, cid), zap.Int("page", page), zap.Int("limit", limit))

	if u.repo == nil {
		return &dto.ListPermissionsResponseDto{
			Items: nil,
			Pagination: dto_common.PaginationMeta{
				Page:  page,
				Limit: limit,
				Total: 0,
			},
		}, nil
	}
	if page < 1 {
		page = constants.DefaultPage
	}
	if limit < 1 {
		limit = constants.DefaultLimit
	}
	if limit != constants.LimitAll && limit > constants.MaxLimit {
		limit = constants.MaxLimit
	}
	search = strings.TrimSpace(search)
	if len(search) > constants.MaxSearchLen {
		search = search[:constants.MaxSearchLen]
	}
	offset := int32((page - 1) * limit)
	limit32 := int32(limit)

	total, err := u.repo.Count(ctx, search)
	if err != nil {
		global.Logger.Error("List: failed to count", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	perms, err := u.repo.List(ctx, search, limit32, offset)
	if err != nil {
		global.Logger.Error("List: failed to list", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	items := make([]dto.PermissionItemDto, 0, len(perms))
	for _, p := range perms {
		items = append(items, permissionTransformer.ToPermissionItemDto(p))
	}
	global.Logger.Info("List: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.Int64("total", total))
	return &dto.ListPermissionsResponseDto{
		Items: items,
		Pagination: dto_common.PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil
}

func (u *permissionUsecase) Update(ctx context.Context, id uuid.UUID, req *dto.UpdatePermissionRequestDto) (*dto.PermissionItemDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("Update: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

	if u.repo == nil {
		global.Logger.Error("Update: repository nil", zap.String(global.KeyCorrelationID, cid))
		return nil, ErrPermissionNotFound
	}
	desc := ""
	if req.Description != nil {
		desc = *req.Description
	}
	params := pgdb.UpdatePermissionParams{
		ID:       id,
		Resource: *req.Resource,
		Action:   *req.Action,
		Name:     *req.Name,
		Description: pgtype.Text{
			String: desc,
			Valid:  desc != "",
		},
	}

	perm, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (*websystem_model.Permission, error) {
			code := permissionCode(*req.Resource, *req.Action)
			existing, err := u.repo.GetByCode(txCtx, code)
			if err == nil && existing.ID != id {
				return nil, ErrPermissionCodeExists
			}
			if err != nil && !errors.Is(err, pgx.ErrNoRows) {
				return nil, err
			}
			updated, err := u.repo.Update(txCtx, params)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return nil, ErrPermissionNotFound
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
	item := permissionTransformer.ToPermissionItemDto(perm)
	return &item, nil
}

func (u *permissionUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("Delete: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

	if u.repo == nil {
		global.Logger.Error("Delete: repository nil", zap.String(global.KeyCorrelationID, cid))
		return ErrPermissionNotFound
	}
	_, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (struct{}, error) {
			return struct{}{}, u.repo.Delete(txCtx, id)
		},
	)
	if err != nil {
		global.Logger.Error("Delete: failed to delete permission", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return err
	}
	global.Logger.Info("Delete: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	return nil
}
