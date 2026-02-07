package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"go-structure/internal/helper/database"
	"go-structure/internal/mapper"
	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ISessionRepository interface {
	CreateSession(ctx context.Context, session *model.Session) (*model.Session, error)
	GetSessionByID(ctx context.Context, id uuid.UUID) (*model.Session, error)
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*model.Session, error)
	UpdateSessionActivity(ctx context.Context, id uuid.UUID) error
	RevokeSessionByID(ctx context.Context, id uuid.UUID, reason string) error
	RevokeSessionByRefreshToken(ctx context.Context, refreshToken string, reason string) error
	RevokeAllSessionsByAccount(ctx context.Context, accountID uuid.UUID, reason string) error
	RevokeAllSessionsByAccountAppDevice(ctx context.Context, accountAppDeviceID uuid.UUID, reason string) error
}

type sessionRepository struct {
	pool *pgxpool.Pool
}

func NewSessionRepository(pool *pgxpool.Pool) ISessionRepository {
	return &sessionRepository{pool: pool}
}

func (r *sessionRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, r.pool)
}

func hashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (r *sessionRepository) CreateSession(ctx context.Context, session *model.Session) (*model.Session, error) {
	db := r.getDB(ctx)

	row, err := db.CreateSession(ctx, pgdb.CreateSessionParams{
		AccountAppDeviceID: session.AccountAppDeviceID,
		RefreshTokenHash:   hashRefreshToken(session.RefreshTokenHash),
		ExpiresAt: pgtype.Timestamptz{
			Time:  session.ExpiresAt,
			Valid: true,
		},
		IpAddress: pgtype.Text{
			String: session.IpAddress,
			Valid:  session.IpAddress != "",
		},
		UserAgent: pgtype.Text{
			String: session.UserAgent,
			Valid:  session.UserAgent != "",
		},
		Metadata: session.Metadata,
	})
	if err != nil {
		return nil, err
	}

	return mapper.ToSession(row), nil
}

func (r *sessionRepository) GetSessionByID(ctx context.Context, id uuid.UUID) (*model.Session, error) {
	db := r.getDB(ctx)
	row, err := db.GetSessionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapper.ToSession(row), nil
}

func (r *sessionRepository) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*model.Session, error) {
	db := r.getDB(ctx)
	hash := hashRefreshToken(refreshToken)
	row, err := db.GetSessionByRefreshTokenHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return mapper.ToSession(row), nil
}

func (r *sessionRepository) UpdateSessionActivity(ctx context.Context, id uuid.UUID) error {
	db := r.getDB(ctx)
	return db.UpdateSessionActivity(ctx, id)
}

func (r *sessionRepository) RevokeSessionByID(ctx context.Context, id uuid.UUID, reason string) error {
	db := r.getDB(ctx)
	return db.RevokeSession(ctx, pgdb.RevokeSessionParams{
		ID: id,
		RevokedReason: pgtype.Text{
			String: reason,
			Valid:  reason != "",
		},
	})
}

func (r *sessionRepository) RevokeSessionByRefreshToken(ctx context.Context, refreshToken string, reason string) error {
	db := r.getDB(ctx)
	hash := hashRefreshToken(refreshToken)
	return db.RevokeSessionByRefreshToken(ctx, pgdb.RevokeSessionByRefreshTokenParams{
		RefreshTokenHash: hash,
		RevokedReason: pgtype.Text{
			String: reason,
			Valid:  reason != "",
		},
	})
}

func (r *sessionRepository) RevokeAllSessionsByAccount(ctx context.Context, accountID uuid.UUID, reason string) error {
	db := r.getDB(ctx)
	return db.RevokeAllSessionsByAccount(ctx, pgdb.RevokeAllSessionsByAccountParams{
		AccountID: accountID,
		RevokedReason: pgtype.Text{
			String: reason,
			Valid:  reason != "",
		},
	})
}

func (r *sessionRepository) RevokeAllSessionsByAccountAppDevice(ctx context.Context, accountAppDeviceID uuid.UUID, reason string) error {
	db := r.getDB(ctx)
	return db.RevokeAllSessionsByAccountAppDevice(ctx, pgdb.RevokeAllSessionsByAccountAppDeviceParams{
		AccountAppDeviceID: accountAppDeviceID,
		RevokedReason: pgtype.Text{
			String: reason,
			Valid:  reason != "",
		},
	})
}
