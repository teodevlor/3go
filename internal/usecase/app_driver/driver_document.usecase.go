package app_driver

import (
	"context"
	"errors"

	common "go-structure/internal/common"
	dto "go-structure/internal/dto/app_driver"
	"go-structure/internal/helper/database"
	pgdb "go-structure/internal/orm/db/postgres"
	appdriverrepo "go-structure/internal/repository/app_driver"
	appdrivermodel "go-structure/internal/repository/model/app_driver"
	appdrivertransformer "go-structure/internal/transformer/app_driver"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrDriverDocumentNotFound = errors.New("không tìm thấy giấy tờ tài xế")
)

type (
	IDriverDocumentUsecase interface {
		Create(ctx context.Context, req *dto.CreateDriverDocumentRequestDto) (*dto.CreateDriverDocumentResponseDto, error)
		BulkCreate(ctx context.Context, req *dto.BulkCreateDriverDocumentsRequestDto) (*dto.BulkCreateDriverDocumentsResponseDto, error)
		GetByID(ctx context.Context, id uuid.UUID) (*dto.DriverDocumentItemDto, error)
		ListByDriverID(ctx context.Context, driverID uuid.UUID) (*dto.ListDriverDocumentsResponseDto, error)
		Update(ctx context.Context, id uuid.UUID, req *dto.UpdateDriverDocumentRequestDto) (*dto.DriverDocumentItemDto, error)
		BulkUpdate(ctx context.Context, req *dto.BulkUpdateDriverDocumentsRequestDto) (*dto.BulkUpdateDriverDocumentsResponseDto, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	driverDocumentUsecase struct {
		repo      appdriverrepo.IDriverDocumentRepository
		txManager database.TransactionManager
	}
)

func NewDriverDocumentUsecase(repo appdriverrepo.IDriverDocumentRepository, txManager database.TransactionManager) IDriverDocumentUsecase {
	return &driverDocumentUsecase{repo: repo, txManager: txManager}
}

func parseDate(s *string) pgtype.Date {
	if s == nil || *s == "" {
		return pgtype.Date{}
	}
	t, err := common.ParseYYYYMMDDToTime(*s)
	if err != nil {
		return pgtype.Date{}
	}
	return pgtype.Date{Time: t, Valid: true}
}

func parseStatus(s string) pgdb.DriverDocumentStatus {
	switch s {
	case dto.DriverDocumentStatusAPPROVED:
		return pgdb.DriverDocumentStatusAPPROVED
	case dto.DriverDocumentStatusREJECTED:
		return pgdb.DriverDocumentStatusREJECTED
	default:
		return pgdb.DriverDocumentStatusPENDING
	}
}

func (u *driverDocumentUsecase) Create(ctx context.Context, req *dto.CreateDriverDocumentRequestDto) (*dto.CreateDriverDocumentResponseDto, error) {
	driverID, err := uuid.Parse(req.DriverID)
	if err != nil {
		return nil, err
	}
	documentTypeID, err := uuid.Parse(req.DocumentTypeID)
	if err != nil {
		return nil, err
	}
	arg := pgdb.CreateDriverDocumentParams{
		DriverID:       driverID,
		DocumentTypeID: documentTypeID,
		FileUrl:        req.FileUrl,
		ExpireAt:       parseDate(req.ExpireAt),
		Status:         pgdb.DriverDocumentStatusPENDING,
	}

	created, err := database.WithTransaction(
		u.txManager,
		ctx,
		func(txCtx context.Context) (*appdrivermodel.DriverDocument, error) {
			return u.repo.Create(txCtx, arg)
		},
	)
	if err != nil {
		return nil, err
	}
	res := appdrivertransformer.ToCreateDriverDocumentResponseDto(created)
	return &res, nil
}

func (u *driverDocumentUsecase) BulkCreate(ctx context.Context, req *dto.BulkCreateDriverDocumentsRequestDto) (*dto.BulkCreateDriverDocumentsResponseDto, error) {
	driverID, err := uuid.Parse(req.DriverID)
	if err != nil {
		return nil, err
	}

	items, err := database.WithTransaction(
		u.txManager,
		ctx,
		func(txCtx context.Context) ([]dto.CreateDriverDocumentResponseDto, error) {
			var items []dto.CreateDriverDocumentResponseDto
			for _, it := range req.Items {
				documentTypeID, err := uuid.Parse(it.DocumentTypeID)
				if err != nil {
					return nil, err
				}
				arg := pgdb.CreateDriverDocumentParams{
					DriverID:       driverID,
					DocumentTypeID: documentTypeID,
					FileUrl:        it.FileUrl,
					ExpireAt:       parseDate(it.ExpireAt),
					Status:         pgdb.DriverDocumentStatusPENDING,
				}
				created, err := u.repo.Create(txCtx, arg)
				if err != nil {
					return nil, err
				}
				items = append(items, appdrivertransformer.ToCreateDriverDocumentResponseDto(created))
			}
			return items, nil
		},
	)
	if err != nil {
		return nil, err
	}
	return &dto.BulkCreateDriverDocumentsResponseDto{Items: items}, nil
}

func (u *driverDocumentUsecase) GetByID(ctx context.Context, id uuid.UUID) (*dto.DriverDocumentItemDto, error) {
	m, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDriverDocumentNotFound
		}
		return nil, err
	}
	if m == nil {
		return nil, ErrDriverDocumentNotFound
	}
	item := appdrivertransformer.ToDriverDocumentItemDto(m)
	return &item, nil
}

func (u *driverDocumentUsecase) ListByDriverID(ctx context.Context, driverID uuid.UUID) (*dto.ListDriverDocumentsResponseDto, error) {
	list, err := u.repo.ListByDriverID(ctx, driverID)
	if err != nil {
		return nil, err
	}
	items := make([]dto.DriverDocumentItemDto, 0, len(list))
	for _, m := range list {
		items = append(items, appdrivertransformer.ToDriverDocumentItemDto(m))
	}
	return &dto.ListDriverDocumentsResponseDto{Items: items}, nil
}

func (u *driverDocumentUsecase) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateDriverDocumentRequestDto) (*dto.DriverDocumentItemDto, error) {
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDriverDocumentNotFound
		}
		return nil, err
	}
	if existing == nil {
		return nil, ErrDriverDocumentNotFound
	}
	arg := pgdb.UpdateDriverDocumentParams{
		ID:           id,
		FileUrl:      req.FileUrl,
		ExpireAt:     parseDate(req.ExpireAt),
		Status:       parseStatus(req.Status),
		RejectReason: pgtype.Text{},
		VerifiedAt:   pgtype.Timestamptz{},
		VerifiedBy:   pgtype.UUID{},
	}
	if req.RejectReason != nil {
		arg.RejectReason = pgtype.Text{String: *req.RejectReason, Valid: true}
	}

	updated, err := database.WithTransaction(
		u.txManager,
		ctx,
		func(txCtx context.Context) (*appdrivermodel.DriverDocument, error) {
			return u.repo.Update(txCtx, arg)
		},
	)
	if err != nil {
		return nil, err
	}
	item := appdrivertransformer.ToDriverDocumentItemDto(updated)
	return &item, nil
}

func (u *driverDocumentUsecase) BulkUpdate(ctx context.Context, req *dto.BulkUpdateDriverDocumentsRequestDto) (*dto.BulkUpdateDriverDocumentsResponseDto, error) {
	items, err := database.WithTransaction(
		u.txManager,
		ctx,
		func(txCtx context.Context) ([]dto.DriverDocumentItemDto, error) {
			var items []dto.DriverDocumentItemDto
			for _, it := range req.Items {
				id, err := uuid.Parse(it.ID)
				if err != nil {
					return nil, err
				}
				arg := pgdb.UpdateDriverDocumentPartialParams{ID: id}
				if it.FileUrl != nil {
					arg.FileUrl = pgtype.Text{String: *it.FileUrl, Valid: true}
				}
				if it.ExpireAt != nil {
					arg.ExpireAt = parseDate(it.ExpireAt)
				}
				if it.Status != nil {
					arg.Status = pgdb.NullDriverDocumentStatus{
						DriverDocumentStatus: parseStatus(*it.Status),
						Valid:                true,
					}
				}
				if it.RejectReason != nil {
					arg.RejectReason = pgtype.Text{String: *it.RejectReason, Valid: true}
				}
				updated, err := u.repo.UpdatePartial(txCtx, arg)
				if err != nil {
					return nil, err
				}
				items = append(items, appdrivertransformer.ToDriverDocumentItemDto(updated))
			}
			return items, nil
		},
	)
	if err != nil {
		return nil, err
	}
	return &dto.BulkUpdateDriverDocumentsResponseDto{Items: items}, nil
}

func (u *driverDocumentUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := database.WithTransaction(
		u.txManager,
		ctx,
		func(txCtx context.Context) (struct{}, error) {
			_, err := u.repo.GetByID(txCtx, id)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return struct{}{}, ErrDriverDocumentNotFound
				}
				return struct{}{}, err
			}
			return struct{}{}, u.repo.Delete(txCtx, id)
		},
	)
	return err
}
