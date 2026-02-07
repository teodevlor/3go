package appuser

import "github.com/google/uuid"

type (
	UserRegisterRequestDto struct {
		Phone    string `json:"phone" binding:"required"`
		FullName string `json:"full_name" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	UserRegisterResponseDto struct {
		UserMessage string `json:"user_message" example:"Đăng ký tài khoản thành công, vui lòng kiểm tra điện thoại để nhận mã OTP"`
	}

	UserActiveRequestDto struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}

	UserActiveResponseDto struct {
		UserMessage string `json:"user_message" example:"Kích hoạt tài khoản thành công"`
	}

	UserLoginRequestDto struct {
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" binding:"required"`
		Device   Device `json:"device" binding:"required"`
	}

	Device struct {
		DeviceUID  string `json:"device_uid" binding:"required"`
		Platform   string `json:"platform" binding:"required"`
		DeviceName string `json:"device_name" binding:"required"`
		OsVersion  string `json:"os_version" binding:"required"`
		AppVersion string `json:"app_version" binding:"required"`
		FCMToken   string `json:"fcm_token"`
	}

	UserLoginResponseDto struct {
		AccessToken  string                 `json:"access_token" example:"access_token"`
		RefreshToken string                 `json:"refresh_token" example:"refresh_token"`
		UserProfile  UserProfileResponseDto `json:"user_profile"`
	}

	UserProfileResponseDto struct {
		ID        uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
		FullName  string    `json:"full_name" example:"Nguyễn Văn A"`
		AvatarURL string    `json:"avatar_url" example:"https://example.com/avatar.jpg"`
		IsActive  bool      `json:"is_active" example:"true"`
		Phone     string    `json:"phone" example:"+84123456789"`
		Email     string    `json:"email" example:"user@example.com"`
	}

	RefreshTokenRequestDto struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	RefreshTokenResponseDto struct {
		AccessToken  string `json:"access_token" example:"new_access_token"`
		RefreshToken string `json:"refresh_token" example:"new_refresh_token"`
	}

	LogoutResponseDto struct {
		UserMessage string `json:"user_message" example:"Đăng xuất thành công"`
	}

	UpdateUserProfile struct {
		FullName  string `json:"full_name" binding:"required"`
		AvatarURL string `json:"avatar_url" binding:"required"`
	}

	UpdateUserProfileResponseDto struct {
		UserMessage string                 `json:"user_message" example:"Cập nhật thông tin thành công"`
		UserProfile UserProfileResponseDto `json:"user_profile"`
	}

	ForgotPasswordRequestDto struct {
		Phone string `json:"phone" binding:"required"`
	}

	ForgotPasswordResponseDto struct {
		UserMessage string `json:"user_message" example:"Vui lòng kiểm tra điện thoại để nhận mã OTP đặt lại mật khẩu"`
	}

	ResetPasswordRequestDto struct {
		Phone           string `json:"phone" binding:"required"`
		Code            string `json:"code" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
	}

	ResetPasswordResponseDto struct {
		UserMessage string `json:"user_message" example:"Đặt lại mật khẩu thành công"`
	}
)
