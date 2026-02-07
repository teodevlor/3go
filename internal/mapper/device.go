package mapper

import (
	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"
)

func ToDevice(row pgdb.Device) *model.Device {
	return &model.Device{
		ID:         row.ID,
		DeviceUID:  row.DeviceUid,
		Platform:   row.Platform,
		DeviceName: row.DeviceName.String,
		OsVersion:  row.OsVersion.String,
		AppVersion: row.AppVersion.String,
		Metadata:   row.Metadata,
		BaseModel: model.BaseModel{
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
			DeletedAt: row.DeletedAt.Time,
		},
	}
}
