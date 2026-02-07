package repository

import (
	"context"
	"net/netip"

	"go-structure/internal/dto"
	"go-structure/internal/helper/database"
	pgdb "go-structure/internal/orm/db/postgres"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	IOTPAuditRepository interface {
		CreateOTPAudit(ctx context.Context, input dto.CreateOTPAuditRequestData) error
	}

	otpAuditRepository struct {
		pool *pgxpool.Pool
	}
)

func NewOTPAuditRepository(pool *pgxpool.Pool) IOTPAuditRepository {
	return &otpAuditRepository{pool: pool}
}

func (r *otpAuditRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, r.pool)
}

func (r *otpAuditRepository) CreateOTPAudit(ctx context.Context, input dto.CreateOTPAuditRequestData) error {
	db := r.getDB(ctx)

	var ipAddr *netip.Addr
	if input.IPAddress != "" {
		addr, err := netip.ParseAddr(input.IPAddress)
		if err == nil {
			ipAddr = &addr
		}
	}

	failureReason := pgtype.Text{Valid: input.FailureReason != "", String: input.FailureReason}
	userAgent := pgtype.Text{Valid: input.UserAgent != "", String: input.UserAgent}

	_, err := db.CreateOTPAudit(ctx, pgdb.CreateOTPAuditParams{
		OtpID:         input.OTPId,
		Target:        input.Target,
		Purpose:       input.Purpose,
		AttemptNumber: int32(input.AttemptNumber),
		Result:        input.Result,
		FailureReason: failureReason,
		IpAddress:     ipAddr,
		UserAgent:     userAgent,
		Metadata:      input.Metadata,
	})
	return err
}
