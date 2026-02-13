package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"go-structure/internal/helper/database"
	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"
	"go-structure/internal/mapper"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	ISystemAdminRefreshTokenRepository interface {
		CreateRefreshToken(ctx context.Context, refreshToken *model.SystemAdminRefreshToken) (*model.SystemAdminRefreshToken, error)
		GetRefreshTokenByHash(ctx context.Context, refreshTokenHash string) (*model.SystemAdminRefreshToken, error)
		RevokeRefreshTokenByHash(ctx context.Context, refreshTokenHash string, reason string) error
		RevokeAllRefreshTokensByAdmin(ctx context.Context, adminID uuid.UUID, reason string) error
		UpdateRefreshTokenActivity(ctx context.Context, id uuid.UUID) error
	}

	systemAdminRefreshTokenRepository struct {
		pool *pgxpool.Pool
	}
)

func NewSystemAdminRefreshTokenRepository(pool *pgxpool.Pool) ISystemAdminRefreshTokenRepository {
	return &systemAdminRefreshTokenRepository{pool: pool}
}

func (r *systemAdminRefreshTokenRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, r.pool)
}

func hashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (r *systemAdminRefreshTokenRepository) CreateRefreshToken(ctx context.Context, refreshToken *model.SystemAdminRefreshToken) (*model.SystemAdminRefreshToken, error) {
	db := r.getDB(ctx)

	params := pgdb.CreateSystemAdminRefreshTokenParams{
		AdminID:         refreshToken.AdminID,
		RefreshTokenHash: hashRefreshToken(refreshToken.RefreshTokenHash),
		ExpiresAt: pgtype.Timestamptz{
			Time:  refreshToken.ExpiresAt,
			Valid: true,
		},
		IpAddress: pgtype.Text{
			String: refreshToken.IpAddress,
			Valid:  refreshToken.IpAddress != "",
		},
		UserAgent: pgtype.Text{
			String: refreshToken.UserAgent,
			Valid:  refreshToken.UserAgent != "",
		},
		Metadata: refreshToken.Metadata,
	}

	row, err := db.CreateSystemAdminRefreshToken(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapper.ToSystemAdminRefreshToken(row), nil
}

func (r *systemAdminRefreshTokenRepository) GetRefreshTokenByHash(ctx context.Context, refreshTokenHash string) (*model.SystemAdminRefreshToken, error) {
	db := r.getDB(ctx)
	hashedToken := hashRefreshToken(refreshTokenHash)

	row, err := db.GetSystemAdminRefreshTokenByHash(ctx, hashedToken)
	if err != nil {
		return nil, err
	}

	return mapper.ToSystemAdminRefreshToken(row), nil
}

func (r *systemAdminRefreshTokenRepository) RevokeRefreshTokenByHash(ctx context.Context, refreshTokenHash string, reason string) error {
	db := r.getDB(ctx)
	hashedToken := hashRefreshToken(refreshTokenHash)
	return db.RevokeSystemAdminRefreshTokenByHash(ctx, pgdb.RevokeSystemAdminRefreshTokenByHashParams{
		RefreshTokenHash: hashedToken,
		RevokedReason: pgtype.Text{
			String: reason,
			Valid:  reason != "",
		},
	})
}

func (r *systemAdminRefreshTokenRepository) RevokeAllRefreshTokensByAdmin(ctx context.Context, adminID uuid.UUID, reason string) error {
	db := r.getDB(ctx)
	return db.RevokeAllSystemAdminRefreshTokens(ctx, pgdb.RevokeAllSystemAdminRefreshTokensParams{
		AdminID: adminID,
		RevokedReason: pgtype.Text{
			String: reason,
			Valid:  reason != "",
		},
	})
}

func (r *systemAdminRefreshTokenRepository) UpdateRefreshTokenActivity(ctx context.Context, id uuid.UUID) error {
	db := r.getDB(ctx)
	return db.UpdateSystemAdminRefreshTokenActivity(ctx, id)
}
