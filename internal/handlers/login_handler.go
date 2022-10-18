package handlers

import (
	"errors"
	"fmt"

	er "github.com/Lerner17/gophermart/internal/errors"
	"github.com/Lerner17/gophermart/internal/helpers"
	"github.com/Lerner17/gophermart/internal/models"
	"github.com/labstack/echo/v4"
)

type DBLoginer interface {
	LoginUser(string, string) error
}

var ErrInvalidCredentials = &er.HTTPError{
	Code: 401,
	Msg:  "invalid credentials",
}

func LoginHandler(db DBLoginer) echo.HandlerFunc {
	return func(c echo.Context) error {

		user := &models.User{}

		if err := c.Bind(user); err != nil {
			return fmt.Errorf("could not bind body: %v:", err)
		}
		if err := db.LoginUser(user.Login, user.Password); err != nil {
			if errors.Is(err, er.InvalidLoginOrPassword) {
				return fmt.Errorf("could not bind body: %v: %w", err, ErrInvalidCredentials)
			}
			return err
		}
		helpers.WriteCookie(c, "auth", "foo123")
		return nil
	}
}
