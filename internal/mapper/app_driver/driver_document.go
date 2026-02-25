package appdriver

import (
	common "go-structure/internal/common"
	pgdb "go-structure/internal/orm/db/postgres"
	appdrivermodel "go-structure/internal/repository/model/app_driver"
)

func ToDriverDocumentFromRow(row *pgdb.DriverDocument) *appdrivermodel.DriverDocument {
	if row == nil {
		return nil
	}

	out := &appdrivermodel.DriverDocument{
		ID:             row.ID,
		DriverID:       row.DriverID,
		DocumentTypeID: row.DocumentTypeID,
		FileUrl:        row.FileUrl,
		Status:         string(row.Status),
		RejectReason:   common.NullableTextToStringPtr(row.RejectReason),
		VerifiedBy:     common.NullableUUIDToUUIDPtr(row.VerifiedBy),
		CreatedAt:      row.CreatedAt.Time,
		UpdatedAt:      row.UpdatedAt.Time,
		ExpireAt:       common.NullableDateToTimePtr(row.ExpireAt),
		VerifiedAt:     common.NullableTimestamptzToTimePtr(row.VerifiedAt),
		DeletedAt:      common.NullableTimestamptzToTimePtr(row.DeletedAt),
	}

	return out
}
