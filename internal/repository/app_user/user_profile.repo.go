package appuser

import (
	"context"
	"time"

	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type (
	IUserProfileRepository interface {
		RegisterUserProfile(ctx context.Context, userProfile *model.UserProfile) (*model.UserProfile, error)
		GetByAccountID(ctx context.Context, accountID uuid.UUID) (*model.UserProfile, error)
	}

	userProfileRepo struct {
		db *pgdb.Queries
	}
)

func NewUserProfileRepository(db *pgdb.Queries) IUserProfileRepository {
	return &userProfileRepo{db: db}
}

func (r *userProfileRepo) RegisterUserProfile(ctx context.Context, userProfile *model.UserProfile) (*model.UserProfile, error) {
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

	up, err := r.db.CreateUserProfile(ctx, params)
	if err != nil {
		return nil, err
	}

	return &model.UserProfile{
		ID:        up.ID,
		AccountID: up.AccountID,
		FullName:  up.FullName,
		AvatarURL: up.AvatarUrl.String,
		IsActive:  up.IsActive.Bool,
		Metadata:  up.Metadata,
		CreatedAt: up.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt: up.UpdatedAt.Time.Format(time.RFC3339),
	}, nil
}

func (r *userProfileRepo) GetByAccountID(ctx context.Context, accountID uuid.UUID) (*model.UserProfile, error) {
	row, err := r.db.GetUserProfileByAccountId(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return &model.UserProfile{
		ID:        row.ID,
		AccountID: row.AccountID,
		FullName:  row.FullName,
		AvatarURL: row.AvatarUrl.String,
		IsActive:  row.IsActive.Bool,
		Metadata:  row.Metadata,
		CreatedAt: row.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt: row.UpdatedAt.Time.Format(time.RFC3339),
	}, nil
}
