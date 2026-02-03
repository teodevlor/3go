package common

import "time"

type ResponseData struct {
	Status    bool     `json:"status"`
	Code      int      `json:"code"`
	Message   string   `json:"message"`
	Errors    []string `json:"errors,omitempty"`
	Data      any      `json:"data,omitempty"`
	Timestamp string   `json:"timestamp"`
}

func SuccessResponse(code int, data any, message string) *ResponseData {
	if message == "" {
		message = HTTPMessage(code)
	}
	return &ResponseData{
		Status:    true,
		Code:      code,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func ErrorResponse(code int, errors []string, message string) *ResponseData {
	if message == "" {
		message = HTTPMessage(code)
	}
	if len(errors) == 0 && message != "" {
		errors = []string{message}
	}
	return &ResponseData{
		Status:    false,
		Code:      code,
		Message:   message,
		Errors:    errors,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
