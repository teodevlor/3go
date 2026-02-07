package common

import "errors"

// Sentinel errors dùng cho resend OTP (usecase trả về, controller so sánh errors.Is).
var (
	ErrResendTooSoon     = errors.New(BaseMessageResendOTPTooSoon)
	ErrResendMaxExceeded = errors.New(BaseMessageResendOTPMaxExceeded)
)

type ErrorWithRetryAfter struct {
	Err               error
	RetryAfterSeconds int64
}

func (e *ErrorWithRetryAfter) Error() string { return e.Err.Error() }
func (e *ErrorWithRetryAfter) Unwrap() error { return e.Err }
