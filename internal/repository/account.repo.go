package repository

import (
	"context"
	"errors"

	"go-structure/internal/helper/database"
	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	IAccountRepository interface {
		GetByPhone(ctx context.Context, phone string) (*model.Account, error)
		GetById(ctx context.Context, id string) (*model.Account, error)
		CreateAccount(ctx context.Context, account *model.Account) (*model.Account, error)
		UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
	}

	accountRepository struct {
		pool *pgxpool.Pool
	}
)

func NewAccountRepository(pool *pgxpool.Pool) IAccountRepository {
	return &accountRepository{pool: pool}
}

// getDB trả về queries object với executor phù hợp (transaction-aware)
func (r *accountRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, r.pool)
}

func (r *accountRepository) CreateAccount(ctx context.Context, account *model.Account) (*model.Account, error) {
	db := r.getDB(ctx)
	params := pgdb.CreateAccountParams{
		Email: pgtype.Text{
			String: account.Email,
			Valid:  account.Email != "",
		},
		PasswordHash: account.PasswordHash,
		Phone:        account.Phone,
	}

	acc, err := db.CreateAccount(ctx, params)
	if err != nil {
		return nil, err
	}

	return &model.Account{
		ID:           acc.ID,
		Email:        acc.Email.String,
		Phone:        acc.Phone,
		PasswordHash: acc.PasswordHash,
	}, nil
}

func (r *accountRepository) GetByPhone(ctx context.Context, phone string) (*model.Account, error) {
	db := r.getDB(ctx)
	acc, err := db.GetAccountByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &model.Account{
		ID:           acc.ID,
		Email:        acc.Email.String,
		Phone:        acc.Phone,
		PasswordHash: acc.PasswordHash,
	}, nil
}

func (r *accountRepository) GetById(ctx context.Context, id string) (*model.Account, error) {
	db := r.getDB(ctx)
	uuidId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	acc, err := db.GetAccountByID(ctx, uuidId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &model.Account{
		ID:           acc.ID,
		Email:        acc.Email.String,
		Phone:        acc.Phone,
		PasswordHash: acc.PasswordHash,
	}, nil
}

func (r *accountRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	db := r.getDB(ctx)
	params := pgdb.UpdatePasswordParams{
		ID:           id,
		PasswordHash: passwordHash,
	}
	return db.UpdatePassword(ctx, params)
}
