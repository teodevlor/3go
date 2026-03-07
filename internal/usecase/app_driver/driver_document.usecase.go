package app_driver

import (
	"context"
	"errors"

	"go-structure/global"
	common "go-structure/internal/common"
	dto "go-structure/internal/dto/app_driver"
	"go-structure/internal/helper/database"
	"go-structure/internal/middleware"
	appdriverrepo "go-structure/internal/repository/app_driver"
	appdrivermodel "go-structure/internal/repository/model/app_driver"
	appdrivertransformer "go-structure/internal/transformer/app_driver"
	pgdb "go-structure/orm/db/postgres"
	"go-structure/pkg/validator"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

var (
	ErrDriverDocumentNotFound = errors.New("không tìm thấy giấy tờ tài xế")
	ErrDriverProfileNotFound  = errors.New("không tìm thấy hồ sơ tài xế")
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
		repo             appdriverrepo.IDriverDocumentRepository
		driverProfileRepo appdriverrepo.IDriverProfileRepository
		documentTypeRepo  appdriverrepo.IDriverDocumentTypeRepository
		txManager        database.TransactionManager
	}
)

func NewDriverDocumentUsecase(
	repo appdriverrepo.IDriverDocumentRepository,
	driverProfileRepo appdriverrepo.IDriverProfileRepository,
	documentTypeRepo appdriverrepo.IDriverDocumentTypeRepository,
	txManager database.TransactionManager,
) IDriverDocumentUsecase {
	return &driverDocumentUsecase{
		repo:              repo,
		driverProfileRepo: driverProfileRepo,
		documentTypeRepo:  documentTypeRepo,
		txManager:         txManager,
	}
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
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("Create: start", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", req.DriverID), zap.String("document_type_id", req.DocumentTypeID))

	driverID, err := uuid.Parse(req.DriverID)
	if err != nil {
		global.Logger.Error("Create: failed to parse driver_id", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	documentTypeID, err := uuid.Parse(req.DocumentTypeID)
	if err != nil {
		global.Logger.Error("Create: failed to parse document_type_id", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}

	if u.driverProfileRepo != nil {
		if err := validator.CheckExists(ctx, driverID, func(ctx context.Context, id uuid.UUID) error {
			_, err := u.driverProfileRepo.GetByID(ctx, id)
			return err
		}, ErrDriverProfileNotFound); err != nil {
			global.Logger.Error("Create: failed to validate driver profile", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
			return nil, err
		}
	}
	if u.documentTypeRepo != nil {
		if err := validator.CheckExists(ctx, documentTypeID, func(ctx context.Context, id uuid.UUID) error {
			_, err := u.documentTypeRepo.GetByID(ctx, id)
			return err
		}, ErrDriverDocumentTypeNotFound); err != nil {
			global.Logger.Error("Create: failed to validate document type", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
			return nil, err
		}
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
		global.Logger.Error("Create: transaction failed", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("Create: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("document_id", created.ID.String()))
	res := appdrivertransformer.ToCreateDriverDocumentResponseDto(created)
	return &res, nil
}

func (u *driverDocumentUsecase) BulkCreate(ctx context.Context, req *dto.BulkCreateDriverDocumentsRequestDto) (*dto.BulkCreateDriverDocumentsResponseDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("BulkCreate: start", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", req.DriverID), zap.Int("item_count", len(req.Items)))

	driverID, err := uuid.Parse(req.DriverID)
	if err != nil {
		global.Logger.Error("BulkCreate: failed to parse driver_id", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}

	if u.driverProfileRepo != nil {
		if err := validator.CheckExists(ctx, driverID, func(ctx context.Context, id uuid.UUID) error {
			_, err := u.driverProfileRepo.GetByID(ctx, id)
			return err
		}, ErrDriverProfileNotFound); err != nil {
			global.Logger.Error("BulkCreate: failed to validate driver profile", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
			return nil, err
		}
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
				if u.documentTypeRepo != nil {
					if err := validator.CheckExists(txCtx, documentTypeID, func(ctx context.Context, id uuid.UUID) error {
						_, err := u.documentTypeRepo.GetByID(ctx, id)
						return err
					}, ErrDriverDocumentTypeNotFound); err != nil {
						return nil, err
					}
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
		global.Logger.Error("BulkCreate: transaction failed", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("BulkCreate: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.Int("item_count", len(items)))
	return &dto.BulkCreateDriverDocumentsResponseDto{Items: items}, nil
}

func (u *driverDocumentUsecase) GetByID(ctx context.Context, id uuid.UUID) (*dto.DriverDocumentItemDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("GetByID: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

	m, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			global.Logger.Error("GetByID: document not found", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
			return nil, ErrDriverDocumentNotFound
		}
		global.Logger.Error("GetByID: failed to get document", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	if m == nil {
		global.Logger.Error("GetByID: document not found", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
		return nil, ErrDriverDocumentNotFound
	}
	global.Logger.Info("GetByID: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	item := appdrivertransformer.ToDriverDocumentItemDto(m)
	return &item, nil
}

func (u *driverDocumentUsecase) ListByDriverID(ctx context.Context, driverID uuid.UUID) (*dto.ListDriverDocumentsResponseDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("ListByDriverID: start", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", driverID.String()))

	list, err := u.repo.ListByDriverID(ctx, driverID)
	if err != nil {
		global.Logger.Error("ListByDriverID: failed to list documents", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	items := make([]dto.DriverDocumentItemDto, 0, len(list))
	typeCache := make(map[uuid.UUID]dto.DriverDocumentTypeItemDto)
	for _, m := range list {
		item := appdrivertransformer.ToDriverDocumentItemDto(m)
		if u.documentTypeRepo != nil {
			var dtDto dto.DriverDocumentTypeItemDto
			if cached, ok := typeCache[m.DocumentTypeID]; ok {
				dtDto = cached
			} else {
				docType, err := u.documentTypeRepo.GetByID(ctx, m.DocumentTypeID)
				if err == nil && docType != nil {
					dtDto = appdrivertransformer.ToDriverDocumentTypeItemDto(docType)
					typeCache[m.DocumentTypeID] = dtDto
				}
			}
			if dtDto.ID != "" {
				ptr := new(dto.DriverDocumentTypeItemDto)
				*ptr = dtDto
				item.DocumentType = ptr
			}
		}
		items = append(items, item)
	}
	global.Logger.Info("ListByDriverID: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("driver_id", driverID.String()), zap.Int("count", len(items)))
	return &dto.ListDriverDocumentsResponseDto{Items: items}, nil
}

func (u *driverDocumentUsecase) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateDriverDocumentRequestDto) (*dto.DriverDocumentItemDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("Update: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			global.Logger.Error("Update: document not found", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
			return nil, ErrDriverDocumentNotFound
		}
		global.Logger.Error("Update: failed to get document", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	if existing == nil {
		global.Logger.Error("Update: document not found", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
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
		global.Logger.Error("Update: transaction failed", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("Update: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	item := appdrivertransformer.ToDriverDocumentItemDto(updated)
	return &item, nil
}

func (u *driverDocumentUsecase) BulkUpdate(ctx context.Context, req *dto.BulkUpdateDriverDocumentsRequestDto) (*dto.BulkUpdateDriverDocumentsResponseDto, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("BulkUpdate: start", zap.String(global.KeyCorrelationID, cid), zap.Int("item_count", len(req.Items)))

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
		global.Logger.Error("BulkUpdate: transaction failed", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("BulkUpdate: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.Int("item_count", len(items)))
	return &dto.BulkUpdateDriverDocumentsResponseDto{Items: items}, nil
}

func (u *driverDocumentUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("Delete: start", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))

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
	if err != nil {
		if errors.Is(err, ErrDriverDocumentNotFound) {
			global.Logger.Error("Delete: document not found", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
		} else {
			global.Logger.Error("Delete: failed to delete document", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		}
		return err
	}
	global.Logger.Info("Delete: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("id", id.String()))
	return nil
}
