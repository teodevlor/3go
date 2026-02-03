package validator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-playground/validator/v10"
)

func Translate(err error) map[string]string {
	result := make(map[string]string)

	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return result
	}

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

		result[field] = msg
	}

	return result
}

func TranslateMessage(err error) string {
	errs := Translate(err)
	if len(errs) == 0 {
		return "Invalid data"
	}

	keys := make([]string, 0, len(errs))
	for k := range errs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s: %s", k, errs[k]))
	}
	return strings.Join(parts, "; ")
}
