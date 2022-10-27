package handlers

import (
	"context"
	"net/http"

	"github.com/Lerner17/gophermart/internal/auth"
	"github.com/Lerner17/gophermart/internal/models"
	"github.com/labstack/echo/v4"
)

type OrdesGetter interface {
	GetOrders(context.Context, int) ([]models.Order, error)
}

func GetOrdersHandler(db OrdesGetter) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := c.Cookie("access-token")

		if err != nil {
			return err
		}

		userID, err := auth.GetUserIdFromToken(token.Value)
		if err != nil {
			return err
		}

		orders, err := db.GetOrders(c.Request().Context(), userID)
		if err != nil {
			return err
		}
		if len(orders) == 0 {
			return c.JSON(204, orders)
		}
		return c.JSON(http.StatusOK, orders)
	}
}
