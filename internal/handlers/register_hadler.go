package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Lerner17/gophermart/internal/auth"
	er "github.com/Lerner17/gophermart/internal/errors"
	"github.com/Lerner17/gophermart/internal/models"

	"github.com/labstack/echo/v4"
)

type DBRegistrator interface {
	RegisterUser(context.Context, string, string) error
}

var ErrUsernameAlreadyExists = &er.HTTPError{
	Code: 409,
	Msg:  "username already exists",
}

var ErrInvalidPassword = &er.HTTPError{
	Code: 400,
	Msg:  "invalid password pattern",
}

var ErrBindBody = &er.HTTPError{
	Code: 400,
	Msg:  "could not parse request body",
}

var ErrInvalidRequestBody = &er.HTTPError{
	Code: 400,
	Msg:  "invalid request body",
}

func Registration(db DBRegistrator) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		user := &models.User{}

		if err := c.Bind(user); err != nil {
			return fmt.Errorf("could not bind body: %v: %w", err, ErrBindBody)
		}

		if user.Login == "" || user.Password == "" {
			return ErrInvalidRequestBody
		}

		// if err := helpers.ValidatePassword(user.Password); err != nil {
		// 	return fmt.Errorf("invalid password provided for registration: %w", ErrInvalidPassword)
		// }

		if err := db.RegisterUser(ctx, user.Login, user.Password); err != nil {
			if errors.Is(err, er.UserNameAlreadyExists) {
				return ErrUsernameAlreadyExists
			}
			return fmt.Errorf("could not insert user into database: %w", err)
		}

		if err := auth.GenerateTokensAndSetCookies(user, c); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, user)
	}
}
