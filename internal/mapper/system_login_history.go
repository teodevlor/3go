package mapper

import (
	pgdb "go-structure/internal/orm/db/postgres"
	"go-structure/internal/repository/model"
)

func ToSystemLoginHistory(row pgdb.SystemLoginHistory) *model.SystemLoginHistory {
	return &model.SystemLoginHistory{
		ID:            row.ID,
		AdminID:       row.AdminID,
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
