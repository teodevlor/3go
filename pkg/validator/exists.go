package validator

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ExistsGetterFn func(ctx context.Context, id uuid.UUID) error

func CheckExists(ctx context.Context, id uuid.UUID, getter ExistsGetterFn, notFoundErr error) error {
	if err := getter(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return notFoundErr
		}
		return err
	}
	return nil
}

func CheckExistsMany(ctx context.Context, ids []uuid.UUID, getter ExistsGetterFn, notFoundErr error) error {
	for _, id := range ids {
		if err := CheckExists(ctx, id, getter, notFoundErr); err != nil {
			return err
		}
	}
	return nil
}

func CheckExistsManyStrings(ctx context.Context, ids []string, getter ExistsGetterFn, notFoundErr error) error {
	for _, s := range ids {
		if s == "" {
			continue
		}
		id, err := uuid.Parse(s)
		if err != nil {
			return err
		}
		if err := CheckExists(ctx, id, getter, notFoundErr); err != nil {
			return err
		}
	}
	return nil
}

func CheckExistsOptionalString(ctx context.Context, s *string, getter ExistsGetterFn, notFoundErr error) error {
	if s == nil || *s == "" {
		return nil
	}
	id, err := uuid.Parse(*s)
	if err != nil {
		return err
	}
	return CheckExists(ctx, id, getter, notFoundErr)
}
