package validator

var defaultMessages = map[string]string{
	"required": "Required field",
	"email":    "Email is not valid",
	"min":      "The field must be at least %s characters",
	"max":      "The field must be at most %s characters",
	"password": "Password must be at least 8 characters, 1 uppercase letter, 1 lowercase letter, 1 number, 1 special character",
}
