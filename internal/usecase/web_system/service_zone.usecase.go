package web_system

import (
	"context"

	"go-structure/internal/helper/database"
	websystem_repo "go-structure/internal/repository/web_system"

	"github.com/google/uuid"
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
	if u.serviceZoneRepo == nil {
		return nil
	}
	if _, inTx := database.TransactionFromContext(ctx); inTx {
		return u.setZonesForServiceInTx(ctx, serviceID, zoneIDs)
	}
	return u.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		return u.setZonesForServiceInTx(txCtx, serviceID, zoneIDs)
	})
}

func (u *serviceZoneUsecase) GetZoneIDsByServiceID(ctx context.Context, serviceID uuid.UUID) ([]uuid.UUID, error) {
	if u.serviceZoneRepo == nil {
		return nil, nil
	}
	return u.serviceZoneRepo.ListServiceZoneIDsByServiceID(ctx, serviceID)
}

