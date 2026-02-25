package app_driver

import (
	"time"

	"github.com/google/uuid"
)

type (
	DriverRegisterRequestDto struct {
		Phone    string `json:"phone" binding:"required"`
		FullName string `json:"full_name" binding:"required"`
		Password string `json:"password" binding:"required"`
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
		ID                   uuid.UUID  `json:"id"`
		AccountID            uuid.UUID  `json:"account_id"`
		Phone                string     `json:"phone"`
		FullName             string     `json:"full_name"`
		DateOfBirth          *time.Time `json:"date_of_birth,omitempty"`
		Gender               string     `json:"gender,omitempty"`
		Address              string     `json:"address,omitempty"`
		GlobalStatus         string     `json:"global_status"`
		Rating               float64    `json:"rating"`
		TotalCompletedOrders int32      `json:"total_completed_orders"`
		CreatedAt            time.Time  `json:"created_at"`
		UpdatedAt            time.Time  `json:"updated_at"`
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
)
