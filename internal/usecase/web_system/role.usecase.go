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
	roleTransformer "go-structure/internal/transformer/web_system"
	"go-structure/pkg/validator"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

var (
	ErrRoleNotFound   = errors.New("không tìm thấy vai trò")
	ErrRoleCodeExists = errors.New("mã vai trò đã tồn tại")
)

type (
	IRoleUsecase interface {
		Create(ctx context.Context, req *dto.CreateRoleRequestDto) (*dto.RoleItemDto, error)
		GetByID(ctx context.Context, id uuid.UUID) (*dto.RoleItemDto, error)
		List(ctx context.Context, page, limit int, search string) (*dto.ListRolesResponseDto, error)
		Update(ctx context.Context, id uuid.UUID, req *dto.UpdateRoleRequestDto) (*dto.RoleItemDto, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	roleUsecase struct {
		repo               websystem_repo.IRoleRepository
		rolePermissionRepo websystem_repo.IRolePermissionRepository
		permissionRepo     websystem_repo.IPermissionRepository
		transactionManager database.TransactionManager
	}
)

func NewRoleUsecase(
	repo websystem_repo.IRoleRepository,
	rolePermissionRepo websystem_repo.IRolePermissionRepository,
	permissionRepo websystem_repo.IPermissionRepository,
	transactionManager database.TransactionManager,
) IRoleUsecase {
	return &roleUsecase{
		repo:               repo,
		rolePermissionRepo: rolePermissionRepo,
		permissionRepo:     permissionRepo,
		transactionManager: transactionManager,
	}
}

func parsePermissionIDs(ss []string) ([]uuid.UUID, error) {
	if len(ss) == 0 {
		return nil, nil
	}
	ids := make([]uuid.UUID, 0, len(ss))
	for _, s := range ss {
		id, err := uuid.Parse(s)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func uuidSliceToStringSlice(uuids []uuid.UUID) []string {
	if len(uuids) == 0 {
		return []string{}
	}
	out := make([]string, 0, len(uuids))
	for _, u := range uuids {
		out = append(out, u.String())
	}
	return out
}

func (u *roleUsecase) Create(ctx context.Context, req *dto.CreateRoleRequestDto) (*dto.RoleItemDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("Create: start", zap.String(global.KeyCorrelationID, cid), zap.String("code", req.Code))

	if u.repo == nil {
		return nil, nil
	}
	permissionIDs, err := parsePermissionIDs(req.PermissionIDs)
	if err != nil {
		global.Logger.Error("Create: failed to parse permission IDs", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	if u.permissionRepo != nil {
		if err := validator.CheckExistsMany(ctx, permissionIDs, func(ctx context.Context, id uuid.UUID) error {
			_, err := u.permissionRepo.GetByID(ctx, id)
			return err
		}, ErrPermissionNotFound); err != nil {
			global.Logger.Error("Create: failed to validate permissions", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
			return nil, err
		}
	}
	params := pgdb.CreateRoleParams{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	role, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (*websystem_model.Role, error) {
			if _, err := u.repo.GetByCode(txCtx, req.Code); err == nil {
				return nil, ErrRoleCodeExists
			} else if !errors.Is(err, pgx.ErrNoRows) {
				return nil, err
			}
			created, err := u.repo.Create(txCtx, params)
			if err != nil {
				return nil, err
			}
			if u.rolePermissionRepo != nil && len(permissionIDs) > 0 {
				if err := u.rolePermissionRepo.CreateRolePermissions(txCtx, created.ID, permissionIDs); err != nil {
					return nil, err
				}
			}
			return created, nil
		},
	)
	if err != nil {
		global.Logger.Error("Create: transaction failed", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("Create: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("role_id", role.ID.String()))
	item := roleTransformer.ToRoleItemDto(role, req.PermissionIDs)
	return &item, nil
}

func (u *roleUsecase) GetByID(ctx context.Context, id uuid.UUID) (*dto.RoleItemDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("GetByID: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

	if u.repo == nil {
		global.Logger.Error("GetByID: repository nil", zap.String(global.KeyCorrelationID, cid))
		return nil, ErrRoleNotFound
	}
	role, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			global.Logger.Error("GetByID: role not found", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
			return nil, ErrRoleNotFound
		}
		global.Logger.Error("GetByID: failed to get role", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	var permIDStrs []string
	if u.rolePermissionRepo != nil {
		permIDs, _ := u.rolePermissionRepo.GetPermissionIDsByRoleID(ctx, id)
		permIDStrs = uuidSliceToStringSlice(permIDs)
	}
	global.Logger.Info("GetByID: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	item := roleTransformer.ToRoleItemDto(role, permIDStrs)
	return &item, nil
}

func (u *roleUsecase) List(ctx context.Context, page, limit int, search string) (*dto.ListRolesResponseDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("List: start", zap.String(global.KeyCorrelationID, cid), zap.Int("page", page), zap.Int("limit", limit))
	if u.repo == nil {
		return &dto.ListRolesResponseDto{
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
	if limit < 1 || limit > constants.MaxLimit {
		limit = constants.DefaultLimit
	}
	search = strings.TrimSpace(search)
	if len(search) > constants.MaxSearchLen {
		search = search[:constants.MaxSearchLen]
	}
	offset := int32((page - 1) * limit)
	limit32 := int32(limit)

	total, err := u.repo.Count(ctx, search)
	if err != nil {
		global.Logger.Error("List: failed to count roles", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	roles, err := u.repo.List(ctx, search, limit32, offset)
	if err != nil {
		global.Logger.Error("List: failed to list roles", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	roleIDs := make([]uuid.UUID, 0, len(roles))
	for _, r := range roles {
		roleIDs = append(roleIDs, r.ID)
	}
	var permMap map[uuid.UUID][]uuid.UUID
	if u.rolePermissionRepo != nil && len(roleIDs) > 0 {
		permMap, _ = u.rolePermissionRepo.GetPermissionIDsByRoleIDs(ctx, roleIDs)
	}
	global.Logger.Info("List: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.Int64("total", total))
	res := roleTransformer.ToListRolesResponseDto(roles, permMap, page, limit, total)
	return &res, nil
}

func (u *roleUsecase) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateRoleRequestDto) (*dto.RoleItemDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("Update: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

	if u.repo == nil {
		global.Logger.Error("Update: repository nil", zap.String(global.KeyCorrelationID, cid))
		return nil, ErrRoleNotFound
	}
	permIDs := *req.PermissionIDs
	permissionIDs, err := parsePermissionIDs(permIDs)
	if err != nil {
		global.Logger.Error("Update: failed to parse permission IDs", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	if u.permissionRepo != nil {
		if err := validator.CheckExistsMany(ctx, permissionIDs, func(ctx context.Context, id uuid.UUID) error {
			_, err := u.permissionRepo.GetByID(ctx, id)
			return err
		}, ErrPermissionNotFound); err != nil {
			global.Logger.Error("Update: failed to validate permissions", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
			return nil, err
		}
	}
	params := pgdb.UpdateRoleParams{
		ID:          id,
		Code:        *req.Code,
		Name:        *req.Name,
		Description: *req.Description,
		IsActive:    *req.IsActive,
	}

	role, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (*websystem_model.Role, error) {
			existing, err := u.repo.GetByCode(txCtx, *req.Code)
			if err == nil && existing.ID != id {
				return nil, ErrRoleCodeExists
			}
			if err != nil && !errors.Is(err, pgx.ErrNoRows) {
				return nil, err
			}
			updated, err := u.repo.Update(txCtx, params)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return nil, ErrRoleNotFound
				}
				return nil, err
			}
			if u.rolePermissionRepo != nil {
				if err := u.rolePermissionRepo.DeleteByRoleID(txCtx, id); err != nil {
					return nil, err
				}
				if len(permissionIDs) > 0 {
					if err := u.rolePermissionRepo.CreateRolePermissions(txCtx, id, permissionIDs); err != nil {
						return nil, err
					}
				}
			}
			return updated, nil
		},
	)
	if err != nil {
		global.Logger.Error("Update: transaction failed", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("Update: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	item := roleTransformer.ToRoleItemDto(role, permIDs)
	return &item, nil
}

func (u *roleUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("Delete: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

	if u.repo == nil {
		global.Logger.Error("Delete: repository nil", zap.String(global.KeyCorrelationID, cid))
		return ErrRoleNotFound
	}
	_, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (struct{}, error) {
			return struct{}{}, u.repo.Delete(txCtx, id)
		},
	)
	if err != nil {
		global.Logger.Error("Delete: failed to delete role", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return err
	}
	global.Logger.Info("Delete: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	return nil
}
