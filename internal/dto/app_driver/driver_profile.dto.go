package app_driver

import (
	"time"

	"go-structure/internal/common"
	dto_common "go-structure/internal/dto/common"

	"github.com/google/uuid"
)

const (
	GenderMale   = common.GenderMale
	GenderFemale = common.GenderFemale
	GenderOther  = common.GenderOther
)

type (
	DriverRegisterRequestDto struct {
		Phone      string      `json:"phone" binding:"required"`
		FullName   string      `json:"full_name" binding:"required"`
		Password   string      `json:"password" binding:"required"`
		ServiceIDs []uuid.UUID `json:"service_ids" binding:"omitempty,dive,required"`
	}

	AdminCreateDriverProfileRequestDto struct {
		Phone       string      `json:"phone" binding:"required"`
		Password    string      `json:"password" binding:"required"`
		FullName    string      `json:"full_name" binding:"required"`
		DateOfBirth *string     `json:"date_of_birth"`
		Gender      string      `json:"gender"`
		Address     string      `json:"address"`
		ServiceIDs  []uuid.UUID `json:"service_ids" binding:"omitempty,dive,required"`
	}

	DriverRegisterResponseDto struct {
		UserMessage string `json:"user_message" example:"Đăng ký tài khoản tài xế thành công, vui lòng kiểm tra điện thoại để nhận mã OTP"`
	}

	DriverVerifyOtpRequestDto struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}

	DriverVerifyOtpResponseDto struct {
		UserMessage string    `json:"user_message" example:"Xác thực OTP thành công. Bạn có thể đăng nhập và upload tài liệu."`
		DriverID    uuid.UUID `json:"driver_id"`
	}

	DriverDeviceDto struct {
		DeviceUID  string `json:"device_uid" binding:"required"`
		Platform   string `json:"platform" binding:"required"`
		DeviceName string `json:"device_name" binding:"required"`
		OsVersion  string `json:"os_version" binding:"required"`
		AppVersion string `json:"app_version" binding:"required"`
		FCMToken   string `json:"fcm_token"`
	}

	DriverLoginRequestDto struct {
		Phone    string          `json:"phone" binding:"required"`
		Password string          `json:"password" binding:"required"`
		Device   DriverDeviceDto `json:"device" binding:"required"`
	}

	DriverProfileItemDto struct {
		ID                   uuid.UUID   `json:"id"`
		AccountID            uuid.UUID   `json:"account_id"`
		Phone                string      `json:"phone"`
		FullName             string      `json:"full_name"`
		DateOfBirth          *time.Time  `json:"date_of_birth,omitempty"`
		Gender               string      `json:"gender,omitempty"`
		Address              string      `json:"address,omitempty"`
		GlobalStatus         string      `json:"global_status"`
		GlobalStatusText     string      `json:"global_status_text"`
		Rating               float64     `json:"rating"`
		TotalCompletedOrders int32       `json:"total_completed_orders"`
		ServiceIDs           []uuid.UUID `json:"service_ids,omitempty"`
		CreatedAt            time.Time   `json:"created_at"`
		UpdatedAt            time.Time   `json:"updated_at"`
	}

	DriverLoginResponseDto struct {
		RequireVerifyOtp bool                  `json:"require_verify_otp"`
		Message          string                `json:"message,omitempty"`
		AccessToken      string                `json:"access_token,omitempty"`
		RefreshToken     string                `json:"refresh_token,omitempty"`
		DriverProfile    *DriverProfileItemDto `json:"driver_profile,omitempty"`
	}

	DriverLocationStatusRequestDto struct {
		Lat    float64 `json:"lat" binding:"required"`
		Lng    float64 `json:"lng" binding:"required"`
		Status string  `json:"status" binding:"required"` // ví dụ: "idle", "busy"
	}

	UpdateDriverProfileRequestDto struct {
		FullName     string  `json:"full_name" binding:"required"`
		DateOfBirth  *string `json:"date_of_birth"`
		Gender       string  `json:"gender"`
		Address      string  `json:"address"`
		GlobalStatus string  `json:"global_status" binding:"omitempty,oneof=PENDING_PROFILE DOCUMENT_INCOMPLETE PENDING_VERIFICATION ACTIVE SUSPENDED REJECTED"`
	}

	ListDriverProfilesResponseDto struct {
		Items      []DriverProfileItemDto    `json:"items"`
		Pagination dto_common.PaginationMeta `json:"pagination"`
	}
)
