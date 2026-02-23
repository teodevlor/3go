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
	permissionTransformer "go-structure/internal/transformer/web_system"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	if u.repo == nil {
		return nil, nil
	}
	params := pgdb.CreatePermissionParams{
		Resource:    req.Resource,
		Action:      req.Action,
		Name:        req.Name,
		Description: req.Description,
	}
	var perm *websystem_model.Permission
	err := u.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		code := permissionCode(req.Resource, req.Action)
		if _, err := u.repo.GetByCode(txCtx, code); err == nil {
			return ErrPermissionCodeExists
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		created, err := u.repo.Create(txCtx, params)
		if err != nil {
			return err
		}
		perm = created
		return nil
	})
	if err != nil {
		return nil, err
	}
	item := permissionTransformer.ToPermissionItemDto(perm)
	return &item, nil
}

func (u *permissionUsecase) GetByID(ctx context.Context, id uuid.UUID) (*dto.PermissionItemDto, error) {
	if u.repo == nil {
		return nil, ErrPermissionNotFound
	}
	perm, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPermissionNotFound
		}
		return nil, err
	}
	item := permissionTransformer.ToPermissionItemDto(perm)
	return &item, nil
}

func (u *permissionUsecase) List(ctx context.Context, page, limit int, search string) (*dto.ListPermissionsResponseDto, error) {
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
		return nil, err
	}
	perms, err := u.repo.List(ctx, search, limit32, offset)
	if err != nil {
		return nil, err
	}
	items := make([]dto.PermissionItemDto, 0, len(perms))
	for _, p := range perms {
		items = append(items, permissionTransformer.ToPermissionItemDto(p))
	}
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
	if u.repo == nil {
		return nil, ErrPermissionNotFound
	}
	params := pgdb.UpdatePermissionParams{
		ID:          id,
		Resource:    *req.Resource,
		Action:      *req.Action,
		Name:        *req.Name,
		Description: *req.Description,
	}
	var perm *websystem_model.Permission
	err := u.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		code := permissionCode(*req.Resource, *req.Action)
		existing, err := u.repo.GetByCode(txCtx, code)
		if err == nil && existing.ID != id {
			return ErrPermissionCodeExists
		}
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		updated, err := u.repo.Update(txCtx, params)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrPermissionNotFound
			}
			return err
		}
		perm = updated
		return nil
	})
	if err != nil {
		return nil, err
	}
	item := permissionTransformer.ToPermissionItemDto(perm)
	return &item, nil
}

func (u *permissionUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if u.repo == nil {
		return ErrPermissionNotFound
	}
	return u.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		return u.repo.Delete(txCtx, id)
	})
}
