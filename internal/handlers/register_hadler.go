package handlers

import (
	"context"
	"fmt"
	"net/http"

	er "github.com/Lerner17/gophermart/internal/errors"
	"github.com/Lerner17/gophermart/internal/helpers"
	"github.com/Lerner17/gophermart/internal/models"

	"github.com/labstack/echo/v4"
)

type DBRegistrator interface {
	RegisterUser(context.Context, string, string) error
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
		user := &models.RegisterUser{}

		if err := c.Bind(user); err != nil {
			return fmt.Errorf("%v %w", err, ErrBindBody)
		}

		if user.Login == "" || user.Password == "" {
			return ErrInvalidRequestBody
		}

		if err := helpers.ValidatePassword(user.Password); err != nil {
			return fmt.Errorf("invalid password provided for registration: %w", ErrInvalidPassword)
		}

		if err := db.RegisterUser(ctx, user.Login, user.Password); err != nil {
			return fmt.Errorf("could not insert user into database: %w", err)
		}

		return c.JSON(http.StatusOK, user)
	}
}
