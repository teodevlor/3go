package mapper

import (
	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"
)

func ToLoginHistory(row pgdb.LoginHistory) *model.LoginHistory {
	return &model.LoginHistory{
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
