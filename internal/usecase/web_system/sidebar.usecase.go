package web_system

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"go-structure/internal/common"
	dto_common "go-structure/internal/dto/common"
	dto "go-structure/internal/dto/web_system"
	"go-structure/internal/helper/database"
	"go-structure/internal/repository/model"
	sidebar_repo "go-structure/internal/repository/web_system"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	ErrSidebarNotFound = errors.New("không tìm thấy sidebar")
)

type (
	ISidebarUsecase interface {
		CreateSidebar(ctx context.Context, req *dto.CreateSidebarRequestDto) (*dto.SidebarResponseDto, error)
		GetSidebar(ctx context.Context, id uuid.UUID) (*dto.SidebarResponseDto, error)
		ListSidebars(ctx context.Context, contextFilter string, page, limit int) (*dto.ListSidebarsResponseDto, error)
		UpdateSidebar(ctx context.Context, id uuid.UUID, req *dto.UpdateSidebarRequestDto) (*dto.SidebarResponseDto, error)
		DeleteSidebar(ctx context.Context, id uuid.UUID) error
	}

	sidebarUsecase struct {
		sidebarRepo        sidebar_repo.ISidebarRepository
		transactionManager database.TransactionManager
	}
)

func NewSidebarUsecase(sidebarRepo sidebar_repo.ISidebarRepository, transactionManager database.TransactionManager) ISidebarUsecase {
	return &sidebarUsecase{
		sidebarRepo:        sidebarRepo,
		transactionManager: transactionManager,
	}
}

func (u *sidebarUsecase) CreateSidebar(ctx context.Context, req *dto.CreateSidebarRequestDto) (*dto.SidebarResponseDto, error) {
	itemsJSON, err := json.Marshal(req.Items)
	if err != nil {
		return nil, err
	}
	version := req.Version
	if version == "" {
		version = "1.0.0"
	}
	var genAt time.Time
	if req.GeneratedAt != nil {
		genAt = *req.GeneratedAt
	}
	input := &model.Sidebar{
		Context:     req.Context,
		Version:     version,
		GeneratedAt: genAt,
		Items:       itemsJSON,
	}
	var sidebar *model.Sidebar
	err = u.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		created, err := u.sidebarRepo.CreateSidebar(txCtx, input)
		if err != nil {
			return err
		}
		sidebar = created
		return nil
	})
	if err != nil {
		return nil, err
	}
	return u.sidebarToResponse(sidebar), nil
}

func (u *sidebarUsecase) GetSidebar(ctx context.Context, id uuid.UUID) (*dto.SidebarResponseDto, error) {
	sidebar, err := u.sidebarRepo.GetSidebarByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSidebarNotFound
		}
		return nil, err
	}
	return u.sidebarToResponse(sidebar), nil
}

func (u *sidebarUsecase) ListSidebars(ctx context.Context, contextFilter string, page, limit int) (*dto.ListSidebarsResponseDto, error) {
	if page < 1 {
		page = common.DefaultPage
	}
	if limit < 1 || limit > common.MaxLimit {
		limit = common.DefaultLimit
	}
	offset := int32((page - 1) * limit)
	limit32 := int32(limit)

	total, err := u.sidebarRepo.CountSidebars(ctx, contextFilter)
	if err != nil {
		return nil, err
	}
	sidebars, err := u.sidebarRepo.ListSidebars(ctx, contextFilter, limit32, offset)
	if err != nil {
		return nil, err
	}
	items := make([]dto.SidebarResponseDto, 0, len(sidebars))
	for _, s := range sidebars {
		items = append(items, *u.sidebarToResponse(s))
	}
	return &dto.ListSidebarsResponseDto{
		Items: items,
		Pagination: dto_common.PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil
}

func (u *sidebarUsecase) UpdateSidebar(ctx context.Context, id uuid.UUID, req *dto.UpdateSidebarRequestDto) (*dto.SidebarResponseDto, error) {
	itemsJSON, err := json.Marshal(req.Items)
	if err != nil {
		return nil, err
	}
	version := req.Version
	if version == "" {
		version = "1.0.0"
	}
	var genAt time.Time
	if req.GeneratedAt != nil {
		genAt = *req.GeneratedAt
	}
	input := &model.Sidebar{
		ID:          id,
		Context:     req.Context,
		Version:     version,
		GeneratedAt: genAt,
		Items:       itemsJSON,
	}
	var sidebar *model.Sidebar
	err = u.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		updated, err := u.sidebarRepo.UpdateSidebar(txCtx, input)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrSidebarNotFound
			}
			return err
		}
		sidebar = updated
		return nil
	})
	if err != nil {
		return nil, err
	}
	return u.sidebarToResponse(sidebar), nil
}

func (u *sidebarUsecase) DeleteSidebar(ctx context.Context, id uuid.UUID) error {
	return u.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		return u.sidebarRepo.DeleteSidebar(txCtx, id)
	})
}

func (u *sidebarUsecase) sidebarToResponse(s *model.Sidebar) *dto.SidebarResponseDto {
	var items []dto.SidebarItemDto
	_ = json.Unmarshal(s.Items, &items)
	var genAt *time.Time
	if !s.GeneratedAt.IsZero() {
		t := s.GeneratedAt
		genAt = &t
	}
	return &dto.SidebarResponseDto{
		ID:          s.ID.String(),
		Context:     s.Context,
		Version:     s.Version,
		GeneratedAt: genAt,
		Items:       items,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}
