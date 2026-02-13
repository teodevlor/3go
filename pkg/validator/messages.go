package validator

var defaultMessages = map[string]string{
	"required": "Trường bắt buộc",
	"email":    "Email không hợp lệ",
	"min":      "Trường phải có ít nhất %s ký tự",
	"max":      "Trường tối đa %s ký tự",
	"uuid":     "Giá trị phải là UUID hợp lệ",
	"gte":      "Giá trị phải lớn hơn hoặc bằng %s",
	"password": "Mật khẩu tối thiểu 8 ký tự, 1 chữ hoa, 1 chữ thường, 1 số, 1 ký tự đặc biệt",
}
