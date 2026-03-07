package web_system

import (
	"context"

	"go-structure/global"
	"go-structure/internal/helper/database"
	"go-structure/internal/middleware"
	websystem_repo "go-structure/internal/repository/web_system"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	IServiceZoneUsecase interface {
		SetZonesForService(ctx context.Context, serviceID uuid.UUID, zoneIDs []uuid.UUID) error
		GetZoneIDsByServiceID(ctx context.Context, serviceID uuid.UUID) ([]uuid.UUID, error)
	}

	serviceZoneUsecase struct {
		serviceZoneRepo    websystem_repo.IServiceZoneRepository
		transactionManager database.TransactionManager
	}
)

func NewServiceZoneUsecase(serviceZoneRepo websystem_repo.IServiceZoneRepository, transactionManager database.TransactionManager) IServiceZoneUsecase {
	return &serviceZoneUsecase{
		serviceZoneRepo:    serviceZoneRepo,
		transactionManager: transactionManager,
	}
}

func (u *serviceZoneUsecase) setZonesForServiceInTx(ctx context.Context, serviceID uuid.UUID, zoneIDs []uuid.UUID) error {
	if err := u.serviceZoneRepo.DeleteServiceZonesByServiceID(ctx, serviceID); err != nil {
		return err
	}
	if len(zoneIDs) == 0 {
		return nil
	}
	return u.serviceZoneRepo.CreateServiceZones(ctx, serviceID, zoneIDs)
}

func (u *serviceZoneUsecase) SetZonesForService(ctx context.Context, serviceID uuid.UUID, zoneIDs []uuid.UUID) error {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("SetZonesForService: start", zap.String(global.KeyCorrelationID, cid), zap.String("service_id", serviceID.String()), zap.Int("zone_count", len(zoneIDs)))

	if u.serviceZoneRepo == nil {
		return nil
	}
	if _, inTx := database.TransactionFromContext(ctx); inTx {
		return u.setZonesForServiceInTx(ctx, serviceID, zoneIDs)
	}
	_, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (struct{}, error) {
			return struct{}{}, u.setZonesForServiceInTx(txCtx, serviceID, zoneIDs)
		},
	)
	if err != nil {
		global.Logger.Error("SetZonesForService: failed to set zones", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return err
	}
	global.Logger.Info("SetZonesForService: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.String("service_id", serviceID.String()))
	return nil
}

func (u *serviceZoneUsecase) GetZoneIDsByServiceID(ctx context.Context, serviceID uuid.UUID) ([]uuid.UUID, error) {
	cid := middleware.CorrelationIDFromContext(ctx)
	global.Logger.Info("GetZoneIDsByServiceID: start", zap.String(global.KeyCorrelationID, cid), zap.String("service_id", serviceID.String()))

	if u.serviceZoneRepo == nil {
		return nil, nil
	}
	ids, err := u.serviceZoneRepo.ListServiceZoneIDsByServiceID(ctx, serviceID)
	if err != nil {
		global.Logger.Error("GetZoneIDsByServiceID: failed to list zone IDs", zap.String(global.KeyCorrelationID, cid), zap.Error(err))
		return nil, err
	}
	global.Logger.Info("GetZoneIDsByServiceID: completed successfully", zap.String(global.KeyCorrelationID, cid), zap.Int("count", len(ids)))
	return ids, nil
}
