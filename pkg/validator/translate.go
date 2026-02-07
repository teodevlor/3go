package validator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-playground/validator/v10"
)

func Translate(err error) []string {
	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return []string{"Invalid data"}
	}

	errMap := make(map[string]string)

	for _, fe := range ve {
		field := strings.ToLower(fe.Field())
		tag := fe.Tag()

		msg, exists := defaultMessages[tag]
		if !exists {
			msg = "Invalid value"
		}

		switch tag {
		case "min", "max":
			msg = fmt.Sprintf(msg, fe.Param())
		}

		// mỗi field giữ 1 lỗi cuối cùng
		errMap[field] = msg
	}

	if len(errMap) == 0 {
		return []string{"Invalid data"}
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
