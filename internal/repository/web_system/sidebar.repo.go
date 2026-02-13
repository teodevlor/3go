package repository

import (
	"context"
	"errors"

	"go-structure/internal/helper/database"
	webmapper "go-structure/internal/mapper/web_system"
	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	ISidebarRepository interface {
		CreateSidebar(ctx context.Context, sidebar *model.Sidebar) (*model.Sidebar, error)
		GetSidebarByID(ctx context.Context, id uuid.UUID) (*model.Sidebar, error)
		ListSidebars(ctx context.Context, contextFilter string, limit, offset int32) ([]*model.Sidebar, error)
		CountSidebars(ctx context.Context, contextFilter string) (int64, error)
		UpdateSidebar(ctx context.Context, sidebar *model.Sidebar) (*model.Sidebar, error)
		DeleteSidebar(ctx context.Context, id uuid.UUID) error
	}

	sidebarRepository struct {
		pool *pgxpool.Pool
	}
)

func NewSidebarRepository(pool *pgxpool.Pool) ISidebarRepository {
	return &sidebarRepository{pool: pool}
}

func (r *sidebarRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, r.pool)
}

func (r *sidebarRepository) CreateSidebar(ctx context.Context, sidebar *model.Sidebar) (*model.Sidebar, error) {
	db := r.getDB(ctx)
	var genAt pgtype.Timestamptz
	if !sidebar.GeneratedAt.IsZero() {
		genAt = pgtype.Timestamptz{Time: sidebar.GeneratedAt, Valid: true}
	}
	row, err := db.CreateSidebar(ctx, pgdb.CreateSidebarParams{
		Context:     sidebar.Context,
		Version:     sidebar.Version,
		GeneratedAt: genAt,
		Items:       sidebar.Items,
	})
	if err != nil {
		return nil, err
	}
	return webmapper.ToSidebar(row), nil
}

func (r *sidebarRepository) GetSidebarByID(ctx context.Context, id uuid.UUID) (*model.Sidebar, error) {
	db := r.getDB(ctx)
	row, err := db.GetSidebarByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}
	return webmapper.ToSidebar(row), nil
}

func (r *sidebarRepository) ListSidebars(ctx context.Context, contextFilter string, limit, offset int32) ([]*model.Sidebar, error) {
	db := r.getDB(ctx)
	rows, err := db.ListSidebars(ctx, pgdb.ListSidebarsParams{
		Column1: contextFilter,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, err
	}
	out := make([]*model.Sidebar, 0, len(rows))
	for _, row := range rows {
		out = append(out, webmapper.ToSidebar(row))
	}
	return out, nil
}

func (r *sidebarRepository) CountSidebars(ctx context.Context, contextFilter string) (int64, error) {
	db := r.getDB(ctx)
	return db.CountSidebars(ctx, contextFilter)
}

func (r *sidebarRepository) UpdateSidebar(ctx context.Context, sidebar *model.Sidebar) (*model.Sidebar, error) {
	db := r.getDB(ctx)
	var genAt pgtype.Timestamptz
	if !sidebar.GeneratedAt.IsZero() {
		genAt = pgtype.Timestamptz{Time: sidebar.GeneratedAt, Valid: true}
	}
	row, err := db.UpdateSidebar(ctx, pgdb.UpdateSidebarParams{
		ID:          sidebar.ID,
		Context:     sidebar.Context,
		Version:     sidebar.Version,
		GeneratedAt: genAt,
		Items:       sidebar.Items,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}
	return webmapper.ToSidebar(row), nil
}

func (r *sidebarRepository) DeleteSidebar(ctx context.Context, id uuid.UUID) error {
	db := r.getDB(ctx)
	return db.DeleteSidebar(ctx, id)
}
