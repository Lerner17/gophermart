package handlers

import (
	"context"
	"net/http"

	"github.com/Lerner17/gophermart/internal/auth"
	"github.com/Lerner17/gophermart/internal/models"
	"github.com/labstack/echo/v4"
)

type WithdrawGetter interface {
	GetWithdraws(context.Context, int) ([]models.Withdraw, error)
}

func GetWithdrawListHandler(db WithdrawGetter) echo.HandlerFunc {
	return func(c echo.Context) error {

		token, err := c.Cookie("access-token")

		if err != nil {
			return err
		}

		userID, err := auth.GetUserIDFromToken(token.Value)

		w, err := db.GetWithdraws(c.Request().Context(), userID)

		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, w)
	}
}
