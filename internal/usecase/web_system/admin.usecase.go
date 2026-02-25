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
	"go-structure/internal/repository/model"
	websystem_repo "go-structure/internal/repository/web_system"
	adminTransformer "go-structure/internal/transformer/web_system"
	"go-structure/pkg/validator"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var ErrAdminEmailUsed = errors.New("email admin đã được sử dụng")

type (
	IAdminUsecase interface {
		Create(ctx context.Context, req *dto.CreateAdminRequestDto) (*dto.AdminItemDto, error)
		GetByID(ctx context.Context, id uuid.UUID) (*dto.AdminItemDto, error)
		List(ctx context.Context, page, limit int, search string) (*dto.ListAdminsResponseDto, error)
		Update(ctx context.Context, id uuid.UUID, req *dto.UpdateAdminRequestDto) (*dto.AdminItemDto, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	adminUsecase struct {
		adminRepo      websystem_repo.ISystemAdminRepository
		adminRoleRepo  websystem_repo.ISystemAdminRoleRepository
		roleRepo       websystem_repo.IRoleRepository
		transactionMgr database.TransactionManager
	}
)

func NewAdminUsecase(
	adminRepo websystem_repo.ISystemAdminRepository,
	adminRoleRepo websystem_repo.ISystemAdminRoleRepository,
	roleRepo websystem_repo.IRoleRepository,
	transactionMgr database.TransactionManager,
) IAdminUsecase {
	return &adminUsecase{
		adminRepo:      adminRepo,
		adminRoleRepo:  adminRoleRepo,
		roleRepo:       roleRepo,
		transactionMgr: transactionMgr,
	}
}

func (u *adminUsecase) getRoleItemDtos(ctx context.Context, roleIDs []uuid.UUID) []dto.RoleItemDto {
	if len(roleIDs) == 0 {
		return nil
	}
	roles, err := u.roleRepo.GetByIDs(ctx, roleIDs)
	if err != nil || len(roles) == 0 {
		return nil
	}
	out := make([]dto.RoleItemDto, 0, len(roles))
	for _, r := range roles {
		out = append(out, adminTransformer.ToRoleItemDto(r, nil))
	}
	return out
}

func (u *adminUsecase) validateRoleIDs(ctx context.Context, roleIDs []string) ([]uuid.UUID, error) {
	if len(roleIDs) == 0 {
		return nil, nil
	}
	ids := make([]uuid.UUID, 0, len(roleIDs))
	for _, s := range roleIDs {
		id, err := uuid.Parse(s)
		if err != nil {
			return nil, ErrRoleNotFound
		}
		_, err = u.roleRepo.GetByID(ctx, id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, ErrRoleNotFound
			}
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (u *adminUsecase) Create(ctx context.Context, req *dto.CreateAdminRequestDto) (*dto.AdminItemDto, error) {
	if u.adminRepo == nil {
		return nil, nil
	}
	roleUUIDs, err := u.validateRoleIDs(ctx, req.RoleIDs)
	if err != nil {
		return nil, err
	}
	hashedPassword, err := validator.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	fullName := pgtype.Text{String: req.FullName, Valid: true}
	isActive := pgtype.Bool{Bool: req.IsActive, Valid: true}
	params := pgdb.CreateSystemAdminParams{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     fullName,
		Department:   req.Department,
		IsActive:     isActive,
	}

	admin, err := database.WithTransaction(
		u.transactionMgr,
		ctx,
		func(txCtx context.Context) (*model.SystemAdmin, error) {
			if _, err := u.adminRepo.GetByEmail(txCtx, req.Email); err == nil {
				return nil, ErrAdminEmailUsed
			} else if !errors.Is(err, pgx.ErrNoRows) {
				return nil, err
			}
			created, err := u.adminRepo.Create(txCtx, params)
			if err != nil {
				return nil, err
			}
			if len(roleUUIDs) > 0 {
				if err := u.adminRoleRepo.SetAdminRoles(txCtx, created.ID, roleUUIDs, nil); err != nil {
					return nil, err
				}
			}
			return created, nil
		},
	)
	if err != nil {
		return nil, err
	}
	roleIDs, _ := u.adminRoleRepo.GetRoleIDsByAdminID(ctx, admin.ID)
	roles := u.getRoleItemDtos(ctx, roleIDs)
	item := adminTransformer.ToAdminItemDto(admin, roles)
	return &item, nil
}

func (u *adminUsecase) GetByID(ctx context.Context, id uuid.UUID) (*dto.AdminItemDto, error) {
	if u.adminRepo == nil {
		return nil, ErrAdminNotFound
	}
	admin, err := u.adminRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAdminNotFound
		}
		return nil, err
	}
	roleIDs, _ := u.adminRoleRepo.GetRoleIDsByAdminID(ctx, id)
	roles := u.getRoleItemDtos(ctx, roleIDs)
	item := adminTransformer.ToAdminItemDto(admin, roles)
	return &item, nil
}

func (u *adminUsecase) List(ctx context.Context, page, limit int, search string) (*dto.ListAdminsResponseDto, error) {
	if u.adminRepo == nil {
		return &dto.ListAdminsResponseDto{Items: nil, Pagination: dto_common.PaginationMeta{Page: page, Limit: limit, Total: 0}}, nil
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

	total, err := u.adminRepo.Count(ctx, search)
	if err != nil {
		return nil, err
	}
	admins, err := u.adminRepo.List(ctx, search, limit32, offset)
	if err != nil {
		return nil, err
	}
	items := make([]dto.AdminItemDto, 0, len(admins))
	for _, a := range admins {
		roleIDs, _ := u.adminRoleRepo.GetRoleIDsByAdminID(ctx, a.ID)
		roles := u.getRoleItemDtos(ctx, roleIDs)
		items = append(items, adminTransformer.ToAdminItemDto(a, roles))
	}
	return &dto.ListAdminsResponseDto{
		Items: items,
		Pagination: dto_common.PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil
}

func (u *adminUsecase) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateAdminRequestDto) (*dto.AdminItemDto, error) {
	if u.adminRepo == nil {
		return nil, ErrAdminNotFound
	}
	roleIDStrs := *req.RoleIDs
	roleUUIDs, err := u.validateRoleIDs(ctx, roleIDStrs)
	if err != nil {
		return nil, err
	}
	fullName := pgtype.Text{String: *req.FullName, Valid: true}
	isActive := pgtype.Bool{Bool: *req.IsActive, Valid: true}
	params := pgdb.UpdateSystemAdminParams{
		ID:         id,
		Email:      *req.Email,
		FullName:   fullName,
		Department: *req.Department,
		IsActive:   isActive,
	}

	admin, err := database.WithTransaction(
		u.transactionMgr,
		ctx,
		func(txCtx context.Context) (*model.SystemAdmin, error) {
			existing, err := u.adminRepo.GetByEmail(txCtx, *req.Email)
			if err == nil && existing.ID != id {
				return nil, ErrAdminEmailUsed
			}
			if err != nil && !errors.Is(err, pgx.ErrNoRows) {
				return nil, err
			}
			updated, err := u.adminRepo.Update(txCtx, params)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return nil, ErrAdminNotFound
				}
				return nil, err
			}
			if err := u.adminRoleRepo.SetAdminRoles(txCtx, id, roleUUIDs, nil); err != nil {
				return nil, err
			}
			return updated, nil
		},
	)
	if err != nil {
		return nil, err
	}
	adminRoleIDs, _ := u.adminRoleRepo.GetRoleIDsByAdminID(ctx, admin.ID)
	roles := u.getRoleItemDtos(ctx, adminRoleIDs)
	item := adminTransformer.ToAdminItemDto(admin, roles)
	return &item, nil
}

func (u *adminUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if u.adminRepo == nil {
		return ErrAdminNotFound
	}
	_, err := database.WithTransaction(
		u.transactionMgr,
		ctx,
		func(txCtx context.Context) (struct{}, error) {
			return struct{}{}, u.adminRepo.Delete(txCtx, id)
		},
	)
	return err
}
