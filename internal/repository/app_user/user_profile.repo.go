package appuser

import (
	"context"
	"errors"

	"go-structure/internal/helper/database"
	appuser "go-structure/internal/mapper/app_user"
	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	IUserProfileRepository interface {
		RegisterUserProfile(ctx context.Context, userProfile *model.UserProfile) (*model.UserProfile, error)
		GetByAccountID(ctx context.Context, accountID uuid.UUID) (*model.UserProfile, error)
		GetByID(ctx context.Context, profileID uuid.UUID) (*model.UserProfile, error)
		UpdateUserProfile(ctx context.Context, userProfile *model.UserProfile) (*model.UserProfile, error)
	}

	userProfileRepo struct {
		pool *pgxpool.Pool
	}
)

func NewUserProfileRepository(pool *pgxpool.Pool) IUserProfileRepository {
	return &userProfileRepo{pool: pool}
}

// getDB trả về queries object với executor phù hợp (transaction-aware)
func (r *userProfileRepo) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, r.pool)
}

func (r *userProfileRepo) RegisterUserProfile(ctx context.Context, userProfile *model.UserProfile) (*model.UserProfile, error) {
	db := r.getDB(ctx)
	params := pgdb.CreateUserProfileParams{
		AccountID: userProfile.AccountID,
		FullName:  userProfile.FullName,
		AvatarUrl: pgtype.Text{
			String: userProfile.AvatarURL,
			Valid:  userProfile.AvatarURL != "",
		},
		IsActive: pgtype.Bool{
			Bool:  userProfile.IsActive,
			Valid: true,
		},
		Metadata: userProfile.Metadata,
	}

	up, err := db.CreateUserProfile(ctx, params)
	if err != nil {
		return nil, err
	}

	return appuser.ToUserProfile(up), nil
}

func (r *userProfileRepo) GetByAccountID(ctx context.Context, accountID uuid.UUID) (*model.UserProfile, error) {
	db := r.getDB(ctx)
	row, err := db.GetUserProfileByAccountId(ctx, accountID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return appuser.ToUserProfile(row), nil
}

func (r *userProfileRepo) GetByID(ctx context.Context, profileID uuid.UUID) (*model.UserProfile, error) {
	db := r.getDB(ctx)
	row, err := db.GetUserProfileByID(ctx, profileID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return appuser.ToUserProfile(row), nil
}

func (r *userProfileRepo) UpdateUserProfile(ctx context.Context, userProfile *model.UserProfile) (*model.UserProfile, error) {
	db := r.getDB(ctx)
	params := pgdb.UpdateUserProfileParams{
		ID:        userProfile.ID,
		FullName:  userProfile.FullName,
		AvatarUrl: pgtype.Text{String: userProfile.AvatarURL, Valid: userProfile.AvatarURL != ""},
		IsActive:  pgtype.Bool{Bool: userProfile.IsActive, Valid: true},
		Metadata:  userProfile.Metadata,
		UpdatedAt: pgtype.Timestamptz{Time: userProfile.UpdatedAt, Valid: true},
	}
	up, err := db.UpdateUserProfile(ctx, params)
	if err != nil {
		return nil, err
	}
	return appuser.ToUserProfile(up), nil
}
