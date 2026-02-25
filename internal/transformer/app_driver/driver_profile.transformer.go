package app_driver

import (
	dto "go-structure/internal/dto/app_driver"
	"go-structure/internal/repository/model"
	appdrivermodel "go-structure/internal/repository/model/app_driver"
)

func ToDriverProfileItemDto(account *model.Account, profile *appdrivermodel.DriverProfile) dto.DriverProfileItemDto {
	if profile == nil {
		return dto.DriverProfileItemDto{}
	}
	out := dto.DriverProfileItemDto{
		ID:                   profile.ID,
		AccountID:            profile.AccountID,
		FullName:             profile.FullName,
		DateOfBirth:          profile.DateOfBirth,
		Gender:               profile.Gender,
		Address:              profile.Address,
		GlobalStatus:         profile.GlobalStatus,
		Rating:               profile.Rating,
		TotalCompletedOrders: profile.TotalCompletedOrders,
		CreatedAt:            profile.CreatedAt,
		UpdatedAt:            profile.UpdatedAt,
	}
	if account != nil {
		out.Phone = account.Phone
	}
	return out
}
