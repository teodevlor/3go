package common

import (
	"errors"
	"go-structure/internal/constants"
)

var (
	ErrResendTooSoon     = errors.New(constants.BaseMessageResendOTPTooSoon)
	ErrResendMaxExceeded = errors.New(constants.BaseMessageResendOTPMaxExceeded)
)

type ErrorWithRetryAfter struct {
	Err               error
	RetryAfterSeconds int64
}

func (e *ErrorWithRetryAfter) Error() string { return e.Err.Error() }
func (e *ErrorWithRetryAfter) Unwrap() error { return e.Err }
