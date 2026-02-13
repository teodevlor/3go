package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/go-playground/validator/v10"
)

func Translate(err error) []string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		return translateValidationErrors(ve)
	}

	var ute *json.UnmarshalTypeError
	if errors.As(err, &ute) {
		field := strings.ToLower(ute.Field)
		if field == "" {
			field = "body"
		}
		return []string{fmt.Sprintf("%s: kiểu dữ liệu không hợp lệ", field)}
	}

	msg := err.Error()
	if strings.Contains(strings.ToLower(msg), "unmarshal") || strings.Contains(strings.ToLower(msg), "invalid character") {
		return []string{"Kiểu dữ liệu JSON không hợp lệ (kiểm tra định dạng từng trường)"}
	}

	return []string{"Dữ liệu không hợp lệ"}
}

func translateValidationErrors(ve validator.ValidationErrors) []string {

	errMap := make(map[string]string)

	for _, fe := range ve {
		field := strings.ToLower(fe.Field())
		tag := fe.Tag()

		msg, exists := defaultMessages[tag]
		if !exists {
			msg = "Invalid value"
		}

		switch tag {
		case "min", "max", "gte":
			msg = fmt.Sprintf(msg, fe.Param())
		}

		errMap[field] = msg
	}

	if len(errMap) == 0 {
		return []string{"Dữ liệu không hợp lệ"}
	}

	keys := make([]string, 0, len(errMap))
	for k := range errMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	messages := make([]string, 0, len(keys))
	for _, k := range keys {
		messages = append(messages, fmt.Sprintf("%s: %s", k, errMap[k]))
	}

	return messages
}
