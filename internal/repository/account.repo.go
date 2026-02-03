package repository

import (
	"context"
	"errors"

	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type (
	IAccountRepository interface {
		// GetByPhone trả về account theo số điện thoại, hoặc nil nếu chưa tồn tại.
		GetByPhone(ctx context.Context, phone string) (*model.Account, error)
		CreateAccount(ctx context.Context, account *model.Account) (*model.Account, error)
	}

	accountRepository struct {
		db *pgdb.Queries
	}
)

func NewAccountRepository(db *pgdb.Queries) IAccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) CreateAccount(ctx context.Context, account *model.Account) (*model.Account, error) {
	params := pgdb.CreateAccountParams{
		Email: pgtype.Text{
			String: account.Email,
			Valid:  account.Email != "",
		},
		PasswordHash: account.PasswordHash,
		Phone:        account.Phone,
	}

	acc, err := r.db.CreateAccount(ctx, params)
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
	acc, err := r.db.GetAccountByPhone(ctx, phone)
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
