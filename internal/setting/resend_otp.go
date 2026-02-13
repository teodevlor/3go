package setting

type ResendOTPConfig struct {
	MaxCount       int32              `json:"maxCount"`       // Số lần tối đa được gửi lại OTP trong khoảng trackingTTL
	TimeOutExpired int32              `json:"timeOutExpired"` // Thời gian hết hạn OTP (giây)
	TimeOutResent  int32              `json:"timeOutResent"`  // Cooldown giữa hai lần gửi lại (giây)
	BlockDurations ResendOTPBlockDurs `json:"blockDurations"` // Thời gian block (giây) theo số lần vi phạm
	TrackingTTL    int64              `json:"trackingTTL"`    // Cửa sổ thời gian đếm số lần resend (giây)
}

// ResendOTPBlockDurs thời gian block theo mức vi phạm (giây).
type ResendOTPBlockDurs struct {
	Violation1     int32 `json:"violation1"`
	Violation2     int32 `json:"violation2"`
	Violation3     int32 `json:"violation3"`
	Violation4Plus int32 `json:"violation4Plus"`
}
