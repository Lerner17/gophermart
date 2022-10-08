package models

import (
	"encoding/json"
	"unicode"

	er "github.com/Lerner17/gophermart/internal/errors"
	"github.com/Lerner17/gophermart/internal/helpers"
)

type RegisterUser struct {
	Login          string `json:"login"`
	Password       string `json:"password"`
	HashedPassword string
}

func (ru *RegisterUser) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{
		Login:    ru.Login,
		Password: ru.HashedPassword,
	})
}

func (ru *RegisterUser) HashPassword() {
	passwd, err := helpers.HashAndSalt([]byte(ru.Password))
	if err != nil {
		panic(err)
	}
	ru.HashedPassword = passwd
}

func (ru *RegisterUser) ValidateUsernameAvailability() (bool, error) {
	// return false, er.UserNameAlreadyExists
	return true, nil
}

func (ru *RegisterUser) IsValid() (bool, error) {

	var isUsernameAvalibe bool = false
	var isPasswordValid bool = false
	var err error = nil
	if isPasswordValid, err = ru.ValidatePassword(); err != nil {
		return false, err
	}
	if isUsernameAvalibe, err = ru.ValidateUsernameAvailability(); err != nil {
		return false, err
	}

	return isUsernameAvalibe && isPasswordValid, nil
}

func (ru *RegisterUser) ValidatePassword() (bool, error) {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	if len(ru.Password) >= 7 {
		hasMinLen = true
	}
	for _, char := range ru.Password {
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
		return true, nil
	}
	return false, er.InvalidPasswordPattern
}
