package common

// ─── Driver Profile Status ──────────────────────────────────────────────────
// Khớp với DB ENUM: driver_profile_status
const (
	DriverProfileStatusPendingProfile      = "PENDING_PROFILE"
	DriverProfileStatusDocumentIncomplete  = "DOCUMENT_INCOMPLETE"
	DriverProfileStatusPendingVerification = "PENDING_VERIFICATION"
	DriverProfileStatusActive              = "ACTIVE"
	DriverProfileStatusSuspended           = "SUSPENDED"
	DriverProfileStatusRejected            = "REJECTED"

	DriverProfileStatusTextPendingProfile      = "Chưa hoàn tất hồ sơ"
	DriverProfileStatusTextDocumentIncomplete  = "Thiếu tài liệu"
	DriverProfileStatusTextPendingVerification = "Chờ xác minh"
	DriverProfileStatusTextActive              = "Kích hoạt"
	DriverProfileStatusTextSuspended           = "Tạm khóa"
	DriverProfileStatusTextRejected            = "Từ chối"
)

// ─── Driver Document Status ──────────────────────────────────────────────────
// Khớp với DB ENUM: driver_document_status
const (
	DriverDocumentStatusPending  = "PENDING"
	DriverDocumentStatusApproved = "APPROVED"
	DriverDocumentStatusRejected = "REJECTED"
)

// ─── Driver Service Status ───────────────────────────────────────────────────
// Khớp với DB ENUM: driver_service_status
const (
	DriverServiceStatusPendingDocument = "PENDING_DOCUMENT"
	DriverServiceStatusPendingApproval = "PENDING_APPROVAL"
	DriverServiceStatusActive          = "ACTIVE"
	DriverServiceStatusSuspended       = "SUSPENDED"
	DriverServiceStatusRejected        = "REJECTED"
)

// ─── OTP Status ──────────────────────────────────────────────────────────────
// Giá trị cột status trong bảng system_otps
const (
	OTPStatusActive  = "active"
	OTPStatusUsed    = "used"
	OTPStatusExpired = "expired"
	OTPStatusLocked  = "locked"
)

// ─── OTP Purpose ─────────────────────────────────────────────────────────────
const (
	OTPPurposeUserRegister   = "user_register"
	OTPPurposeDriverRegister = "driver_register"
	OTPPurposeResetPassword  = "reset_password"
)

// ─── OTP Audit Result ────────────────────────────────────────────────────────
const (
	OTPResultSuccess = "success"
	OTPResultFailed  = "failed"
	OTPResultLocked  = "locked"
)

// ─── OTP Failure Reason ──────────────────────────────────────────────────────
const (
	OTPFailureReasonInvalidCode = "invalid_code"
	OTPFailureReasonMaxAttempt  = "max_attempt"
)

// ─── Gender ───────────────────────────────────────────────────────────────────
const (
	GenderMale   = "MALE"
	GenderFemale = "FEMALE"
	GenderOther  = "OTHER"
)

// ─── App Type ────────────────────────────────────────────────────────────────
// Dùng trong account_app_devices.app_type và app_login_histories.app_type
const (
	AppTypeUser   = "user"
	AppTypeDriver = "driver"
)

// ─── Setting Type ─────────────────────────────────────────────────────────────
// Khớp với DB ENUM: setting_type
const (
	SettingTypeWebSystem = "web_system"
	SettingTypeAppUser   = "app_user"
	SettingTypeAppDriver = "app_driver"
)

// ─── Admin Department ─────────────────────────────────────────────────────────
// Khớp với DB ENUM: admin_department
const (
	AdminDepartmentEmployee = "employee"
	AdminDepartmentAdmin    = "admin"
	AdminDepartmentSeller   = "seller"
	AdminDepartmentMarketer = "marketer"
)

// ─── Surcharge Unit ───────────────────────────────────────────────────────────
const (
	SurchargeUnitPercent = "percent"
	SurchargeUnitFixed   = "fixed"
)

// ─── Surcharge Condition Type ─────────────────────────────────────────────────
const (
	ConditionTypeTimeWindow = "time_window"
	ConditionTypeWeather    = "weather"
	ConditionTypeTraffic    = "traffic"
	ConditionTypeHoliday    = "holiday"
)

// ─── Traffic Level ────────────────────────────────────────────────────────────
const (
	TrafficLevelLow    = "low"
	TrafficLevelMedium = "medium"
	TrafficLevelHigh   = "high"
)
