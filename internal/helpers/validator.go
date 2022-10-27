package helpers

import (
	"unicode"

	er "github.com/Lerner17/gophermart/internal/errors"
)

func ValidatePassword(password string) error {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	if len(password) >= 7 {
		hasMinLen = true
	}
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	if hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial {
		return nil
	}
	return er.ErrInvalidPasswordPattern
}
