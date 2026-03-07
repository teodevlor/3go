package web_system

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"go-structure/global"
	"go-structure/internal/constants"
	dto_common "go-structure/internal/dto/common"
	dto "go-structure/internal/dto/web_system"
	"go-structure/internal/helper/database"
	"go-structure/internal/middleware"
	"go-structure/internal/repository/model"
	sidebar_repo "go-structure/internal/repository/web_system"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

var (
	ErrSidebarNotFound = errors.New("không tìm thấy sidebar")
)

type (
	ISidebarUsecase interface {
		CreateSidebar(ctx context.Context, req *dto.CreateSidebarRequestDto) (*dto.SidebarResponseDto, error)
		GetSidebar(ctx context.Context, id uuid.UUID) (*dto.SidebarResponseDto, error)
		GetSidebarByContext(ctx context.Context, sidebarContext string) (*dto.SidebarResponseDto, error)
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
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("CreateSidebar: start", zap.String(global.KeyCorrelationID, cid), zap.String("context", req.Context))

	itemsJSON, err := json.Marshal(req.Items)
	if err != nil {
		global.Logger.Error("CreateSidebar: failed to marshal items", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
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

	sidebar, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (*model.Sidebar, error) {
			return u.sidebarRepo.CreateSidebar(txCtx, input)
		},
	)
	if err != nil {
		global.Logger.Error("CreateSidebar: transaction failed", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("CreateSidebar: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("sidebar_id", sidebar.ID.String()))
	return u.sidebarToResponse(sidebar), nil
}

func (u *sidebarUsecase) GetSidebar(ctx context.Context, id uuid.UUID) (*dto.SidebarResponseDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("GetSidebar: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	sidebar, err := u.sidebarRepo.GetSidebarByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			global.Logger.Error("GetSidebar: sidebar not found", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
			return nil, ErrSidebarNotFound
		}
		global.Logger.Error("GetSidebar: failed to get sidebar", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("GetSidebar: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	return u.sidebarToResponse(sidebar), nil
}

func (u *sidebarUsecase) GetSidebarByContext(ctx context.Context, sidebarContext string) (*dto.SidebarResponseDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("GetSidebarByContext: start", zap.String(global.KeyCorrelationID, cid), zap.String("context", sidebarContext))
	sidebar, err := u.sidebarRepo.GetSidebarByContext(ctx, sidebarContext)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			global.Logger.Error("GetSidebarByContext: sidebar not found", zap.String(global.KeyCorrelationID, cid), zap.String("context", sidebarContext))
			return nil, ErrSidebarNotFound
		}
		global.Logger.Error("GetSidebarByContext: failed to get sidebar", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("GetSidebarByContext: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("context", sidebarContext))
	return u.sidebarToResponse(sidebar), nil
}

func (u *sidebarUsecase) ListSidebars(ctx context.Context, contextFilter string, page, limit int) (*dto.ListSidebarsResponseDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("ListSidebars: start", zap.String(global.KeyCorrelationID, cid), zap.Int("page", page), zap.Int("limit", limit))
	if page < 1 {
		page = constants.DefaultPage
	}
	if limit < 1 || limit > constants.MaxLimit {
		limit = constants.DefaultLimit
	}
	offset := int32((page - 1) * limit)
	limit32 := int32(limit)

	total, err := u.sidebarRepo.CountSidebars(ctx, contextFilter)
	if err != nil {
		global.Logger.Error("ListSidebars: failed to count sidebars", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	sidebars, err := u.sidebarRepo.ListSidebars(ctx, contextFilter, limit32, offset)
	if err != nil {
		global.Logger.Error("ListSidebars: failed to list sidebars", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	items := make([]dto.SidebarResponseDto, 0, len(sidebars))
	for _, s := range sidebars {
		items = append(items, *u.sidebarToResponse(s))
	}
	global.Logger.Info("ListSidebars: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.Int64("total", total))
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
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("UpdateSidebar: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	itemsJSON, err := json.Marshal(req.Items)
	if err != nil {
		global.Logger.Error("UpdateSidebar: failed to marshal items", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
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

	sidebar, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (*model.Sidebar, error) {
			updated, err := u.sidebarRepo.UpdateSidebar(txCtx, input)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return nil, ErrSidebarNotFound
				}
				return nil, err
			}
			return updated, nil
		},
	)
	if err != nil {
		global.Logger.Error("UpdateSidebar: transaction failed", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("UpdateSidebar: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	return u.sidebarToResponse(sidebar), nil
}

func (u *sidebarUsecase) DeleteSidebar(ctx context.Context, id uuid.UUID) error {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("DeleteSidebar: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	_, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (struct{}, error) {
			return struct{}{}, u.sidebarRepo.DeleteSidebar(txCtx, id)
		},
	)
	if err != nil {
		global.Logger.Error("DeleteSidebar: failed to delete sidebar", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return err
	}
	global.Logger.Info("DeleteSidebar: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	return nil
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
