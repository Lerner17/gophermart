package handlers

import (
	"context"
	"net/http"

	"github.com/Lerner17/gophermart/internal/auth"
	"github.com/Lerner17/gophermart/internal/models"
	"github.com/labstack/echo/v4"
)

type BalanceGetter interface {
	GetUserBalance(context.Context, int) (models.Balance, error)
}

func BalanceHandler(db BalanceGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := c.Cookie("access-token")

		if err != nil {
			return err
		}

		userID, err := auth.GetUserIDFromToken(token.Value)
		if err != nil {
			return err
		}

		balance, err := db.GetUserBalance(c.Request().Context(), userID)

		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, balance)
	}
}
