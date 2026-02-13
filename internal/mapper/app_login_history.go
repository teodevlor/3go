package mapper

import (
	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"
)

func ToAppLoginHistory(row pgdb.AppLoginHistory) *model.AppLoginHistory {
	return &model.AppLoginHistory{
		ID:            row.ID,
		AccountID:     row.AccountID,
		DeviceID:      row.DeviceID,
		AppType:       row.AppType,
		LoginAt:       row.LoginAt.Time,
		Result:        row.Result,
		FailureReason: row.FailureReason.String,
		IpAddress:     row.IpAddress.String,
		UserAgent:     row.UserAgent.String,
		Location:      row.Location.String,
		Metadata:      row.Metadata,
		CreatedAt:     row.CreatedAt.Time,
	}
}
