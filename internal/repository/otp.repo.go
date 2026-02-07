package repository

import (
	"context"
	"errors"
	"time"

	"go-structure/internal/dto"
	"go-structure/internal/helper/database"
	pgdb "go-structure/internal/orm/db/postgres"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	IOTPRepository interface {
		CreateOTP(ctx context.Context, input dto.CreateOTPRequestData) error
		ExpireOldOTPs(ctx context.Context) error
		GetOTP(ctx context.Context, target, purpose string) (*dto.ActiveOTPResponseData, error)
		GetLastOTPCreatedAt(ctx context.Context, target, purpose string) (*time.Time, error)
		CountOTPsCreatedSince(ctx context.Context, target, purpose string, since time.Time) (int32, error)
		GetOldestOTPCreatedAtSince(ctx context.Context, target, purpose string, since time.Time) (*time.Time, error)
		MarkOTPAsUsed(ctx context.Context, id uuid.UUID) error
		IncrementOTPAttempt(ctx context.Context, id uuid.UUID) error
		LockOTP(ctx context.Context, id uuid.UUID) error
	}

	otpRepository struct {
		pool *pgxpool.Pool
	}
)

func NewOTPRepository(pool *pgxpool.Pool) IOTPRepository {
	return &otpRepository{pool: pool}
}

func (r *otpRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, r.pool)
}

func (r *otpRepository) CreateOTP(ctx context.Context, input dto.CreateOTPRequestData) error {
	db := r.getDB(ctx)
	_, err := db.CreateOTP(ctx, pgdb.CreateOTPParams{
		Target:     input.Target,
		OtpCode:    input.OtpCode,
		Purpose:    input.Purpose,
		MaxAttempt: pgtype.Int4{Int32: int32(input.MaxAttempt), Valid: true},
		ExpiresAt: pgtype.Timestamptz{
			Time:  input.ExpiresAt,
			Valid: true,
		},
	})
	return err
}

func (r *otpRepository) ExpireOldOTPs(ctx context.Context) error {
	db := r.getDB(ctx)
	return db.ExpireOldOTPs(ctx)
}

func (r *otpRepository) GetOTP(ctx context.Context, target, purpose string) (*dto.ActiveOTPResponseData, error) {
	db := r.getDB(ctx)
	row, err := db.GetActiveOTP(ctx, pgdb.GetActiveOTPParams{
		Target:  target,
		Purpose: purpose,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	attemptCount := int32(0)
	if row.AttemptCount.Valid {
		attemptCount = row.AttemptCount.Int32
	}
	maxAttempt := int32(0)
	if row.MaxAttempt.Valid {
		maxAttempt = row.MaxAttempt.Int32
	}

	return &dto.ActiveOTPResponseData{
		ID:           row.ID,
		OtpCode:      row.OtpCode,
		AttemptCount: attemptCount,
		MaxAttempt:   maxAttempt,
	}, nil
}

func (r *otpRepository) MarkOTPAsUsed(ctx context.Context, id uuid.UUID) error {
	db := r.getDB(ctx)
	return db.MarkOTPAsUsed(ctx, id)
}

func (r *otpRepository) IncrementOTPAttempt(ctx context.Context, id uuid.UUID) error {
	db := r.getDB(ctx)
	return db.IncrementOTPAttempt(ctx, id)
}

func (r *otpRepository) LockOTP(ctx context.Context, id uuid.UUID) error {
	db := r.getDB(ctx)
	return db.LockOTP(ctx, id)
}

func (r *otpRepository) GetLastOTPCreatedAt(ctx context.Context, target, purpose string) (*time.Time, error) {
	db := r.getDB(ctx)
	t, err := db.GetLastOTPCreatedAt(ctx, pgdb.GetLastOTPCreatedAtParams{Target: target, Purpose: purpose})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if !t.Valid {
		return nil, nil
	}
	return &t.Time, nil
}

func (r *otpRepository) CountOTPsCreatedSince(ctx context.Context, target, purpose string, since time.Time) (int32, error) {
	db := r.getDB(ctx)
	return db.CountOTPsCreatedSince(ctx, pgdb.CountOTPsCreatedSinceParams{
		Target:    target,
		Purpose:   purpose,
		CreatedAt: pgtype.Timestamptz{Time: since, Valid: true},
	})
}

func (r *otpRepository) GetOldestOTPCreatedAtSince(ctx context.Context, target, purpose string, since time.Time) (*time.Time, error) {
	db := r.getDB(ctx)
	t, err := db.GetOldestOTPCreatedAtSince(ctx, pgdb.GetOldestOTPCreatedAtSinceParams{
		Target:    target,
		Purpose:   purpose,
		CreatedAt: pgtype.Timestamptz{Time: since, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if !t.Valid {
		return nil, nil
	}
	return &t.Time, nil
}
