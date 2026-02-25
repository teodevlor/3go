package appdriver

import (
	common "go-structure/internal/common"
	pgdb "go-structure/internal/orm/db/postgres"
	appdrivermodel "go-structure/internal/repository/model/app_driver"
)

func ToDriverProfileFromRow(row *pgdb.DriverProfile) *appdrivermodel.DriverProfile {
	if row == nil {
		return nil
	}

	return &appdrivermodel.DriverProfile{
		ID:                   row.ID,
		AccountID:            row.AccountID,
		FullName:             row.FullName,
		DateOfBirth:          common.NullableDateToTimePtr(row.DateOfBirth),
		Gender:               common.NullableTextToString(row.Gender),
		Address:              common.NullableTextToString(row.Address),
		GlobalStatus:         string(row.GlobalStatus),
		Rating:               common.NullableNumericToFloat64(row.Rating),
		TotalCompletedOrders: common.NullableInt4ToInt32(row.TotalCompletedOrders),
		CreatedAt:            row.CreatedAt.Time,
		UpdatedAt:            row.UpdatedAt.Time,
		DeletedAt:            common.NullableTimestamptzToTimePtr(row.DeletedAt),
	}
}
