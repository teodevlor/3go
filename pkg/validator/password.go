package validator

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

const PasswordMinLength = 8

var (
	reUppercase = regexp.MustCompile(`[A-Z]`)
	reLowercase = regexp.MustCompile(`[a-z]`)
	reDigit     = regexp.MustCompile(`[0-9]`)
	reSpecial   = regexp.MustCompile(`[^A-Za-z0-9]`)
)

func HashPassword(pw string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CheckPassword(plain, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

func ValidatePassword(pw string) error {
	if len(pw) < PasswordMinLength {
		return fmt.Errorf("mật khẩu phải có ít nhất %d ký tự", PasswordMinLength)
	}
	if !reUppercase.MatchString(pw) {
		return fmt.Errorf("mật khẩu phải có ít nhất 1 chữ in hoa")
	}
	if !reLowercase.MatchString(pw) {
		return fmt.Errorf("mật khẩu phải có ít nhất 1 chữ thường")
	}
	if !reDigit.MatchString(pw) {
		return fmt.Errorf("mật khẩu phải có ít nhất 1 chữ số")
	}
	if !reSpecial.MatchString(pw) {
		return fmt.Errorf("mật khẩu phải có ít nhất 1 ký tự đặc biệt")
	}
	return nil
}
