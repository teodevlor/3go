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

func resolveMessage(code int, message ...string) string {
	if len(message) > 0 && message[0] != "" {
		return message[0]
	}
	return HTTPMessage(code)
}

func SuccessResponse(code int, data any, message ...string) *ResponseData {
	msg := resolveMessage(code, message...)
	return &ResponseData{
		Status:    true,
		Code:      code,
		Message:   msg,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func ErrorResponse(code int, errors []string, message ...string) *ResponseData {
	msg := resolveMessage(code, message...)
	return &ResponseData{
		Status:    false,
		Code:      code,
		Message:   msg,
		Errors:    errors,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
