package handlers

import (
	"errors"
	"fmt"

	"github.com/Lerner17/gophermart/internal/auth"
	er "github.com/Lerner17/gophermart/internal/errors"
	"github.com/Lerner17/gophermart/internal/models"
	"github.com/labstack/echo/v4"
)

type DBLoginer interface {
	LoginUser(string, string) (int, error)
}

var ErrInvalidCredentials = &er.HTTPError{
	Code: 401,
	Msg:  "invalid credentials",
}

func LoginHandler(db DBLoginer) echo.HandlerFunc {
	return func(c echo.Context) error {

		user := &models.User{}
		var userID int

		if err := c.Bind(user); err != nil {
			return fmt.Errorf("could not bind body: %v:", err)
		}
		userID, err := db.LoginUser(user.Login, user.Password)
		if err != nil {
			if errors.Is(err, er.ErrInvalidLoginOrPassword) {
				return fmt.Errorf("invalid login or password: %v: %w", err, ErrInvalidCredentials)
			}
			return err
		}
		user.ID = userID
		err = auth.GenerateTokensAndSetCookies(user, c)
		if err != nil {
			panic(err)
		}
		return nil
	}
}
