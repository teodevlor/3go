package appdriver

import (
	"time"

	pgdb "go-structure/internal/orm/db/postgres"
	appdrivermodel "go-structure/internal/repository/model/app_driver"
)

func ToDriverDocumentTypeFromRow(row *pgdb.DriverDocumentType) *appdrivermodel.DriverDocumentType {
	if row == nil {
		return nil
	}
	var deletedAt *time.Time
	if row.DeletedAt.Valid {
		t := row.DeletedAt.Time
		deletedAt = &t
	}
	return &appdrivermodel.DriverDocumentType{
		ID:                row.ID,
		Code:              row.Code,
		Name:              row.Name,
		Description:       row.Description,
		IsRequired:        row.IsRequired,
		RequireExpireDate: row.RequireExpireDate,
		ServiceID:         row.ServiceID,
		IsActive:          row.IsActive,
		CreatedAt:         row.CreatedAt.Time,
		DeletedAt:         deletedAt,
	}
}
