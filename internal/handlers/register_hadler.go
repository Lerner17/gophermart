package handlers

import (
	"errors"
	"net/http"

	er "github.com/Lerner17/gophermart/internal/errors"
	"github.com/Lerner17/gophermart/internal/models"
	"github.com/labstack/echo/v4"
)

type DBRegistrator interface {
	RegisterUser(string, string) error
}

func Registration(db DBRegistrator) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := &models.RegisterUser{}

		if err := c.Bind(user); err != nil {
			return err
		}

		if _, err := user.IsValid(); err != nil {
			if errors.Is(err, er.UserNameAlreadyExists) {
				return c.JSON(http.StatusConflict, struct {
					Error  string `json:"error"`
					Status int    `json:"status_code"`
				}{
					Error:  "username already exists",
					Status: http.StatusConflict,
				})
			}
			if errors.Is(err, er.InvalidPasswordPattern) {
				return c.JSON(http.StatusBadRequest, struct {
					Error  string `json:"error"`
					Status int    `json:"status_code"`
				}{
					Error:  "password pattern not valid",
					Status: http.StatusBadRequest,
				})
			}
		}
		user.HashPassword()
		if err := db.RegisterUser(user.Login, user.HashedPassword); err != nil {
			panic(err)
		}
		return c.JSON(http.StatusOK, user)
	}
}
