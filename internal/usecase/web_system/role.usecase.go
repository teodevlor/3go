package web_system

import (
	"context"
	"errors"
	"strings"

	"go-structure/internal/constants"
	dto_common "go-structure/internal/dto/common"
	dto "go-structure/internal/dto/web_system"
	"go-structure/internal/helper/database"
	pgdb "go-structure/internal/orm/db/postgres"
	websystem_model "go-structure/internal/repository/model/web_system"
	websystem_repo "go-structure/internal/repository/web_system"
	roleTransformer "go-structure/internal/transformer/web_system"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
		transactionManager database.TransactionManager
	}
)

func NewRoleUsecase(repo websystem_repo.IRoleRepository, rolePermissionRepo websystem_repo.IRolePermissionRepository, transactionManager database.TransactionManager) IRoleUsecase {
	return &roleUsecase{
		repo:               repo,
		rolePermissionRepo: rolePermissionRepo,
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
	if u.repo == nil {
		return nil, nil
	}
	permissionIDs, err := parsePermissionIDs(req.PermissionIDs)
	if err != nil {
		return nil, err
	}
	params := pgdb.CreateRoleParams{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
	}
	var role *websystem_model.Role
	err = u.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if _, err := u.repo.GetByCode(txCtx, req.Code); err == nil {
			return ErrRoleCodeExists
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		created, err := u.repo.Create(txCtx, params)
		if err != nil {
			return err
		}
		role = created
		if u.rolePermissionRepo != nil && len(permissionIDs) > 0 {
			if err := u.rolePermissionRepo.CreateRolePermissions(txCtx, role.ID, permissionIDs); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	item := roleTransformer.ToRoleItemDto(role, req.PermissionIDs)
	return &item, nil
}

func (u *roleUsecase) GetByID(ctx context.Context, id uuid.UUID) (*dto.RoleItemDto, error) {
	if u.repo == nil {
		return nil, ErrRoleNotFound
	}
	role, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRoleNotFound
		}
		return nil, err
	}
	var permIDStrs []string
	if u.rolePermissionRepo != nil {
		permIDs, _ := u.rolePermissionRepo.GetPermissionIDsByRoleID(ctx, id)
		permIDStrs = uuidSliceToStringSlice(permIDs)
	}
	item := roleTransformer.ToRoleItemDto(role, permIDStrs)
	return &item, nil
}

func (u *roleUsecase) List(ctx context.Context, page, limit int, search string) (*dto.ListRolesResponseDto, error) {
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
		return nil, err
	}
	roles, err := u.repo.List(ctx, search, limit32, offset)
	if err != nil {
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
	res := roleTransformer.ToListRolesResponseDto(roles, permMap, page, limit, total)
	return &res, nil
}

func (u *roleUsecase) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateRoleRequestDto) (*dto.RoleItemDto, error) {
	if u.repo == nil {
		return nil, ErrRoleNotFound
	}
	permIDs := *req.PermissionIDs
	permissionIDs, err := parsePermissionIDs(permIDs)
	if err != nil {
		return nil, err
	}
	params := pgdb.UpdateRoleParams{
		ID:          id,
		Code:        *req.Code,
		Name:        *req.Name,
		Description: *req.Description,
		IsActive:    *req.IsActive,
	}
	var role *websystem_model.Role
	err = u.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		existing, err := u.repo.GetByCode(txCtx, *req.Code)
		if err == nil && existing.ID != id {
			return ErrRoleCodeExists
		}
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		updated, err := u.repo.Update(txCtx, params)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrRoleNotFound
			}
			return err
		}
		role = updated
		if u.rolePermissionRepo != nil {
			if err := u.rolePermissionRepo.DeleteByRoleID(txCtx, id); err != nil {
				return err
			}
			if len(permissionIDs) > 0 {
				if err := u.rolePermissionRepo.CreateRolePermissions(txCtx, id, permissionIDs); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	item := roleTransformer.ToRoleItemDto(role, permIDs)
	return &item, nil
}

func (u *roleUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if u.repo == nil {
		return ErrRoleNotFound
	}
	return u.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		return u.repo.Delete(txCtx, id)
	})
}
